package ascripts

import (
	"bytes"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"strings"
	"sync"
)

// CompilerGolangYaegi is a struct that holds an interpreter instance for executing Go code.
type CompilerGolangYaegi struct {
	SandBox *interp.Interpreter // The Yaegi interpreter instance.
	BIn     *bytes.Buffer       // Buffer for interpreter's standard input.
	BOut    *bytes.Buffer       // Buffer for interpreter's standard output and error.
}

// RunWithExports compiles and runs the Go code with the provided exports.
func (c *CompilerGolangYaegi) RunWithExports(body string, exports ...interp.Exports) (interface{}, error) {
	// Check if the body is empty after trimming whitespace.
	if strings.TrimSpace(body) == "" {
		return nil, fmt.Errorf("body of code to compile is empty")
	}

	// Compile the Go code using Yaegi with the provided exports.
	sandbox, bOut, bIn, err := CompileGolangYaegi(body, exports...)
	if err != nil {
		return nil, err
	}

	// Store the interpreter and buffers in the struct.
	c.SandBox = sandbox
	c.BOut = bOut
	c.BIn = bIn

	return sandbox, nil
}

// Run compiles and runs the Go code with the provided parameters.
func (c *CompilerGolangYaegi) Run(body string, params ...interface{}) (interface{}, error) {
	var exports []interp.Exports
	// Convert params to a slice of interp.Exports.
	if len(params) > 0 {
		for _, param := range params {
			pmap, ok := param.(interp.Exports)
			if !ok {
				return nil, fmt.Errorf("unable to cast param as type interp.Exports")
			}
			exports = append(exports, pmap)
		}
	}
	return c.RunWithExports(body, exports...)
}

// Render executes the Go code and retrieves the rendered output.
func (c *CompilerGolangYaegi) Render(body string, params ...interface{}) (string, error) {
	// Run the Go code.
	ro, err := c.Run(body, params...)
	if err != nil {
		return "", err
	}

	// Cast the result to an interpreter.
	embed, err := CastToGolangYaegi(ro)
	if err != nil {
		return "", err
	}

	// Evaluate the script to get the rendered output.
	res, err := embed.Eval("ascript.GetRendered")
	if err != nil {
		return "", fmt.Errorf("failed ascript.GetRendered: %v", err)
	}

	// Assert the result is a function and call it to get the string output.
	fn := res.Interface().(func() string)
	return fn(), nil
}

// CompileGolangYaegi compiles the Go code using Yaegi and returns the interpreter and buffers.
func CompileGolangYaegi(src string, params ...interp.Exports) (*interp.Interpreter, *bytes.Buffer, *bytes.Buffer, error) {
	bOut := &bytes.Buffer{}
	bIn := &bytes.Buffer{}

	// Create a new Yaegi interpreter with custom options.
	embed := interp.New(interp.Options{
		GoPath: "./_pkg",
		Stdin:  bIn,
		Stdout: bOut,
		Stderr: bOut,
	})

	// Use the predefined sandbox symbols.
	if err := embed.Use(GOSYMBOLSSANDBOX()); err != nil {
		return nil, bOut, bIn, fmt.Errorf("failed to use GoSymbolsSandbox: %v", err)
	}

	// Add any additional exports provided as parameters.
	for _, param := range params {
		if err := embed.Use(param); err != nil {
			return nil, bOut, bIn, fmt.Errorf("failed to use custom params: %v", err)
		}
	}

	// Evaluate the source code.
	_, err := embed.Eval(src)
	if err != nil {
		return nil, bOut, bIn, fmt.Errorf("failed to eval embedded Go script: %v", err)
	}

	return embed, bOut, bIn, nil
}

// CastToGolangYaegi casts the target to an interp.Interpreter.
func CastToGolangYaegi(target interface{}) (*interp.Interpreter, error) {
	if target == nil {
		return nil, fmt.Errorf("target is nil")
	}
	embed, ok := target.(*interp.Interpreter)
	if !ok {
		return nil, fmt.Errorf("failed to cast to *interp.Interpreter")
	}
	return embed, nil
}

// goSymbolsSandbox holds the allowed symbols for the Yaegi interpreter.
var goSymbolsSandbox interp.Exports

// once is used to ensure the sandbox is initialized only once.
var once sync.Once

// GOSYMBOLSSANDBOX initializes the GoSymbolsSandbox with allowed standard library symbols and returns the map.
func GOSYMBOLSSANDBOX() interp.Exports {
	once.Do(func() {
		useLib := []string{
			"net/netip/netip",
			"time/time",
			"strings/strings",
			"regexp/regexp",
			"encoding/json/json",
			"sort/sort",
			"math/math",
			"text/template/parse/parse",
			"bytes/bytes",
			"errors/errors",
			"strconv/strconv",
			"fmt/fmt",
			"image/color/color",
		}

		goSymbolsSandbox = make(interp.Exports)
		// Populate the sandbox with symbols from the allowed libraries.
		for key, val := range stdlib.Symbols {
			for _, kMatch := range useLib {
				if key == kMatch {
					goSymbolsSandbox[key] = val
				}
			}
		}
	})
	return goSymbolsSandbox
}
