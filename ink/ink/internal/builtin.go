package internal

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	crand "crypto/rand"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rprtr258/fun"
)

func vs(s string) ValueString {
	b := []byte(s)
	return ValueString{&b}
}

// TODO: export as builtin.ink signatured functions

// NativeFunctionValue represents a function whose implementation is written
// in Go and built-into the runtime.
type NativeFunctionValue struct {
	name string
	exec func(*Context, Pos, []Value) Value
	ctx  *Context // runtime context to dispatch async errors
}

func (v NativeFunctionValue) String() string {
	return fmt.Sprintf("NATIVE_FUNCTION(%s)", v.name)
}

func (v NativeFunctionValue) Equals(other Value) bool {
	if _, isEmpty := other.(ValueEmpty); isEmpty {
		return true
	}

	if ov, ok := other.(NativeFunctionValue); ok {
		return v.name == ov.name
	}

	return false
}

var operators = map[Kind]func(Pos, []Value) Value{
	OpSubtract: func(pos Pos, args []Value) Value { // TODO: - :: (number -> number) | (boolean -> boolean)
		assert(len(args) == 1)

		operand := args[0]
		if isErr(operand) {
			return operand
		}

		switch o := operand.(type) {
		case ValueNumber:
			return -o
		case ValueBoolean:
			return ValueBoolean(!o)
		default:
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot negate non-boolean and non-number value %s", o), pos}}
		}
	},
}

func operatorFunc(kind Kind) NativeFunctionValue {
	f := operators[kind]
	return NativeFunctionValue{
		name: kind.String(),
		exec: func(_ *Context, p Pos, v []Value) Value {
			return f(p, v)
		},
		ctx: nil,
	}
}

var ctxBuiltins typeContext = typeContext{fun.Ptr(0), nil, map[int]Type{}}

// LoadEnvironment loads all builtins (functions and constants) to a given Context.
func (ctx *Context) LoadEnvironment() {
	// TODO: kal ebaniy, extract to var and use on init, cant be done right now because import func makes loop
	fillCtxBuiltins := ctxBuiltins.bindings == nil
	for _, sig := range []struct {
		name     string
		fn       func(*Context, Pos, []Value) Value
		tyArgs   []Type
		tyResult Type
	}{
		{"import", inkImport, []Type{typeString}, typeAny}, // TODO: string -> infer at compile time
		{"par", inkPar, []Type{typeAny}, typeAny},          // TODO: [](() -> $T) -> $T

		// system interfaces
		{"args", inkArgs, []Type{}, typeString},
		{"in", inkIn, []Type{TypeFunction{[]Type{typeAny}, typeAny}}, typeNull}, // TODO: ???
		{"out", inkOut, []Type{typeString}, typeNull},
		{"dir", inkDir, []Type{typeString, TypeFunction{[]Type{typeAny}, typeNull}}, typeNull},                               // TODO: ???
		{"make", inkMake, []Type{typeString}, typeAny},                                                                       // TODO: ???
		{"stat", inkStat, []Type{typeString}, typeAny},                                                                       // TODO: ???
		{"read", inkRead, []Type{typeString, typeNumber, typeString, TypeFunction{[]Type{typeString}, typeNull}}, typeAny},   // TODO: ???
		{"write", inkWrite, []Type{typeString, typeNumber, typeString, TypeFunction{[]Type{}, typeNull}}, typeAny},           // TODO: ???
		{"delete", inkDelete, []Type{typeString, TypeFunction{[]Type{typeAny}, typeNull}}, typeAny},                          // TODO: ???
		{"listen", inkListen, []Type{typeString, TypeFunction{[]Type{typeAny}, typeNull}}, TypeFunction{[]Type{}, typeNull}}, // TODO: (string, {type: "req", data: {method: string, url: string, headers: map[string]string, body: string | ()}} -> (string, number)) -> (() -> ())
		{"req", inkReq, []Type{typeAny}, typeAny},                                                                            // TODO: {method: string, url: string, headers: map[string]string, body: string | ()} -> ???
		{"rand", inkRand, []Type{}, typeNumber},
		{"urand", inkUrand, []Type{typeNumber}, typeString},
		{"time", inkTime, []Type{}, typeNumber},
		{"wait", inkWait, []Type{typeNumber}, typeNull},
		{"exec", inkExec, []Type{typeAny}, typeAny}, // TODO: (string, []string, string, string -> ()) -> ()|error
		{"env", inkEnv, []Type{}, typeAny},          // TODO: () -> map[string, string]
		{"exit", inkExit, []Type{typeNumber}, typeVoid},

		// math
		{"sin", inkSin, []Type{typeNumber}, typeNumber},
		{"cos", inkCos, []Type{typeNumber}, typeNumber},
		{"asin", inkAsin, []Type{typeNumber}, typeNumber},
		{"acos", inkAcos, []Type{typeNumber}, typeNumber},
		{"pow", inkPow, []Type{typeNumber, typeNumber}, typeNumber}, // TODO: handle errors effect
		{"ln", inkLn, []Type{typeNumber}, typeNumber},
		{"floor", inkFloor, []Type{typeNumber}, typeNumber},

		// type conversions
		{"string", inkString, []Type{typeAny}, typeString},
		{"number", inkNumber, []Type{typeAny}, typeNumber}, // TODO: string -> number|error, number -> number, bool -> 0|1
		{"point", inkPoint, []Type{typeString}, typeNumber},
		{"char", inkChar, []Type{typeNumber}, typeString},

		// introspection
		{"type", inkType, []Type{typeAny}, typeString},
		{"len", inkLen, []Type{typeAny}, typeAny},   // TODO: []$T -> number, map[$K, $V] -> number
		{"keys", inkKeys, []Type{typeAny}, typeAny}, // TODO: []$T -> []number, map[$K, $V] -> []$K
	} {
		if fillCtxBuiltins {
			ctxBuiltins = ctxBuiltins.Append(sig.name, TypeFunction{sig.tyArgs, sig.tyResult})
		}
		ctx.LoadFunc(sig.name, sig.fn)
	}

	// side effects
	rand.Seed(time.Now().UTC().UnixNano())
}

// LoadFunc loads a single Go-implemented function into a Context.
func (ctx *Context) LoadFunc(
	name string,
	exec func(*Context, Pos, []Value) Value,
) {
	ctx.Scope.Set(name, NativeFunctionValue{name, exec, ctx})
}

// Create and return a standard error callback response with the given message
func errMsg(message string) ValueComposite {
	return ValueComposite{
		"type":    vs("error"),
		"message": vs(message),
	}
}

// TODO: remove type checks, just assert them
func validate(pos Pos, errs ...string) *Err {
	for _, err := range errs {
		if err != "" {
			return &Err{nil, ErrAssert, err, pos}
		}
	}
	return nil
}

func validateArgsLen(in []Value, expected int) string {
	if len(in) != expected {
		return fmt.Sprintf("takes expected %d arguments, but got %d", expected, len(in))
	}
	return ""
}

func validateArgType[T any](in []Value, i int, dest *T) string {
	if i >= len(in) {
		// skip like it is optional
		return ""
	}

	res, ok := in[i].(T)
	if !ok {
		if err, ok := in[i].(ValueError); ok {
			return fmt.Sprintf(
				"%d-th argument must be %T, but got ERROR: %s",
				i, *new(T), err.Error(),
			)
		}
		return fmt.Sprintf(
			"%d-th argument must be %T, but got %T",
			i, *new(T), in[i],
		)
	}

	*dest = res
	return ""
}

func validateArgListOf[T any](in []Value, i int, dest *[]T) string {
	if i >= len(in) {
		// skip like it is optional
		return ""
	}

	var args ValueList
	if err := validateArgType(in, i, &args); err != "" {
		return err
	}

	argsList := make([]T, len(*args.xs))
	for i, v := range *args.xs {
		if a, ok := v.(T); ok {
			argsList[i] = a
		} else {
			return fmt.Sprintf(
				"%d-th argument must contain %Ts, got %s",
				i, *new(T), v)
		}
	}

	*dest = argsList
	return ""
}

func validateCustom(condition bool, msg string) string {
	if !condition {
		return msg
	}

	return ""
}

func inkImport(ctx *Context, pos Pos, in []Value) Value {
	var givenPath ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &givenPath),
		validateCustom(len(*givenPath.b) > 0, "arg must be path"),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "import()", pos}}
	}

	// imports via import() are assumed to be relative
	// TODO: separate type and operations over import paths
	importPath := string(*givenPath.b)
	if u, err := url.Parse(importPath); err == nil && u.Scheme != "" {
	} else if !filepath.IsAbs(importPath) {
		if u, err := url.Parse(ctx.WorkingDirectory); err == nil && u.Scheme != "" {
			u.Path = path.Join(path.Dir(u.Path), importPath)
			importPath = u.String()
		} else {
			importPath = filepath.Join(ctx.WorkingDirectory, importPath)
		}
	}

	// evalLock blocks file eval; temporary unlock it for the import to run.
	// Calling import() from within a running program is not supported, so we
	// don't really care if catastrophic things happen because of unlocked evalLock.
	// ctx.Engine.mu.Unlock()
	// defer ctx.Engine.mu.Lock()

	if _, ok := ctx.Engine.Contexts[importPath]; !ok {
		// The imported program runs in a "child context", a distinct context from
		// the importing program. The "child" term is a bit of a misnomer as Contexts
		// do not exist in a hierarchy, but conceptually makes sense here.
		childCtx := ctx.Engine.CreateContext()
		ctx.Engine.Contexts[importPath] = childCtx

		// Execution here follows updating ctx.Engine.Contexts
		// to behave correctly in the case where A imports B imports A again,
		// and still only import one instance of A.
		nodes, err := childCtx.ExecPath(importPath)
		if err != nil {
			return ValueError{&Err{err, ErrRuntime, fmt.Sprintf("error importing file %s", importPath), pos}}
		}
		value := childCtx.Eval(nodes)
		if err, ok := value.(ValueError); ok {
			return ValueError{&Err{err.Err, ErrRuntime, fmt.Sprintf("error evaluating importing file %s", importPath), pos}}
		}

		ctx.Engine.values[importPath] = value
	}

	return ctx.Engine.values[importPath]
}

func inkArgs(_ *Context, _ Pos, _ []Value) Value {
	comp := make([]Value, len(os.Args))
	for i, v := range os.Args {
		comp[i] = vs(v)
	}
	return ValueList{&comp}
}

func inkIn(ctx *Context, pos Pos, in []Value) Value {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "in()", pos}}
	}

	cbErr := func(err *Err) {
		LogError(&Err{err, ErrRuntime, "error in callback to in()", pos})
	}

	ctx.ExecListener(func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				// also captures io.EOF
				break
			}

			rv := evalInkFunction(ctx, in[0], pos, ValueComposite{
				"type": vs("data"),
				"data": vs(str),
			})
			if errEval, ok := rv.(ValueError); ok {
				cbErr(errEval.Err)
				return
			}

			if boolValue, isBool := rv.(ValueBoolean); isBool {
				if !boolValue {
					break
				}
			} else {
				LogError(&Err{nil, ErrRuntime, fmt.Sprintf("callback to in() should return a boolean, but got %s", rv), pos})
				return
			}
		}

		if err, ok := evalInkFunction(ctx, in[0], pos, ValueComposite{
			"type": vs("end"),
		}).(ValueError); ok {
			cbErr(err.Err)
		}
	})

	return Null
}

// TODO: replace with write('/proc/self/fd/1', ~1, output, e => ())
func inkOut(ctx *Context, pos Pos, in []Value) Value {
	var output ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &output),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "out()", pos}}
	}

	os.Stdout.Write(*output.b)
	return Null
}

func inkDir(ctx *Context, pos Pos, in []Value) Value {
	var (
		dirPath ValueString
		cb      ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &dirPath),
		validateArgType(in, 1, &cb),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "dir()", pos}}
	}

	cbMaybeErr := func(v Value) {
		if err, ok := v.(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to dir()", pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		fileInfos, err := os.ReadDir(string(*dirPath.b))
		if err != nil {
			ctx.ExecListener(func() {
				cbMaybeErr(evalInkFunction(ctx, cb, pos, errMsg(
					fmt.Sprintf("error listing directory contents in dir(), %s", err.Error()),
				)))
			})
			return
		}

		fileList := make(ValueComposite, len(fileInfos))
		for i, fi := range fileInfos {
			info, err := fi.Info()
			if err != nil {
				ctx.ExecListener(func() {
					cbMaybeErr(evalInkFunction(ctx, cb, pos, errMsg(
						fmt.Sprintf("error listing directory contents in dir(), %s", err.Error()),
					)))
				})
				return
			}

			fileList[strconv.Itoa(i)] = ValueComposite{
				"name": vs(info.Name()),
				"len":  ValueNumber(info.Size()),
				"dir":  ValueBoolean(info.IsDir()),
				"mod":  ValueNumber(info.ModTime().Unix()),
			}
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ctx, cb, pos, ValueComposite{
				"type": vs("data"),
				"data": fileList,
			}))
		})
	}()

	return Null
}

func inkMake(ctx *Context, pos Pos, in []Value) Value {
	var dirPath ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &dirPath),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "make()", pos}}
	}

	err := os.MkdirAll(string(*dirPath.b), 0o755)
	if err != nil {
		return errMsg(
			fmt.Sprintf("error making a new directory in make(), %s", err.Error()),
		)
	}

	return ValueComposite{
		"type": vs("end"),
	}
}

func inkStat(ctx *Context, pos Pos, in []Value) Value {
	var statPath ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &statPath),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "stat()", pos}}
	}

	var res Value
	fi, err := os.Stat(string(*statPath.b))
	if err != nil {
		if os.IsNotExist(err) {
			res = ValueComposite{
				"type": vs("data"),
				"data": Null,
			}
		} else {
			res = errMsg(fmt.Sprintf("error getting file data in stat(), %s", err.Error()))
		}
	} else {
		res = ValueComposite{
			"type": vs("data"),
			"data": ValueComposite{
				"name": vs(fi.Name()),
				"len":  ValueNumber(fi.Size()),
				"dir":  ValueBoolean(fi.IsDir()),
				"mod":  ValueNumber(fi.ModTime().Unix()),
			},
		}
	}

	return res
}

func inkRead(ctx *Context, pos Pos, in []Value) Value {
	var (
		filePath ValueString
		offset   ValueNumber
		length   ValueNumber
		cb       ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 4),
		validateArgType(in, 0, &filePath),
		validateArgType(in, 1, &offset),
		validateArgType(in, 2, &length),
		validateArgType(in, 3, &cb),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "read()", pos}}
	}

	cbMaybeErr := func(v Value) {
		if err, ok := v.(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to read()", pos})
		}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ctx, cb, pos, errMsg(msg)))
		})
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// open
		file, err := os.OpenFile(string(*filePath.b), os.O_RDONLY, 0o644)
		if err != nil {
			sendErr(fmt.Sprintf("error opening requested file in read(), %s", err.Error()))
			return
		}
		defer file.Close()

		// seek
		ofs := int64(offset)
		if ofs != 0 {
			_, err := file.Seek(ofs, 0) // 0 means relative to start of file
			if err != nil {
				sendErr(fmt.Sprintf("error seeking requested file in read(), %s", err.Error()))
				return
			}
		}

		// read
		buf := make([]byte, int64(length))
		count, err := file.Read(buf)
		if err == io.EOF && count == 0 {
			// if first read returns EOF, it may just be an empty file
			ctx.ExecListener(func() {
				cbMaybeErr(evalInkFunction(ctx, cb, pos, ValueComposite{
					"type": vs("data"),
					"data": ValueString{},
				}))
			})
			return
		} else if err != nil {
			sendErr(fmt.Sprintf("error reading requested file in read(), %s", err.Error()))
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ctx, cb, pos, ValueComposite{
				"type": vs("data"),
				"data": vs(string(buf[:count])),
			}))
		})
	}()

	return Null
}

func inkWrite(ctx *Context, pos Pos, in []Value) Value {
	var (
		filePath ValueString
		offset   ValueNumber
		buf      ValueString
		cb       ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 4),
		validateArgType(in, 0, &filePath),
		validateArgType(in, 1, &offset),
		validateArgType(in, 2, &buf),
		validateArgType(in, 3, &cb),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "write()", pos}}
	}

	cbMaybeErr := func(v Value) {
		if err, ok := v.(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to write()", pos})
		}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ctx, cb, pos, errMsg(msg)))
		})
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// open
		var flag int
		if offset == -1 {
			// -1 offset is append
			flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		} else {
			// all other offsets are writing
			flag = os.O_CREATE | os.O_WRONLY
		}
		file, err := os.OpenFile(string(*filePath.b), flag, 0o644)
		if err != nil {
			sendErr(fmt.Sprintf("error opening requested file in write(), %s", err.Error()))
			return
		}
		defer file.Close()

		// seek
		if offset != -1 {
			ofs := int64(offset)
			if _, err := file.Seek(ofs, 0); err != nil { // 0 means relative to start of file
				sendErr(fmt.Sprintf("error seeking requested file in write(), %s", err.Error()))
				return
			}
		}

		// write
		if _, err := file.Write(*buf.b); err != nil {
			sendErr(fmt.Sprintf("error writing to requested file in write(), %s", err.Error()))
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ctx, cb, pos, ValueComposite{
				"type": vs("end"),
			}))
		})
	}()

	return Null
}

func inkDelete(ctx *Context, pos Pos, in []Value) Value {
	var (
		filePath ValueString
		cb       ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &filePath),
		validateArgType(in, 1, &cb),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "delete()", pos}}
	}

	cbMaybeErr := func(v Value) {
		if err, ok := v.(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to delete()", pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// delete
		err := os.RemoveAll(string(*filePath.b))
		if err != nil {
			ctx.ExecListener(func() {
				cbMaybeErr(evalInkFunction(ctx, cb, pos, errMsg(
					fmt.Sprintf("error removing requested file in delete(), %s", err.Error()),
				)))
			})
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ctx, cb, pos, ValueComposite{
				"type": vs("end"),
			}))
		})
	}()

	return Null
}

// inkHTTPHandler fulfills the Handler interface for inkListen() to work
func inkHTTPHandler(
	ctx *Context,
	cb ValueFunction,
) http.HandlerFunc {
	pos := Pos{ctx.File, 0, 0} // TODO: pass position here and everywhere

	return func(w http.ResponseWriter, r *http.Request) {
		// unmarshal request
		method := r.Method
		url := r.URL.String()

		headers := make(ValueComposite, len(r.Header))
		for key, values := range r.Header {
			headers[key] = vs(strings.Join(values, ","))
		}

		body := ValueString{new([]byte)}
		if r.ContentLength != 0 {
			bodyBuf, err := io.ReadAll(r.Body)
			if err != nil {
				ctx.ExecListener(func() {
					_ = evalInkFunction(ctx, cb, pos, errMsg(
						fmt.Sprintf("error reading request in listen(), %s", err.Error()),
					))
				})
				return
			}
			body = ValueString{&bodyBuf}
		}

		// construct request object to pass to Ink, and call handler
		responseEnded := false
		responses := make(chan Value, 1)

		ctx.Engine.Listeners.Add(1)
		go func() {
			defer ctx.Engine.Listeners.Done()

			v := evalInkFunction(ctx, cb, pos, ValueComposite{
				"type": vs("req"),
				"data": ValueComposite{
					"method":  vs(method),
					"url":     vs(url),
					"headers": headers,
					"body":    body,
				},
			})
			// if len(v) != 1 {
			// 	LogError(&Err{nil, ErrRuntime, "end() callback to listen() must have one argument", pos})
			// }
			if responseEnded {
				LogError(&Err{nil, ErrRuntime, "end() callback to listen() was called more than once", pos})
			}
			responseEnded = true
			responses <- v
		}()

		// validate response from Ink callback
		resp := <-responses
		rsp, isComposite := resp.(ValueComposite)
		if !isComposite {
			LogError(&Err{nil, ErrRuntime, fmt.Sprintf("callback to listen() should return a response, got %s", resp), pos})
			return
		}

		// unmarshal response from the return value
		// response = {status, headers, body}
		resStatus, okStatus := rsp["status"].(ValueNumber)
		resHeaders, okHeaders := rsp["headers"].(ValueComposite)
		resBody, okBody := rsp["body"].(ValueString)

		if !okStatus || !okHeaders || !okBody {
			LogError(&Err{nil, ErrRuntime, fmt.Sprintf("callback to listen() returned malformed response, %s", rsp), pos})
			return
		}

		// write values to response
		// Content-Length is automatically set for us by Go
		for k, v := range resHeaders {
			if str, isStr := v.(ValueString); isStr {
				w.Header().Set(k, string(*str.b))
			} else {
				LogError(&Err{nil, ErrRuntime, fmt.Sprintf("could not set response header, value %s was not a string", v), pos})
				return
			}
		}

		code := int(resStatus)
		// guard against invalid HTTP codes, which cause Go panics.
		// https://golang.org/src/net/http/server.go
		if code < 100 || code > 599 {
			LogError(&Err{nil, ErrRuntime, fmt.Sprintf("could not set response status code, code %d is not valid", code), pos})
			return
		}

		// status code write must follow all other header writes,
		// since it sends the status
		w.WriteHeader(int(resStatus))
		if _, err := w.Write(*resBody.b); err != nil {
			_ = evalInkFunction(ctx, cb, pos, errMsg(
				fmt.Sprintf("error writing request body in listen(), %s", err.Error()),
			))
			return
		}
	}
}

func inkListen(ctx *Context, pos Pos, in []Value) Value {
	var (
		host ValueString
		cb   ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &host),
		validateArgType(in, 1, &cb),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "listen()", pos}}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			if err, ok := evalInkFunction(ctx, cb, pos, errMsg(msg)).(ValueError); ok {
				LogError(&Err{err.Err, ErrRuntime, "error in callback to listen()", pos})
			}
		})
	}

	server := &http.Server{
		Addr:    string(*host.b),
		Handler: http.HandlerFunc(inkHTTPHandler(ctx, cb)),
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()
		fmt.Fprintf(os.Stderr, "listening on %s\n", string(server.Addr))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			sendErr(fmt.Sprintf("error starting http server in listen(), %s", err.Error()))
		}
	}()

	closer := func(ctx *Context, _ Pos, in []Value) Value {
		// attempt graceful shutdown, concurrently, without
		// blocking Ink evaluation thread
		if err := server.Shutdown(context.Background()); err != nil {
			sendErr(fmt.Sprintf("error closing server in listen(), %s", err.Error()))
		}
		return Null
	}

	return NativeFunctionValue{
		name: "close",
		exec: closer,
		ctx:  ctx,
	}
}

func inkReq(ctx *Context, pos Pos, in []Value) Value {
	var data ValueComposite
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &data),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "req()", pos}}
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// do not follow redirects
			return http.ErrUseLastResponse
		},
	}
	reqContext, reqCancel := context.WithCancel(context.Background())
	defer reqCancel()

	// closer := func(_ *Context, _ Pos, _ []Value) Value {
	// 	reqCancel()
	// 	return Null
	// }

	// unmarshal request contents
	methodVal := cmp.Or(data["method"], Value(vs("GET")))
	urlVal := data["url"]
	headersVal := cmp.Or(data["headers"], Value(ValueComposite{}))
	bodyVal := cmp.Or(data["body"], Value(vs("")))

	// TODO: add query params

	reqMethod, okMethod := methodVal.(ValueString)
	reqURL, okURL := urlVal.(ValueString)
	reqHeaders, okHeaders := headersVal.(ValueComposite)
	reqBody, okBody := bodyVal.(ValueString)

	if !okMethod || !okURL || !okHeaders || !okBody {
		return errMsg(fmt.Sprintf("request in req() is malformed, %s", data))
	}
	req, err := http.NewRequestWithContext(
		reqContext,
		string(*reqMethod.b),
		string(*reqURL.b),
		bytes.NewReader(*reqBody.b),
	)
	if err != nil {
		return errMsg(fmt.Sprintf("error creating request in req(), %s", err.Error()))
	}

	// construct headers
	// Content-Length is automatically set for us by Go
	req.Header.Set("User-Agent", "") // remove Go's default user agent header
	for k, v := range reqHeaders {
		if str, isStr := v.(ValueString); isStr {
			req.Header.Set(k, string(*str.b))
		} else {
			return errMsg(fmt.Sprintf("could not set request header, value %s was not a string", v))
		}
	}

	// send request
	resp, err := client.Do(req)
	if err != nil {
		return errMsg(fmt.Sprintf("error processing request in req(), %s", err.Error()))
	}
	defer resp.Body.Close()

	resStatus := ValueNumber(resp.StatusCode)
	resHeaders := make(ValueComposite, len(resp.Header))
	for key, values := range resp.Header {
		resHeaders[key] = vs(strings.Join(values, ","))
	}

	var resBody Value
	if resp.ContentLength == 0 {
		resBody = ValueString{}
	} else {
		bodyBuf, err := io.ReadAll(resp.Body)
		if err != nil {
			return errMsg(fmt.Sprintf("error reading response in req(), %s", err.Error()))
		}
		resBody = ValueString{&bodyBuf}
	}

	return ValueComposite{
		"type": vs("resp"),
		"data": ValueComposite{
			"status":  resStatus,
			"headers": resHeaders,
			"body":    resBody,
		},
	}

	// return NativeFunctionValue{
	// 	name: "close",
	// 	exec: closer,
	// 	ctx:  ctx,
	// }
}

func inkRand(ctx *Context, pos Pos, in []Value) Value {
	return ValueNumber(rand.Float64())
}

// TODO: rewrite in user-space using rand
func inkUrand(ctx *Context, pos Pos, in []Value) Value {
	var bufLength ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &bufLength),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "urand()", pos}}
	}

	buf := make([]byte, int64(bufLength))
	if _, err := crand.Read(buf); err != nil {
		return Null
	}

	return ValueString{&buf}
}

func inkTime(ctx *Context, pos Pos, in []Value) Value {
	unixSeconds := float64(time.Now().UnixNano()) / 1e9
	return ValueNumber(unixSeconds)
}

func inkWait(ctx *Context, pos Pos, in []Value) Value {
	var secs ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &secs),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "wait()", pos}}
	}

	// This is a bit tricky, since we don't want wait() to hold the evalLock
	// on the Context while we're waiting for the timeout, but do want to hold
	// the main goroutine from completing with sync.WaitGroup.
	time.Sleep(time.Duration(
		int64(float64(secs) * float64(time.Second)),
	))

	return Null
}

func inkExec(ctx *Context, pos Pos, in []Value) Value {
	var (
		path     ValueString
		args     []ValueString
		stdin    ValueString
		stdoutFn ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 4),
		validateArgType(in, 0, &path),
		validateArgListOf(in, 1, &args),
		validateArgType(in, 2, &stdin),
		validateArgType(in, 3, &stdoutFn),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "exec()", pos}}
	}

	argsList := make([]string, len(args))
	for i, v := range args {
		argsList[i] = string(*v.b)
	}

	cmd := exec.Command(string(*path.b), argsList...)
	// cmdMutex locks control over reading and modifying child
	// process state, because both the Ink eval thread and exec
	// thread must read from/write to cmd.
	cmdMutex := sync.Mutex{} // TODO: remove as much mutexes as possible
	stdout := bytes.Buffer{}
	cmd.Stdin = strings.NewReader(string(*stdin.b))
	cmd.Stdout = &stdout

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			v := evalInkFunction(ctx, stdoutFn, pos, errMsg(msg))
			if err, ok := v.(ValueError); ok {
				LogError(&Err{err.Err, ErrRuntime, "error in callback to exec()", pos})
			}
		})
	}

	runAndExit := func() {
		cmdMutex.Lock()
		err := cmd.Start()
		cmdMutex.Unlock()
		if err != nil {
			sendErr(fmt.Sprintf("error starting command in exec(), %s", err.Error()))
			return
		}

		if err = cmd.Wait(); err != nil {
			// if there is an err but err is just ExitErr, this means
			// the process ran successfully but exited with an error code.
			// We consider this ok and keep going.
			if _, isExitErr := err.(*exec.ExitError); !isExitErr {
				sendErr(fmt.Sprintf("error waiting for command to exit in exec(), %s", err.Error()))
				return
			}
		}

		output, err := io.ReadAll(&stdout)
		if err != nil {
			sendErr(fmt.Sprintf("error reading command output in exec(), %s", err.Error()))
			return
		}

		ctx.ExecListener(func() {
			in := ValueComposite{
				"type": vs("data"),
				"data": ValueString{&output},
			}
			v := evalInkFunction(ctx, stdoutFn, pos, in)
			if err, ok := v.(ValueError); ok {
				LogError(&Err{err.Err, ErrRuntime, "error in callback to exec()", pos})
			}
		})
	}

	// if the caller closes the cmd before it ever starts running,
	// we need to signal that safely to the cmd-running goroutine
	neverRun := make(chan bool, 1)
	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		select {
		case <-neverRun:
			// do nothing
		default:
			runAndExit()
		}
	}()

	closed := false
	return NativeFunctionValue{
		name: "close",
		exec: func(_ *Context, pos Pos, _ []Value) Value {
			// multiple calls to close() should be idempotent
			if closed {
				return Null
			}

			neverRun <- true
			closed = true

			cmdMutex.Lock()
			if cmd.Process != nil || cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				if err := cmd.Process.Kill(); err != nil {
					return ValueError{&Err{nil, ErrRuntime, err.Error(), pos}}
				}
			}
			cmdMutex.Unlock()

			return Null
		},
		ctx: ctx,
	}
}

func inkEnv(ctx *Context, pos Pos, in []Value) Value {
	envp := os.Environ()

	envVars := make(ValueComposite, len(envp))
	for _, e := range envp {
		kv := strings.SplitN(e, "=", 2)
		envVars[kv[0]] = vs(kv[1])
	}
	return envVars
}

func inkExit(ctx *Context, pos Pos, in []Value) Value {
	var code ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &code),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "exit()", pos}}
	}

	os.Exit(int(code))
	return nil // unreachable
}

func inkSin(ctx *Context, pos Pos, in []Value) Value {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "sin()", pos}}
	}

	return ValueNumber(math.Sin(float64(inNum)))
}

func inkCos(ctx *Context, pos Pos, in []Value) Value {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "cos()", pos}}
	}

	return ValueNumber(math.Cos(float64(inNum)))
}

func inkAsin(ctx *Context, pos Pos, in []Value) Value {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
		validateCustom(inNum >= -1 && inNum <= 1, fmt.Sprintf("number must be in range [-1, 1], got %v", inNum)),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "asin()", pos}}
	}

	return ValueNumber(math.Asin(float64(inNum)))
}

func inkAcos(ctx *Context, pos Pos, in []Value) Value {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
		validateCustom(inNum >= -1 && inNum <= 1, fmt.Sprintf("number must be in range [-1, 1], got %v", inNum)),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "acos()", pos}}
	}

	return ValueNumber(math.Acos(float64(inNum)))
}

func inkPow(ctx *Context, pos Pos, in []Value) Value {
	var base, exp ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &base),
		validateArgType(in, 1, &exp),
		validateCustom(base != 0 || exp != 0, "math error, pow(0, 0) is not defined"),
		validateCustom(base >= 0 || isInteger(exp), "math error, fractional power of negative number"),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "pow()", pos}}
	}

	return ValueNumber(math.Pow(float64(base), float64(exp)))
}

func inkLn(ctx *Context, pos Pos, in []Value) Value {
	var n ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &n),
		validateCustom(n > 0, fmt.Sprintf("cannot take natural logarithm of non-positive number %s", n.String())),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "ln()", pos}}
	}

	return ValueNumber(math.Log(float64(n)))
}

func inkFloor(ctx *Context, pos Pos, in []Value) Value {
	var n ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &n),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "floor()", pos}}
	}

	return ValueNumber(math.Trunc(float64(n)))
}

func inkString(ctx *Context, pos Pos, in []Value) Value {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "string()", pos}}
	}

	s := func() string {
		switch v := in[0].(type) {
		case ValueError:
			return v.Error()
		case ValueString:
			return string(*v.b)
		case ValueNumber:
			return v.String()
		case ValueBoolean:
			return fun.IF(bool(v), "true", "false")
		case ValueNull:
			return "()"
		case ValueComposite:
			return v.String()
		case ValueList:
			return v.String()
		case ValueFunction, NativeFunctionValue:
			return "(function)"
		default:
			return ""
		}
	}()
	return vs(s)
}

func inkNumber(ctx *Context, pos Pos, in []Value) Value {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "number()", pos}}
	}

	switch v := in[0].(type) {
	case ValueString:
		f, err := strconv.ParseFloat(string(*v.b), 64)
		if err != nil {
			return Null
		}
		return ValueNumber(f)
	case ValueNumber:
		return v
	case ValueBoolean:
		return ValueNumber(fun.IF[float64](bool(v), 1, 0))
	default:
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cant convert %v to number", v), pos}}
	}
}

func inkPoint(ctx *Context, pos Pos, in []Value) Value {
	var str ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &str),
		validateCustom(len(*str.b) >= 1, "argument must be of length at least 1"),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "point()", pos}}
	}

	return ValueNumber((*str.b)[0])
}

func inkChar(ctx *Context, pos Pos, in []Value) Value {
	var cp ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &cp),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "char()", pos}}
	}

	return vs(string([]byte{byte(cp)}))
}

func inkType(ctx *Context, pos Pos, in []Value) Value {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return ValueError{&Err{err, ErrAssert, "type()", pos}}
	}

	switch in[0].(type) {
	case ValueString:
		return vs("string")
	case ValueNumber:
		return vs("number")
	case ValueBoolean:
		return vs("boolean")
	case ValueNull:
		return vs("()")
	case ValueComposite:
		return vs("composite")
	case ValueList:
		return vs("list")
	case ValueFunction, NativeFunctionValue:
		return vs("function")
	case ValueError:
		return vs("error")
	default:
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("unknown type: %T", in[0]), pos}}
	}
}

func inkLen(ctx *Context, pos Pos, in []Value) Value {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "len()", pos}}
	}

	switch v := in[0].(type) {
	case ValueComposite:
		return ValueNumber(len(v))
	case ValueList:
		return ValueNumber(len(*v.xs))
	case ValueString:
		// TODO: bytes count/rune count/grapheme clusters count?
		return ValueNumber(len(*v.b))
	default:
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("len() takes a string or list or composite value, but got %s", in[0].String()), pos}}
	}
}

func inkKeys(ctx *Context, pos Pos, in []Value) Value {
	// var obj ValueComposite
	if err := validate(pos,
		validateArgsLen(in, 1),
		// validateArgType(in, 0, &obj),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "keys()", pos}}
	}

	switch obj := in[0].(type) {
	case ValueComposite:
		keys := make([]Value, len(obj))
		i := 0
		for k := range obj {
			keys[i] = vs(k)
			i++
		}
		return ValueList{&keys}
	case ValueList:
		indices := make([]Value, len(*obj.xs))
		for i := range *obj.xs {
			indices[i] = ValueNumber(i)
		}
		return ValueList{&indices}
	default:
		return ValueError{&Err{nil, ErrRuntime, "keys()", pos}}
	}
}

func inkPar(ctx *Context, pos Pos, in []Value) Value {
	var funcs ValueList
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &funcs),
	); err != nil {
		return ValueError{&Err{err, ErrRuntime, "par()", pos}}
	}

	// evalLock blocks file eval; temporary unlock it for the import to run.
	// Calling import() from within a running program is not supported, so we
	// don't really care if catastrophic things happen because of unlocked evalLock.
	// ctx.Engine.mu.Unlock()
	// defer ctx.Engine.mu.Lock()

	var wg sync.WaitGroup
	wg.Add(len(*funcs.xs))
	for _, f := range *funcs.xs {
		go func() {
			_ = evalInkFunction(ctx, f, pos)
			wg.Done()
		}()
	}
	wg.Wait()

	// TODO: composite of results
	return Null
}
