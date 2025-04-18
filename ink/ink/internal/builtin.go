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

// TODO: export as builtin.ink signatured functions

// NativeFunctionValue represents a function whose implementation is written
// in Go and built-into the runtime.
type NativeFunctionValue struct {
	name string
	exec func(*Context, Pos, []Value, Cont) ValueThunk
	ctx  *Context // runtime context to dispatch async errors
}

func (v NativeFunctionValue) String() string {
	return fmt.Sprintf("Native Function (%s)", v.name)
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

// LoadEnvironment loads all builtins (functions and constants) to a given Context.
func (ctx *Context) LoadEnvironment() {
	for name, fn := range map[string]func(*Context, Pos, []Value, Cont) ValueThunk{
		"import": inkImport,
		"par":    inkPar,

		// system interfaces
		"args":   inkArgs,
		"in":     inkIn,
		"out":    inkOut,
		"dir":    inkDir,
		"make":   inkMake,
		"stat":   inkStat,
		"read":   inkRead,
		"write":  inkWrite,
		"delete": inkDelete,
		"listen": inkListen,
		"req":    inkReq,
		"rand":   inkRand,
		"urand":  inkUrand,
		"time":   inkTime,
		"wait":   inkWait,
		"exec":   inkExec,
		"env":    inkEnv,
		"exit":   inkExit,

		// math
		"sin":   inkSin,
		"cos":   inkCos,
		"asin":  inkAsin,
		"acos":  inkAcos,
		"pow":   inkPow,
		"ln":    inkLn,
		"floor": inkFloor,

		// type conversions
		"string": inkString,
		"number": inkNumber,
		"point":  inkPoint,
		"char":   inkChar,

		// introspection
		"type": inkType,
		"len":  inkLen,
		"keys": inkKeys,
	} {
		ctx.LoadFunc(name, fn)
	}

	// side effects
	rand.Seed(time.Now().UTC().UnixNano())
}

// LoadFunc loads a single Go-implemented function into a Context.
func (ctx *Context) LoadFunc(
	name string,
	exec func(*Context, Pos, []Value, Cont) ValueThunk,
) {
	ctx.Scope.Set(name, NativeFunctionValue{name, exec, ctx})
}

// Create and return a standard error callback response with the given message
func errMsg(message string) ValueComposite {
	return ValueComposite{
		"type":    ValueString("error"),
		"message": ValueString(message),
	}
}

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

func inkImport(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var givenPath ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &givenPath),
		validateCustom(len(givenPath) > 0, "arg must be path"),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "import()", pos}})
	}

	// imports via import() are assumed to be relative
	// TODO: separate type and operations over import paths
	importPath := string(givenPath)
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
	ctx.Engine.mu.Unlock()
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
			ctx.Engine.mu.Lock()
			return k(ValueError{&Err{err, ErrRuntime, fmt.Sprintf("error importing file %s", importPath), pos}})
		}
		value, err := childCtx.Eval(nodes)
		if err != nil {
			ctx.Engine.mu.Lock()
			return k(ValueError{&Err{err, ErrRuntime, fmt.Sprintf("error evaluating importing file %s", importPath), pos}})
		}

		ctx.Engine.values[importPath] = value
	}

	ctx.Engine.mu.Lock()
	return k(ctx.Engine.values[importPath])
}

func inkArgs(_ *Context, _ Pos, _ []Value, k Cont) ValueThunk {
	comp := make(ValueComposite, len(os.Args))
	for i, v := range os.Args {
		comp[strconv.Itoa(i)] = ValueString(v)
	}
	return k(comp)
}

func inkIn(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "in()", pos}})
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

			rv := trampoline(evalInkFunction(ast, in[0], pos, k, ValueComposite{
				"type": ValueString("data"),
				"data": ValueString(str),
			}))
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

		_ = trampoline(evalInkFunction(ast, in[0], pos, func(err Value) ValueThunk {
			if isErr(err) {
				cbErr(err.(ValueError).Err)
			}
			return k(Null)
		}, ValueComposite{
			"type": ValueString("end"),
		})).(ValueError)
	})

	return k(Null)
}

// TODO: replace with write('/proc/self/fd/1', ~1, output, e => ())
func inkOut(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var output ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &output),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "out()", pos}})
	}

	os.Stdout.Write([]byte(output))
	return k(Null)
}

func inkDir(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var (
		dirPath ValueString
		cb      ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &dirPath),
		validateArgType(in, 1, &cb),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "dir()", pos}})
	}

	cbMaybeErr := func(v ValueThunk) {
		if err, ok := trampoline(v).(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to dir()", pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		fileInfos, err := os.ReadDir(string(dirPath))
		if err != nil {
			ctx.ExecListener(func() {
				cbMaybeErr(evalInkFunction(ast, cb, pos, k, errMsg(
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
					cbMaybeErr(evalInkFunction(ast, cb, pos, k, errMsg(
						fmt.Sprintf("error listing directory contents in dir(), %s", err.Error()),
					)))
				})
				return
			}

			fileList[strconv.Itoa(i)] = ValueComposite{
				"name": ValueString(info.Name()),
				"len":  ValueNumber(info.Size()),
				"dir":  ValueBoolean(info.IsDir()),
				"mod":  ValueNumber(info.ModTime().Unix()),
			}
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
				"type": ValueString("data"),
				"data": fileList,
			}))
		})
	}()

	return k(Null)
}

func inkMake(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var (
		dirPath ValueString
		cb      ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &dirPath),
		validateArgType(in, 1, &cb),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "make()", pos}})
	}

	cbMaybeErr := func(v ValueThunk) {
		if err, ok := trampoline(v).(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to make()", pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		err := os.MkdirAll(string(dirPath), 0o755)
		if err != nil {
			ctx.ExecListener(func() {
				cbMaybeErr(evalInkFunction(ast, cb, pos, k, errMsg(
					fmt.Sprintf("error making a new directory in make(), %s", err.Error()),
				)))
			})
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
				"type": ValueString("end"),
			}))
		})
	}()

	return k(Null)
}

func inkStat(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var (
		statPath ValueString
		cb       ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &statPath),
		validateArgType(in, 1, &cb),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "stat()", pos}})
	}

	cbMaybeErr := func(v ValueThunk) {
		if err, ok := trampoline(v).(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to stat()", pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		fi, err := os.Stat(string(statPath))
		if err != nil {
			if os.IsNotExist(err) {
				ctx.ExecListener(func() {
					cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
						"type": ValueString("data"),
						"data": Null,
					}))
				})
			} else {
				ctx.ExecListener(func() {
					cbMaybeErr(evalInkFunction(ast, cb, pos, k, errMsg(
						fmt.Sprintf("error getting file data in stat(), %s", err.Error()),
					)))
				})
			}
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
				"type": ValueString("data"),
				"data": ValueComposite{
					"name": ValueString(fi.Name()),
					"len":  ValueNumber(fi.Size()),
					"dir":  ValueBoolean(fi.IsDir()),
					"mod":  ValueNumber(fi.ModTime().Unix()),
				},
			}))
		})
	}()

	return k(Null)
}

func inkRead(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
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
		return k(ValueError{&Err{err, ErrRuntime, "read()", pos}})
	}

	cbMaybeErr := func(v ValueThunk) {
		if err, ok := trampoline(v).(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to read()", pos})
		}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, errMsg(msg)))
		})
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// open
		file, err := os.OpenFile(string(filePath), os.O_RDONLY, 0o644)
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
				cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
					"type": ValueString("data"),
					"data": ValueString{},
				}))
			})
			return
		} else if err != nil {
			sendErr(fmt.Sprintf("error reading requested file in read(), %s", err.Error()))
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
				"type": ValueString("data"),
				"data": ValueString(buf[:count]),
			}))
		})
	}()

	return k(Null)
}

func inkWrite(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
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
		return k(ValueError{&Err{err, ErrRuntime, "write()", pos}})
	}

	cbMaybeErr := func(v ValueThunk) {
		if err, ok := trampoline(v).(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to write()", pos})
		}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, errMsg(msg)))
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
		file, err := os.OpenFile(string(filePath), flag, 0o644)
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
		if _, err := file.Write(buf); err != nil {
			sendErr(fmt.Sprintf("error writing to requested file in write(), %s", err.Error()))
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
				"type": ValueString("end"),
			}))
		})
	}()

	return k(Null)
}

func inkDelete(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var (
		filePath ValueString
		cb       ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &filePath),
		validateArgType(in, 1, &cb),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "delete()", pos}})
	}

	cbMaybeErr := func(v ValueThunk) {
		if err, ok := trampoline(v).(ValueError); ok {
			LogError(&Err{err.Err, ErrRuntime, "error in callback to delete()", pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// delete
		err := os.RemoveAll(string(filePath))
		if err != nil {
			ctx.ExecListener(func() {
				cbMaybeErr(evalInkFunction(ast, cb, pos, k, errMsg(
					fmt.Sprintf("error removing requested file in delete(), %s", err.Error()),
				)))
			})
			return
		}

		ctx.ExecListener(func() {
			cbMaybeErr(evalInkFunction(ast, cb, pos, k, ValueComposite{
				"type": ValueString("end"),
			}))
		})
	}()

	return k(Null)
}

// inkHTTPHandler fulfills the Handler interface for inkListen() to work
type inkHTTPHandler struct {
	ctx         *Context
	ast         *AST
	inkCallback ValueFunction
}

func (h inkHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.ctx
	pos := Pos{ctx.File, 0, 0} // TODO: pass position here and everywhere
	cb := h.inkCallback

	cbMaybeErr := func(v Value) ValueThunk {
		return func() Value {
			if err, ok := v.(ValueError); ok {
				LogError(&Err{err.Err, ErrRuntime, "error in callback to listen()", pos})
			}
			return nil
		}
	}

	// unmarshal request
	method := r.Method
	url := r.URL.String()

	headers := make(ValueComposite, len(r.Header))
	for key, values := range r.Header {
		headers[key] = ValueString(strings.Join(values, ","))
	}

	var body Value
	if r.ContentLength == 0 {
		body = ValueString{}
	} else {
		bodyBuf, err := io.ReadAll(r.Body)
		if err != nil {
			ctx.ExecListener(func() {
				_ = trampoline(evalInkFunction(h.ast, cb, pos, cbMaybeErr, errMsg(
					fmt.Sprintf("error reading request in listen(), %s", err.Error()),
				)))
			})
			return
		}
		body = ValueString(bodyBuf)
	}

	// construct request object to pass to Ink, and call handler
	responseEnded := false
	responses := make(chan Value, 1)
	// this is what Ink's callback calls to send a response
	endHandler := func(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
		if len(in) != 1 {
			LogError(&Err{nil, ErrRuntime, "end() callback to listen() must have one argument", pos})
		}
		if responseEnded {
			LogError(&Err{nil, ErrRuntime, "end() callback to listen() was called more than once", pos})
		}
		responseEnded = true
		responses <- in[0]

		return k(Null)
	}

	ctx.ExecListener(func() {
		_ = trampoline(evalInkFunction(h.ast, cb, pos, cbMaybeErr, ValueComposite{
			"type": ValueString("req"),
			"data": ValueComposite{
				"method":  ValueString(method),
				"url":     ValueString(url),
				"headers": headers,
				"body":    body,
			},
			"end": NativeFunctionValue{
				name: "end",
				exec: endHandler,
				ctx:  ctx,
			},
		}))
	})

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
			w.Header().Set(k, string(str))
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
	if _, err := w.Write(resBody); err != nil {
		ctx.ExecListener(func() {
			_ = trampoline(evalInkFunction(h.ast, cb, pos, cbMaybeErr, errMsg(
				fmt.Sprintf("error writing request body in listen(), %s", err.Error()),
			)))
		})
		return
	}
}

func inkListen(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var (
		host ValueString
		cb   ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &host),
		validateArgType(in, 1, &cb),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "listen()", pos}})
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			if err, ok := evalInkFunction(ast, cb, pos, k, errMsg(msg))().(ValueError); ok {
				LogError(&Err{err.Err, ErrRuntime, "error in callback to listen()", pos})
			}
		})
	}

	server := &http.Server{
		Addr: string(host),
		Handler: inkHTTPHandler{
			ctx:         ctx,
			inkCallback: cb,
			ast:         ast,
		},
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()
		fmt.Fprintf(os.Stderr, "listening on %s\n", string(server.Addr))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			sendErr(fmt.Sprintf("error starting http server in listen(), %s", err.Error()))
		}
	}()

	closer := func(ctx *Context, _ Pos, in []Value, k Cont) ValueThunk {
		// attempt graceful shutdown, concurrently, without
		// blocking Ink evaluation thread
		ctx.Engine.Listeners.Add(1)
		go func() {
			defer ctx.Engine.Listeners.Done()

			if err := server.Shutdown(context.Background()); err != nil {
				sendErr(fmt.Sprintf("error closing server in listen(), %s", err.Error()))
			}
		}()

		return k(Null)
	}

	return k(NativeFunctionValue{
		name: "close",
		exec: closer,
		ctx:  ctx,
	})
}

func inkReq(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var (
		data ValueComposite
		cb   ValueFunction
	)
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &data),
		validateArgType(in, 1, &cb),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "req()", pos}})
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// do not follow redirects
			return http.ErrUseLastResponse
		},
	}
	reqContext, reqCancel := context.WithCancel(context.Background())

	closer := func(_ *Context, _ Pos, _ []Value, k Cont) ValueThunk {
		reqCancel()
		return k(Null)
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			if err, ok := trampoline(evalInkFunction(ast, cb, pos, k, errMsg(msg))).(ValueError); ok {
				LogError(&Err{err.Err, ErrRuntime, "error in callback to req()", pos})
			}
		})
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// unmarshal request contents
		methodVal := cmp.Or(data["method"], Value(ValueString("GET")))
		urlVal := data["url"]
		headersVal := cmp.Or(data["headers"], Value(ValueComposite{}))
		bodyVal := cmp.Or(data["body"], Value(ValueString("")))

		// TODO: add query params

		reqMethod, okMethod := methodVal.(ValueString)
		reqURL, okURL := urlVal.(ValueString)
		reqHeaders, okHeaders := headersVal.(ValueComposite)
		reqBody, okBody := bodyVal.(ValueString)

		if !okMethod || !okURL || !okHeaders || !okBody {
			LogError(&Err{nil, ErrRuntime, fmt.Sprintf("request in req() is malformed, %s", data), pos})
			return
		}

		req, err := http.NewRequest(
			string(reqMethod),
			string(reqURL),
			bytes.NewReader(reqBody),
		)
		if err != nil {
			sendErr(fmt.Sprintf("error creating request in req(), %s", err.Error()))
			return
		}

		req = req.WithContext(reqContext)

		// construct headers
		// Content-Length is automatically set for us by Go
		req.Header.Set("User-Agent", "") // remove Go's default user agent header
		for k, v := range reqHeaders {
			if str, isStr := v.(ValueString); isStr {
				req.Header.Set(k, string(str))
			} else {
				LogError(&Err{nil, ErrRuntime, fmt.Sprintf("could not set request header, value %s was not a string", v), pos})
			}
		}

		// send request
		resp, err := client.Do(req)
		if err != nil {
			sendErr(fmt.Sprintf("error processing request in req(), %s", err.Error()))
			return
		}
		defer resp.Body.Close()

		resStatus := ValueNumber(resp.StatusCode)
		resHeaders := make(ValueComposite, len(resp.Header))
		for key, values := range resp.Header {
			resHeaders[key] = ValueString(strings.Join(values, ","))
		}

		var resBody Value
		if resp.ContentLength == 0 {
			resBody = ValueString{}
		} else {
			bodyBuf, err := io.ReadAll(resp.Body)
			if err != nil {
				sendErr(fmt.Sprintf("error reading response in req(), %s", err.Error()))
				return
			}
			resBody = ValueString(bodyBuf)
		}

		ctx.ExecListener(func() {
			_ = trampoline(evalInkFunction(ast, cb, pos, func(v Value) ValueThunk {
				err, ok := v.(ValueError)
				if ok {
					LogError(&Err{err.Err, ErrRuntime, "error in callback to req()", pos})
				}
				return func() Value { return Null }
			}, ValueComposite{
				"type": ValueString("resp"),
				"data": ValueComposite{
					"status":  resStatus,
					"headers": resHeaders,
					"body":    resBody,
				},
			}))
		})
	}()

	return k(NativeFunctionValue{
		name: "close",
		exec: closer,
		ctx:  ctx,
	})
}

func inkRand(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	return k(ValueNumber(rand.Float64()))
}

func inkUrand(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var bufLength ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &bufLength),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "urand()", pos}})
	}

	buf := make([]byte, int64(bufLength))
	if _, err := crand.Read(buf); err != nil {
		return k(Null)
	}

	return k(ValueString(buf))
}

func inkTime(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	unixSeconds := float64(time.Now().UnixNano()) / 1e9
	return k(ValueNumber(unixSeconds))
}

func inkWait(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var secs ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &secs),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "wait()", pos}})
	}

	// This is a bit tricky, since we don't want wait() to hold the evalLock
	// on the Context while we're waiting for the timeout, but do want to hold
	// the main goroutine from completing with sync.WaitGroup.
	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		time.Sleep(time.Duration(
			int64(float64(secs) * float64(time.Second)),
		))

		ctx.ExecListener(func() {
			_ = trampoline(evalInkFunction(ast, in[1], pos, func(v Value) ValueThunk {
				if err, ok := v.(ValueError); ok {
					LogError(err.Err)
				}
				return func() Value { return Null }
			}))
		})
	}()

	return k(Null)
}

func inkExec(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
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
		return k(ValueError{&Err{err, ErrRuntime, "exec()", pos}})
	}

	argsList := make([]string, len(args))
	for i, v := range args {
		argsList[i] = string(v)
	}

	cmd := exec.Command(string(path), argsList...)
	// cmdMutex locks control over reading and modifying child
	// process state, because both the Ink eval thread and exec
	// thread must read from/write to cmd.
	cmdMutex := sync.Mutex{} // TODO: remove as much mutexes as possible
	stdout := bytes.Buffer{}
	cmd.Stdin = strings.NewReader(string(stdin))
	cmd.Stdout = &stdout

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			_ = trampoline(evalInkFunction(ast, stdoutFn, pos, func(v Value) ValueThunk {
				if err, ok := v.(ValueError); ok {
					LogError(&Err{err.Err, ErrRuntime, "error in callback to exec()", pos})
				}
				return func() Value { return Null }
			}, errMsg(msg)))
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
				"type": ValueString("data"),
				"data": ValueString(output),
			}
			_ = trampoline(evalInkFunction(ast, stdoutFn, pos, func(v Value) ValueThunk {
				if err, ok := v.(ValueError); ok {
					LogError(&Err{err.Err, ErrRuntime, "error in callback to exec()", pos})
				}
				return func() Value { return Null }
			}, in))
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
	return k(NativeFunctionValue{
		name: "close",
		exec: func(_ *Context, pos Pos, _ []Value, k Cont) ValueThunk {
			// multiple calls to close() should be idempotent
			if closed {
				return k(Null)
			}

			neverRun <- true
			closed = true

			cmdMutex.Lock()
			if cmd.Process != nil || cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				if err := cmd.Process.Kill(); err != nil {
					return k(ValueError{&Err{nil, ErrRuntime, err.Error(), pos}})
				}
			}
			cmdMutex.Unlock()

			return k(Null)
		},
		ctx: ctx,
	})
}

func inkEnv(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	envp := os.Environ()

	envVars := make(ValueComposite, len(envp))
	for _, e := range envp {
		kv := strings.SplitN(e, "=", 2)
		envVars[kv[0]] = ValueString(kv[1])
	}
	return k(envVars)
}

func inkExit(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var code ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &code),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "exit()", pos}})
	}

	os.Exit(int(code))
	return nil // unreachable
}

func inkSin(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "sin()", pos}})
	}

	return k(ValueNumber(math.Sin(float64(inNum))))
}

func inkCos(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "cos()", pos}})
	}

	return k(ValueNumber(math.Cos(float64(inNum))))
}

func inkAsin(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
		validateCustom(inNum >= -1 && inNum <= 1, fmt.Sprintf("number must be in range [-1, 1], got %v", inNum)),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "asin()", pos}})
	}

	return k(ValueNumber(math.Asin(float64(inNum))))
}

func inkAcos(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var inNum ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
		validateCustom(inNum >= -1 && inNum <= 1, fmt.Sprintf("number must be in range [-1, 1], got %v", inNum)),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "acos()", pos}})
	}

	return k(ValueNumber(math.Acos(float64(inNum))))
}

func inkPow(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var base, exp ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 2),
		validateArgType(in, 0, &base),
		validateArgType(in, 1, &exp),
		validateCustom(base != 0 || exp != 0, "math error, pow(0, 0) is not defined"),
		validateCustom(base >= 0 || isInteger(exp), "math error, fractional power of negative number"),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "pow()", pos}})
	}

	return k(ValueNumber(math.Pow(float64(base), float64(exp))))
}

func inkLn(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var n ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &n),
		validateCustom(n > 0, fmt.Sprintf("cannot take natural logarithm of non-positive number %s", n.String())),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "ln()", pos}})
	}

	return k(ValueNumber(math.Log(float64(n))))
}

func inkFloor(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var n ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &n),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "floor()", pos}})
	}

	return k(ValueNumber(math.Trunc(float64(n))))
}

func inkString(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "string()", pos}})
	}

	s := func() string {
		switch v := in[0].(type) {
		case ValueError:
			return v.Error()
		case ValueString:
			return string(v)
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
	return k(ValueString(s))
}

func inkNumber(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "number()", pos}})
	}

	switch v := in[0].(type) {
	case ValueString:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return k(Null)
		}
		return k(ValueNumber(f))
	case ValueNumber:
		return k(v)
	case ValueBoolean:
		return k(ValueNumber(fun.IF[float64](bool(v), 1, 0)))
	default:
		return k(ValueNumber(0))
	}
}

func inkPoint(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var str ValueString
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &str),
		validateCustom(len(str) >= 1, "argument must be of length at least 1"),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "point()", pos}})
	}

	return k(ValueNumber(str[0]))
}

func inkChar(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	var cp ValueNumber
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &cp),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "char()", pos}})
	}

	return k(ValueString([]byte{byte(cp)}))
}

func inkType(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return k(ValueError{&Err{err, ErrAssert, "type()", pos}})
	}

	switch in[0].(type) {
	case ValueString:
		return k(ValueString("string"))
	case ValueNumber:
		return k(ValueString("number"))
	case ValueBoolean:
		return k(ValueString("boolean"))
	case ValueNull:
		return k(ValueString("()"))
	case ValueComposite:
		return k(ValueString("composite"))
	case ValueList:
		return k(ValueString("list"))
	case ValueFunction, NativeFunctionValue:
		return k(ValueString("function"))
	case ValueError:
		return k(ValueString("error"))
	default:
		return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("unknown type: %T", in[0]), pos}})
	}
}

func inkLen(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	if err := validate(pos,
		validateArgsLen(in, 1),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "len()", pos}})
	}

	switch v := in[0].(type) {
	case ValueComposite:
		return k(ValueNumber(len(v)))
	case ValueList:
		return k(ValueNumber(len(*v.xs)))
	case ValueString:
		// TODO: bytes count/rune count/grapheme clusters count?
		return k(ValueNumber(len(v)))
	default:
		return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("len() takes a string or composite value, but got %s", in[0]), pos}})
	}
}

func inkKeys(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	// var obj ValueComposite
	if err := validate(pos,
		validateArgsLen(in, 1),
		// validateArgType(in, 0, &obj),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "keys()", pos}})
	}

	switch obj := in[0].(type) {
	case ValueComposite:
		cv := make(ValueComposite, len(obj))
		i := 0
		for k := range obj {
			cv[strconv.Itoa(i)] = ValueString(k)
			i++
		}
		return k(cv)
	case ValueList:
		xs := make([]Value, len(*obj.xs))
		cv := ValueList{&xs}
		for i := range *obj.xs {
			(*cv.xs)[i] = ValueNumber(i)
		}
		return k(cv)
	default:
		return k(ValueError{&Err{nil, ErrRuntime, "keys()", pos}})
	}
}

func inkPar(ctx *Context, pos Pos, in []Value, k Cont) ValueThunk {
	ast := ctx.Engine.AST
	var funcs ValueComposite
	if err := validate(pos,
		validateArgsLen(in, 1),
		validateArgType(in, 0, &funcs),
	); err != nil {
		return k(ValueError{&Err{err, ErrRuntime, "par()", pos}})
	}

	// evalLock blocks file eval; temporary unlock it for the import to run.
	// Calling import() from within a running program is not supported, so we
	// don't really care if catastrophic things happen because of unlocked evalLock.
	ctx.Engine.mu.Unlock()
	defer ctx.Engine.mu.Lock()

	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		go func() {
			_ = trampoline(evalInkFunction(ast, f, pos, k))
			wg.Done()
		}()
	}
	wg.Wait()

	// TODO: composite of results
	return k(Null)
}
