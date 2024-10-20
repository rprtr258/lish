package internal

import (
	"bufio"
	"bytes"
	"context"
	crand "crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rprtr258/fun"
)

// NativeFunctionValue represents a function whose implementation is written
// in Go and built-into the runtime.
type NativeFunctionValue struct {
	name string
	exec func(*Context, []Value) (Value, *Err)
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
	for name, fn := range map[string]func(*Context, []Value) (Value, *Err){
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
	exec func(*Context, []Value) (Value, *Err),
) {
	ctx.Scope.Set(name, NativeFunctionValue{
		name,
		exec,
		ctx,
	})
}

// Create and return a standard error callback response with the given message
func errMsg(message string) ValueComposite {
	return ValueComposite{
		"type":    ValueString("error"),
		"message": ValueString(message),
	}
}

func validate(errs ...string) (string, bool) {
	for _, err := range errs {
		if err != "" {
			return err, true
		}
	}
	return "", false
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
		return fmt.Sprintf(
			"%d-th argument of must be a %T, but got %T",
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

	var args ValueComposite
	if err := validateArgType(in, i, &args); err != "" {
		return err
	}

	argsList := make([]T, len(args))
	for k, v := range args {
		i, err := strconv.ParseInt(string(k), 10, 64)
		if err != nil {
			return fmt.Sprintf("%d-th argument must be a list", i)
		}

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

func inkImport(ctx *Context, in []Value) (Value, *Err) {
	var givenPath ValueString
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &givenPath),
		validateCustom(len(givenPath) > 0, "arg must be path without the .ink suffix"),
	); ok {
		return nil, &Err{ErrRuntime, "import(): " + err, position{ctx.File, 0, 0}} // TODO: pass position here and everywhere
	}

	// imports via import() are assumed to be relative
	importPath := string(givenPath) + ".ink"
	if !filepath.IsAbs(importPath) {
		importPath = filepath.Join(ctx.WorkingDirectory, importPath)
	}

	// evalLock blocks file eval; temporary unlock it for the import to run.
	// Calling import() from within a running program is not supported, so we
	// don't really care if catastrophic things happen because of unlocked evalLock.
	ctx.Engine.mu.Unlock()
	defer ctx.Engine.mu.Lock()

	childCtx, ok := ctx.Engine.Contexts[importPath]
	if !ok {
		// The imported program runs in a "child context", a distinct context from
		// the importing program. The "child" term is a bit of a misnomer as Contexts
		// do not exist in a hierarchy, but conceptually makes sense here.
		childCtx = ctx.Engine.CreateContext()
		ctx.Engine.Contexts[importPath] = childCtx

		// Execution here follows updating ctx.Engine.Contexts
		// to behave correctly in the case where A imports B imports A again,
		// and still only import one instance of A.
		value, err := childCtx.ExecPath(importPath)
		if err != nil {
			return nil, &Err{ErrRuntime, fmt.Sprintf("error importing file %s: %s", importPath, err.Error()), position{ctx.File, 0, 0}} // TODO: pass position here and everywhere
		}

		ctx.Engine.values[importPath] = value
	}

	return ctx.Engine.values[importPath], nil
}

func inkArgs(*Context, []Value) (Value, *Err) {
	comp := make(ValueComposite, len(os.Args))
	for i, v := range os.Args {
		comp[strconv.Itoa(i)] = ValueString(v)
	}
	return comp, nil
}

func inkIn(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	if err, ok := validate(
		validateArgsLen(in, 1),
	); ok {
		return nil, &Err{ErrRuntime, "in(): " + err, pos}
	}

	cbErr := func(err *Err) {
		LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to in(), %s", err.Error()), pos})
	}

	ctx.ExecListener(func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				// also captures io.EOF
				break
			}

			rv, errEval := evalInkFunction(in[0], false, pos, ValueComposite{
				"type": ValueString("data"),
				"data": ValueString(str),
			})
			if errEval != nil {
				cbErr(errEval)
				return
			}

			if boolValue, isBool := rv.(ValueBoolean); isBool {
				if !boolValue {
					break
				}
			} else {
				LogError(&Err{ErrRuntime, fmt.Sprintf("callback to in() should return a boolean, but got %s", rv), pos})
				return
			}
		}

		_, err := evalInkFunction(in[0], false, pos, ValueComposite{
			"type": ValueString("end"),
		})
		if err != nil {
			cbErr(err)
			return
		}
	})

	return Null, nil
}

func inkOut(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var output ValueString
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &output),
	); ok {
		return nil, &Err{ErrRuntime, "out(): " + err, pos}
	}

	os.Stdout.Write([]byte(output))
	return Null, nil
}

func inkDir(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		dirPath ValueString
		cb      ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &dirPath),
		validateArgType(in, 1, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "dir(): " + err, pos}
	}

	cbMaybeErr := func(err *Err) {
		if err != nil {
			LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to dir(), %s", err.Error()), pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		fileInfos, err := ioutil.ReadDir(string(dirPath))
		if err != nil {
			ctx.ExecListener(func() {
				_, err := evalInkFunction(cb, false, pos, errMsg(
					fmt.Sprintf("error listing directory contents in dir(), %s", err.Error()),
				))
				cbMaybeErr(err)
			})
			return
		}

		fileList := make(ValueComposite, len(fileInfos))
		for i, fi := range fileInfos {
			fileList[strconv.Itoa(i)] = ValueComposite{
				"name": ValueString(fi.Name()),
				"len":  ValueNumber(fi.Size()),
				"dir":  ValueBoolean(fi.IsDir()),
				"mod":  ValueNumber(fi.ModTime().Unix()),
			}
		}

		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, ValueComposite{
				"type": ValueString("data"),
				"data": fileList,
			})
			cbMaybeErr(err)
		})
	}()

	return Null, nil
}

func inkMake(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		dirPath ValueString
		cb      ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &dirPath),
		validateArgType(in, 1, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "make(): " + err, pos}
	}

	cbMaybeErr := func(err *Err) {
		if err != nil {
			LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to make(), %s", err.Error()), pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		err := os.MkdirAll(string(dirPath), 0o755)
		if err != nil {
			ctx.ExecListener(func() {
				_, err := evalInkFunction(cb, false, pos, errMsg(
					fmt.Sprintf("error making a new directory in make(), %s", err.Error()),
				))
				cbMaybeErr(err)
			})
			return
		}

		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, ValueComposite{
				"type": ValueString("end"),
			})
			cbMaybeErr(err)
		})
	}()

	return Null, nil
}

func inkStat(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		statPath ValueString
		cb       ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &statPath),
		validateArgType(in, 1, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "stat(): " + err, pos}
	}

	cbMaybeErr := func(err *Err) {
		if err != nil {
			LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to stat(): %s", err.Error()), pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		fi, err := os.Stat(string(statPath))
		if err != nil {
			if os.IsNotExist(err) {
				ctx.ExecListener(func() {
					_, err := evalInkFunction(cb, false, pos, ValueComposite{
						"type": ValueString("data"),
						"data": Null,
					})
					cbMaybeErr(err)
				})
			} else {
				ctx.ExecListener(func() {
					_, err := evalInkFunction(cb, false, pos, errMsg(
						fmt.Sprintf("error getting file data in stat(), %s", err.Error()),
					))
					cbMaybeErr(err)
				})
			}
			return
		}

		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, ValueComposite{
				"type": ValueString("data"),
				"data": ValueComposite{
					"name": ValueString(fi.Name()),
					"len":  ValueNumber(fi.Size()),
					"dir":  ValueBoolean(fi.IsDir()),
					"mod":  ValueNumber(fi.ModTime().Unix()),
				},
			})
			cbMaybeErr(err)
		})
	}()

	return Null, nil
}

func inkRead(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		filePath ValueString
		offset   ValueNumber
		length   ValueNumber
		cb       ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 4),
		validateArgType(in, 0, &filePath),
		validateArgType(in, 1, &offset),
		validateArgType(in, 2, &length),
		validateArgType(in, 3, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "read(): " + err, pos}
	}

	cbMaybeErr := func(err *Err) {
		if err != nil {
			LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to read(): %s", err.Error()), pos})
		}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, errMsg(msg))
			cbMaybeErr(err)
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
				_, err := evalInkFunction(cb, false, pos, ValueComposite{
					"type": ValueString("data"),
					"data": ValueString{},
				})
				cbMaybeErr(err)
			})
			return
		} else if err != nil {
			sendErr(fmt.Sprintf("error reading requested file in read(), %s", err.Error()))
			return
		}

		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, ValueComposite{
				"type": ValueString("data"),
				"data": ValueString(buf[:count]),
			})
			cbMaybeErr(err)
		})
	}()

	return Null, nil
}

func inkWrite(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		filePath ValueString
		offset   ValueNumber
		buf      ValueString
		cb       ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 4),
		validateArgType(in, 0, &filePath),
		validateArgType(in, 1, &offset),
		validateArgType(in, 2, &buf),
		validateArgType(in, 3, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "write(): " + err, pos}
	}

	cbMaybeErr := func(err *Err) {
		if err != nil {
			LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to write(): %s", err.Error()), pos})
		}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, errMsg(msg))
			cbMaybeErr(err)
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
			_, err := file.Seek(ofs, 0) // 0 means relative to start of file
			if err != nil {
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
			_, err := evalInkFunction(cb, false, pos, ValueComposite{
				"type": ValueString("end"),
			})
			cbMaybeErr(err)
		})
	}()

	return Null, nil
}

func inkDelete(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		filePath ValueString
		cb       ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &filePath),
		validateArgType(in, 1, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "delete(): " + err, pos}
	}

	cbMaybeErr := func(err *Err) {
		if err != nil {
			LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to delete(): %s", err.Error()), pos})
		}
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// delete
		err := os.RemoveAll(string(filePath))
		if err != nil {
			ctx.ExecListener(func() {
				_, err := evalInkFunction(cb, false, pos, errMsg(
					fmt.Sprintf("error removing requested file in delete(), %s", err.Error()),
				))
				cbMaybeErr(err)
			})
			return
		}

		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, ValueComposite{
				"type": ValueString("end"),
			})
			cbMaybeErr(err)
		})
	}()

	return Null, nil
}

// inkHTTPHandler fulfills the Handler interface for inkListen() to work
type inkHTTPHandler struct {
	ctx         *Context
	inkCallback ValueFunction
}

func (h inkHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pos := position{h.ctx.File, 0, 0} // TODO: pass position here and everywhere
	ctx := h.ctx
	cb := h.inkCallback

	cbMaybeErr := func(err *Err) {
		if err != nil {
			LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to listen(): %s", err.Error()), pos})
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
				_, err := evalInkFunction(cb, false, pos, errMsg(
					fmt.Sprintf("error reading request in listen(), %s", err.Error()),
				))
				cbMaybeErr(err)
			})
			return
		}
		body = ValueString(bodyBuf)
	}

	// construct request object to pass to Ink, and call handler
	responseEnded := false
	responses := make(chan Value, 1)
	// this is what Ink's callback calls to send a response
	endHandler := func(ctx *Context, in []Value) (Value, *Err) {
		if len(in) != 1 {
			LogError(&Err{ErrRuntime, "end() callback to listen() must have one argument", pos})
		}
		if responseEnded {
			LogError(&Err{ErrRuntime, "end() callback to listen() was called more than once", pos})
		}
		responseEnded = true
		responses <- in[0]

		return Null, nil
	}

	ctx.ExecListener(func() {
		_, err := evalInkFunction(cb, false, pos, ValueComposite{
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
		})
		cbMaybeErr(err)
	})

	// validate response from Ink callback
	resp := <-responses
	rsp, isComposite := resp.(ValueComposite)
	if !isComposite {
		LogError(&Err{ErrRuntime, fmt.Sprintf("callback to listen() should return a response, got %s", resp), pos})
		return
	}

	// unmarshal response from the return value
	// response = {status, headers, body}
	statusVal, okStatus := rsp["status"]
	headersVal, okHeaders := rsp["headers"]
	bodyVal, okBody := rsp["body"]

	resStatus, okStatus := statusVal.(ValueNumber)
	resHeaders, okHeaders := headersVal.(ValueComposite)
	resBody, okBody := bodyVal.(ValueString)

	if !okStatus || !okHeaders || !okBody {
		LogError(&Err{ErrRuntime, fmt.Sprintf("callback to listen() returned malformed response, %s", rsp), pos})
		return
	}

	// write values to response
	// Content-Length is automatically set for us by Go
	for k, v := range resHeaders {
		if str, isStr := v.(ValueString); isStr {
			w.Header().Set(k, string(str))
		} else {
			LogError(&Err{ErrRuntime, fmt.Sprintf("could not set response header, value %s was not a string", v), pos})
			return
		}
	}

	code := int(resStatus)
	// guard against invalid HTTP codes, which cause Go panics.
	// https://golang.org/src/net/http/server.go
	if code < 100 || code > 599 {
		LogError(&Err{ErrRuntime, fmt.Sprintf("could not set response status code, code %d is not valid", code), pos})
		return
	}

	// status code write must follow all other header writes,
	// since it sends the status
	w.WriteHeader(int(resStatus))
	_, err := w.Write(resBody)
	if err != nil {
		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, errMsg(
				fmt.Sprintf("error writing request body in listen(), %s", err.Error()),
			))
			cbMaybeErr(err)
		})
		return
	}
}

func inkListen(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		host ValueString
		cb   ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &host),
		validateArgType(in, 1, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "listen(): " + err, pos}
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			_, err := evalInkFunction(cb, false, pos, errMsg(msg))
			if err != nil {
				LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to listen(), %s", err.Error()), pos})
			}
		})
	}

	server := &http.Server{
		Addr: string(host),
		Handler: inkHTTPHandler{
			ctx:         ctx,
			inkCallback: cb,
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

	closer := func(ctx *Context, in []Value) (Value, *Err) {
		// attempt graceful shutdown, concurrently, without
		// blocking Ink evaluation thread
		ctx.Engine.Listeners.Add(1)
		go func() {
			defer ctx.Engine.Listeners.Done()

			err := server.Shutdown(context.Background())
			if err != nil {
				sendErr(fmt.Sprintf("error closing server in listen(), %s", err.Error()))
			}
		}()

		return Null, nil
	}

	return NativeFunctionValue{
		name: "close",
		exec: closer,
		ctx:  ctx,
	}, nil
}

func inkReq(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		data ValueComposite
		cb   ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &data),
		validateArgType(in, 1, &cb),
	); ok {
		return nil, &Err{ErrRuntime, "req(): " + err, pos}
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// do not follow redirects
			return http.ErrUseLastResponse
		},
	}
	reqContext, reqCancel := context.WithCancel(context.Background())

	closer := func(ctx *Context, in []Value) (Value, *Err) {
		reqCancel()
		return Null, nil
	}

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			if _, err := evalInkFunction(cb, false, pos, errMsg(msg)); err != nil {
				LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to req(), %s", err.Error()), pos})
			}
		})
	}

	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// unmarshal request contents
		methodVal, okMethod := data["method"]
		urlVal, okURL := data["url"]
		headersVal, okHeaders := data["headers"]
		bodyVal, okBody := data["body"]
		// TODO: add query params

		if !okMethod {
			methodVal = ValueString("GET")
			okMethod = true
		}
		if !okHeaders {
			headersVal = ValueComposite{}
			okHeaders = true
		}
		if !okBody {
			bodyVal = ValueString("")
			okBody = true
		}

		reqMethod, okMethod := methodVal.(ValueString)
		reqURL, okURL := urlVal.(ValueString)
		reqHeaders, okHeaders := headersVal.(ValueComposite)
		reqBody, okBody := bodyVal.(ValueString)

		if !okMethod || !okURL || !okHeaders || !okBody {
			LogError(&Err{ErrRuntime, fmt.Sprintf("request in req() is malformed, %s", data), pos})
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
				LogError(&Err{ErrRuntime, fmt.Sprintf("could not set request header, value %s was not a string", v), pos})
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
			_, err := evalInkFunction(cb, false, pos, ValueComposite{
				"type": ValueString("resp"),
				"data": ValueComposite{
					"status":  resStatus,
					"headers": resHeaders,
					"body":    resBody,
				},
			})
			if err != nil {
				LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to req(), %s", err.Error()), pos})
			}
		})
	}()

	return NativeFunctionValue{
		name: "close",
		exec: closer,
		ctx:  ctx,
	}, nil
}

func inkRand(ctx *Context, in []Value) (Value, *Err) {
	return ValueNumber(rand.Float64()), nil
}

func inkUrand(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var bufLength ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &bufLength),
	); ok {
		return nil, &Err{ErrRuntime, "urand(): " + err, pos}
	}

	buf := make([]byte, int64(bufLength))
	if _, err := crand.Read(buf); err != nil {
		return Null, nil
	}

	return ValueString(buf), nil
}

func inkTime(ctx *Context, in []Value) (Value, *Err) {
	unixSeconds := float64(time.Now().UnixNano()) / 1e9
	return ValueNumber(unixSeconds), nil
}

func inkWait(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var secs ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &secs),
	); ok {
		return nil, &Err{ErrRuntime, "wait(): " + err, pos}
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
			if _, err := evalInkFunction(in[1], false, pos); err != nil {
				LogError(err)
			}
		})
	}()

	return Null, nil
}

func inkExec(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var (
		path     ValueString
		args     []ValueString
		stdin    ValueString
		stdoutFn ValueFunction
	)
	if err, ok := validate(
		validateArgsLen(in, 4),
		validateArgType(in, 0, &path),
		validateArgListOf(in, 1, &args),
		validateArgType(in, 2, &stdin),
		validateArgType(in, 3, &stdoutFn),
	); ok {
		return nil, &Err{ErrRuntime, "exec(): " + err, pos}
	}

	argsList := make([]string, len(args))
	for i, v := range args {
		argsList[i] = string(v)
	}

	cmd := exec.Command(string(path), argsList...)
	// cmdMutex locks control over reading and modifying child
	// process state, because both the Ink eval thread and exec
	// thread must read from/write to cmd.
	cmdMutex := sync.Mutex{}
	stdout := bytes.Buffer{}
	cmd.Stdin = strings.NewReader(string(stdin))
	cmd.Stdout = &stdout

	sendErr := func(msg string) {
		ctx.ExecListener(func() {
			_, err := evalInkFunction(stdoutFn, false, pos, errMsg(msg))
			if err != nil {
				LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to exec(), %s", err.Error()), pos})
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

		err = cmd.Wait()
		if err != nil {
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
			_, err := evalInkFunction(stdoutFn, false, pos, ValueComposite{
				"type": ValueString("data"),
				"data": ValueString(output),
			})
			if err != nil {
				LogError(&Err{ErrRuntime, fmt.Sprintf("error in callback to exec(), %s", err.Error()), pos})
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
		exec: func(ctx *Context, in []Value) (Value, *Err) {
			// multiple calls to close() should be idempotent
			if closed {
				return Null, nil
			}

			neverRun <- true
			closed = true

			cmdMutex.Lock()
			if cmd.Process != nil || cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				cmd.Process.Kill()
			}
			cmdMutex.Unlock()

			return Null, nil
		},
		ctx: ctx,
	}, nil
}

func inkEnv(ctx *Context, in []Value) (Value, *Err) {
	envp := os.Environ()

	envVars := make(ValueComposite, len(envp))
	for _, e := range envp {
		kv := strings.SplitN(e, "=", 2)
		envVars[kv[0]] = ValueString(kv[1])
	}
	return envVars, nil
}

func inkExit(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var code ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &code),
	); ok {
		return nil, &Err{ErrRuntime, "exit(): " + err, pos}
	}

	os.Exit(int(code))

	// unreachable
	return Null, nil
}

func inkSin(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var inNum ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
	); ok {
		return nil, &Err{ErrRuntime, "sin(): " + err, pos}
	}

	return ValueNumber(math.Sin(float64(inNum))), nil
}

func inkCos(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var inNum ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
	); ok {
		return nil, &Err{ErrRuntime, "cos(): " + err, pos}
	}

	return ValueNumber(math.Cos(float64(inNum))), nil
}

func inkAsin(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var inNum ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
		validateCustom(inNum >= -1 && inNum <= 1, fmt.Sprintf("number must be in range [-1, 1], got %v", inNum)),
	); ok {
		return nil, &Err{ErrRuntime, "asin(): " + err, pos}
	}

	return ValueNumber(math.Asin(float64(inNum))), nil
}

func inkAcos(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var inNum ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &inNum),
		validateCustom(inNum >= -1 && inNum <= 1, fmt.Sprintf("number must be in range [-1, 1], got %v", inNum)),
	); ok {
		return nil, &Err{ErrRuntime, "acos(): " + err, pos}
	}

	return ValueNumber(math.Acos(float64(inNum))), nil
}

func inkPow(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var base, exp ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 2),
		validateArgType(in, 0, &base),
		validateArgType(in, 1, &exp),
		validateCustom(base != 0 || exp != 0, "math error, pow(0, 0) is not defined"),
		validateCustom(base >= 0 || isInteger(exp), "math error, fractional power of negative number"),
	); ok {
		return nil, &Err{ErrRuntime, "pow(): " + err, pos}
	}

	return ValueNumber(math.Pow(float64(base), float64(exp))), nil
}

func inkLn(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var n ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &n),
		validateCustom(n > 0, fmt.Sprintf("cannot take natural logarithm of non-positive number %s", nvToS(n))),
	); ok {
		return nil, &Err{ErrRuntime, "ln(): " + err, pos}
	}

	return ValueNumber(math.Log(float64(n))), nil
}

func inkFloor(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var n ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &n),
	); ok {
		return nil, &Err{ErrRuntime, "floor(): " + err, pos}
	}

	return ValueNumber(math.Trunc(float64(n))), nil
}

func inkString(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	if err, ok := validate(
		validateArgsLen(in, 1),
	); ok {
		return nil, &Err{ErrRuntime, "string(): " + err, pos}
	}

	switch v := in[0].(type) {
	case ValueString:
		return v, nil
	case ValueNumber:
		return ValueString(nvToS(v)), nil
	case ValueBoolean:
		return ValueString(fun.IF(bool(v), "true", "false")), nil
	case ValueNull:
		return ValueString("()"), nil
	case ValueComposite:
		return ValueString(v.String()), nil
	case ValueFunction, NativeFunctionValue:
		return ValueString("(function)"), nil
	default:
		return ValueString(""), nil
	}
}

func inkNumber(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	if err, ok := validate(
		validateArgsLen(in, 1),
	); ok {
		return nil, &Err{ErrRuntime, "number(): " + err, pos}
	}

	switch v := in[0].(type) {
	case ValueString:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return Null, nil
		}
		return ValueNumber(f), nil
	case ValueNumber:
		return v, nil
	case ValueBoolean:
		var res float64
		if v {
			res = 1
		} else {
			res = 0
		}
		return ValueNumber(res), nil
	default:
		return ValueNumber(0), nil
	}
}

func inkPoint(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var str ValueString
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &str),
		validateCustom(len(str) >= 1, "argument must be of length at least 1"),
	); ok {
		return nil, &Err{ErrRuntime, "point(): " + err, pos}
	}

	return ValueNumber(str[0]), nil
}

func inkChar(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var cp ValueNumber
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &cp),
	); ok {
		return nil, &Err{ErrRuntime, "char(): " + err, pos}
	}

	return ValueString([]byte{byte(cp)}), nil
}

func inkType(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	if err, ok := validate(
		validateArgsLen(in, 1),
	); ok {
		return nil, &Err{ErrRuntime, "type(): " + err, pos}
	}

	switch in[0].(type) {
	case ValueString:
		return ValueString("string"), nil
	case ValueNumber:
		return ValueString("number"), nil
	case ValueBoolean:
		return ValueString("boolean"), nil
	case ValueNull:
		return ValueString("()"), nil
	case ValueComposite:
		return ValueString("composite"), nil
	case ValueFunction, NativeFunctionValue:
		return ValueString("function"), nil
	default:
		return ValueString(""), nil
	}
}

func inkLen(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	if err, ok := validate(
		validateArgsLen(in, 1),
	); ok {
		return nil, &Err{ErrRuntime, "len(): " + err, pos}
	}

	if list, isComposite := in[0].(ValueComposite); isComposite {
		return ValueNumber(len(list)), nil
	} else if str, isString := in[0].(ValueString); isString {
		return ValueNumber(len(str)), nil
	}

	return nil, &Err{ErrRuntime, fmt.Sprintf("len() takes a string or composite value, but got %s", in[0]), pos}
}

func inkKeys(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var obj ValueComposite
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &obj),
	); ok {
		return nil, &Err{ErrRuntime, "keys(): " + err, pos}
	}

	cv := make(ValueComposite, len(obj))
	i := 0
	for k := range obj {
		cv[strconv.Itoa(i)] = ValueString(k)
		i++
	}

	return cv, nil
}

func inkPar(ctx *Context, in []Value) (Value, *Err) {
	pos := position{ctx.File, 0, 0} // TODO: pass position here and everywhere
	var funcs ValueComposite
	if err, ok := validate(
		validateArgsLen(in, 1),
		validateArgType(in, 0, &funcs),
	); ok {
		return nil, &Err{ErrRuntime, "par(): " + err, pos}
	}

	// evalLock blocks file eval; temporary unlock it for the import to run.
	// Calling import() from within a running program is not supported, so we
	// don't really care if catastrophic things happen because of unlocked evalLock.
	ctx.Engine.mu.Unlock()
	defer ctx.Engine.mu.Lock()

	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		f := f
		go func() {
			evalInkFunction(f, false, pos)
			wg.Done()
		}()
	}
	wg.Wait()

	return Null, nil
}
