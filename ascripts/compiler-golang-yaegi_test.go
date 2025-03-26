package ascripts

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/traefik/yaegi/interp"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

// IMyObj is an interface that defines a method to get a message.
type IMyObj interface {
	GetMessage() string
}

// MyObj is a struct that implements IMyObj.
type MyObj struct {
	Message string
}

// GetMessage returns the message from MyObj.
func (obj *MyObj) GetMessage() string {
	return obj.Message
}

// scriptBodyGolang is a Go script to be used in the test.
const scriptBodyGolang = `package foo
import (
	"ascripts"
	"fmt"
)

func GetRendered(obj *ascripts.MyObj) string {
	fmt.Printf("MyObj-Buffer-Out-Line-1\n")
	fmt.Printf("MyObj-Buffer-Out-Line-2")
	return obj.Message + "-MyObj"
}

func GetIRendered(obj ascripts.IMyObj) string {
	return obj.GetMessage() + "-MyObj"
}`

// TestYaegi_GolangCompilerRun tests the compilation and execution of Go code using Yaegi.
func TestYaegi_GolangCompilerRun(t *testing.T) {
	// Create a Scriptable with Go code.
	script := &Scriptable{
		Type: SCRIPTTYPE_GO,
		Body: scriptBodyGolang,
	}

	// Prepare the export map with types to be used in the script.
	pmap := interp.Exports{}
	pmap["ascripts/ascripts"] = map[string]reflect.Value{
		"MyObj":  reflect.ValueOf((*MyObj)(nil)),
		"IMyObj": reflect.ValueOf((*IMyObj)(nil)),
	}

	// Compile the script with the provided export map.
	rawcompiled, err := script.Compile(pmap)
	assert.NoError(t, err, "Compilation should not return an error")

	// Cast the compiled result to a Yaegi interpreter.
	embed, err := CastToGolangYaegi(rawcompiled)
	assert.NoError(t, err, "Casting to Yaegi interpreter should not return an error")

	// Evaluate the GetRendered function from the script.
	res, err := embed.Eval("foo.GetRendered")
	assert.NoError(t, err, "Evaluating GetRendered should not return an error")

	// Assert that the function returns the expected result.
	fn1 := res.Interface().(func(*MyObj) string)
	obj := &MyObj{Message: "test"}
	assert.Equal(t, "test-MyObj", fn1(obj), "The function should return the expected message")

	// Evaluate the GetIRendered function from the script.
	res2, err := embed.Eval("foo.GetIRendered")
	assert.NoError(t, err, "Evaluating GetIRendered should not return an error")

	// Assert that the function returns the expected result.
	fn2 := res2.Interface().(func(IMyObj) string)
	assert.Equal(t, "test-MyObj", fn2(obj), "The function should return the expected message")

	// Assert that the compiler is of the correct type and the output buffer contains the expected output.
	goc, ok := script.GetCompiler().(*CompilerGolangYaegi)
	assert.True(t, ok, "Compiler should be of type CompilerGolangYaegi")
	assert.Equal(t, "MyObj-Buffer-Out-Line-1\nMyObj-Buffer-Out-Line-2", goc.BOut.String(), "Output buffer should contain the expected output")
}

var testFilesystem = fstest.MapFS{
	"main.go": &fstest.MapFile{
		Data: []byte(`package main

import (
	"foo/bar"
	"./localfoo"
)

func main() {
	bar.PrintSomething()
	localfoo.PrintSomethingElse()
}
`),
	},
	"_pkg/src/foo/bar/bar.go": &fstest.MapFile{
		Data: []byte(`package bar

import (
	"fmt"
)

func PrintSomething() {
	fmt.Println("I am a virtual filesystem printing something from _pkg/src/foo/bar/bar.go!")
}
`),
	},
	"localfoo/foo.go": &fstest.MapFile{
		Data: []byte(`package localfoo

import (
	"fmt"
)

func PrintSomethingElse() {
	fmt.Println("I am virtual filesystem printing else from localfoo/foo.go!")
}
`),
	},
}

func TestFilesystemMapFS(t *testing.T) {
	i := interp.New(interp.Options{
		GoPath:               "./_pkg",
		SourcecodeFilesystem: testFilesystem,
	})
	if err := i.Use(GOSYMBOLSSANDBOX()); err != nil {
		t.Fatal(err)
	}
	_, err := i.EvalPath(`main.go`)
	if err != nil {
		t.Fatal(err)
	}
}

var testFilesystemEmbedStruct_Fail = fstest.MapFS{
	"bar.go": &fstest.MapFile{
		Data: []byte(`package bar

import (
    "fmt"
)

func NewFoo() *Foo {
	return &Foo{A: "test"}
}

type Foo struct {
	A string
}

func (f *Foo) PrintSomething() {
	fmt.Printf("I am printing '%s' from inside Foo.PrintSomething!\n", f.A)
}

func PrintSomething(myVal string) {
	fmt.Printf("I am printing '%s'\n", myVal)
}
`),
	},
}

// IFoo is an interface that defines a method to print something.
type IFoo interface {
	PrintSomething()
}

// TestFilesystemMapFS_EmbedStruct_Fail tests the failure of embedding a struct with Yaegi.
func TestFilesystemMapFS_EmbedStruct_Fail(t *testing.T) {
	// Initialize a new Yaegi interpreter with the provided options.
	i := interp.New(interp.Options{
		GoPath:               "./_pkg",
		SourcecodeFilesystem: testFilesystemEmbedStruct_Fail,
	})

	// Use the Go symbols sandbox in the interpreter.
	if err := i.Use(GOSYMBOLSSANDBOX()); err != nil {
		t.Fatal(err)
	}

	// Evaluate the 'bar.go' file which should contain the definition of Foo.
	_, err := i.EvalPath(`bar.go`)
	assert.NoError(t, err, "EvalPath should not return an error for 'bar.go'")

	// Attempt to create a new instance of Foo using the NewFoo function.
	val, err := i.Eval(`bar.NewFoo`)
	assert.NoError(t, err, "Eval should not return an error for 'bar.NewFoo'")

	// Call the NewFoo function to get the object values.
	objVals := val.Call(nil)

	// Assert that the returned object cannot be cast to the IFoo interface.
	_, ok := objVals[0].Interface().(IFoo)
	assert.False(t, ok, "objVals[0] should not be castable to IFoo")

	// Assert that the element of the returned object cannot be cast to the IFoo interface.
	_, ok = objVals[0].Elem().Interface().(IFoo)
	assert.False(t, ok, "objVals[0].Elem() should not be castable to IFoo")

	// Evaluate the PrintSomething function and execute it.
	val2, err := i.Eval(`bar.PrintSomething`)
	assert.NoError(t, err, "Eval should not return an error for 'bar.PrintSomething'")

	// Assert that the function executes without error.
	fn2 := val2.Interface().(func(string))
	fn2("test")
}

var testFilesystemEmbedStruct_XYZ = fstest.MapFS{
	"bar.go": &fstest.MapFile{
		Data: []byte(`package bar

import (
    "fmt"
	"xyz"
)

func PrintSomething(myVal string) {
	fmt.Printf("I am printing '%s'\n", myVal)
}

func PrintFoo(d *xyz.Data) string {
	d.ChangeMe = "I have mutated!"
	return d.Message + "-Foo"
}
`),
	},
}

type Data struct {
	Message  string
	ChangeMe string
}

func (d *Data) PrintChangeMe() {
	fmt.Printf(d.printChangeMe())
}

func (d *Data) PrintChangeMeExternal() string {
	return d.printChangeMe()
}

func (d *Data) printChangeMe() string {
	return fmt.Sprintf("'%s' from INSIDE the embedded script!\n", d.ChangeMe)
}

// TestFilesystemMapFS_EmbedStruct_XYZ tests the embedding of a custom struct in the Yaegi interpreter.
func TestFilesystemMapFS_EmbedStruct_XYZ(t *testing.T) {
	// Initialize a new Yaegi interpreter with the provided options.
	i := interp.New(interp.Options{
		GoPath:               "./_pkg",
		SourcecodeFilesystem: testFilesystemEmbedStruct_XYZ,
	})

	// Use the Go symbols sandbox in the interpreter.
	if err := i.Use(GOSYMBOLSSANDBOX()); err != nil {
		t.Fatal(err)
	}

	// Create a custom library containing the Data struct to be available for import.
	custom := interp.Exports{}
	custom["xyz/xyz"] = map[string]reflect.Value{
		"Data": reflect.ValueOf((*Data)(nil)), // Register the Data type.
	}

	// Use the custom library in the interpreter.
	if err := i.Use(custom); err != nil {
		t.Fatal(err)
	}

	// Evaluate the 'bar.go' file which should now have access to the Data struct.
	_, err := i.EvalPath(`bar.go`)
	assert.NoError(t, err, "EvalPath should not return an error for 'bar.go'")

	// Evaluate the PrintSomething function and execute it.
	val2, err := i.Eval(`bar.PrintSomething`)
	assert.NoError(t, err, "Eval should not return an error for 'bar.PrintSomething'")
	fn2 := val2.Interface().(func(string))
	fn2("test")

	// Evaluate the PrintFoo function, which takes a *Data parameter, and execute it.
	valPrintFoo, err := i.Eval("bar.PrintFoo")
	assert.NoError(t, err, "Eval should not return an error for 'bar.PrintFoo'")
	d := Data{Message: "Kung"}
	fnFoo := valPrintFoo.Interface().(func(*Data) string)
	results := fnFoo(&d)
	assert.Equal(t, "Kung-Foo", results, "PrintFoo should return the expected result")
	// This should print the changed message from the embedded script.
	assert.Equal(t, `'I have mutated!' from INSIDE the embedded script!`, strings.TrimSpace(d.PrintChangeMeExternal()))
}

// TestYaegi_FilesystemMapFS tests the virtual filesystem provided by fstest.MapFS.
func TestYaegi_FilesystemMapFS(t *testing.T) {
	// Initialize a new Yaegi interpreter with the virtual filesystem.
	i := interp.New(interp.Options{
		GoPath:               "./_pkg",
		SourcecodeFilesystem: testFilesystem,
	})

	// Use the predefined Go symbols in the interpreter.
	if err := i.Use(GOSYMBOLSSANDBOX()); err != nil {
		t.Fatal(err)
	}

	// Evaluate the 'main.go' file in the virtual filesystem.
	_, err := i.EvalPath(`main.go`)
	assert.NoError(t, err, "EvalPath should not return an error for 'main.go'")
}

var testFilesystemShouldFail = fstest.MapFS{
	"main.go": &fstest.MapFile{
		Data: []byte(`package main

import (
	"foo/bar"
	"./localfoo"
    "fmt"
    "os"
)

func main() {

	// test print args
    fmt.Println(len(os.Args), os.Args)

	// test read file outside of sandbox
	b, err := os.ReadFile("/etc/hosts")
    if err != nil {
        fmt.Println(err)
    } else {
    	fmt.Println(string(b))
	}

	bar.PrintSomething()
	localfoo.PrintSomethingElse()
}
`),
	},
	"_pkg/src/foo/bar/bar.go": &fstest.MapFile{
		Data: []byte(`package bar

import (
	"fmt"
)

func PrintSomething() {
	fmt.Println("I am a virtual filesystem printing something from _pkg/src/foo/bar/bar.go!")
}
`),
	},
	"localfoo/foo.go": &fstest.MapFile{
		Data: []byte(`package localfoo

import (
	"fmt"
)

func PrintSomethingElse() {
	fmt.Println("I am virtual filesystem printing else from localfoo/foo.go!")
}
`),
	},
}

// TestYaegi_FilesystemMapFS_ShouldFail tests the virtual filesystem with operations that should fail.
func TestYaegi_FilesystemMapFS_ShouldFail(t *testing.T) {
	// Initialize a new Yaegi interpreter with the virtual filesystem that should fail.
	i := interp.New(interp.Options{
		GoPath:               "./_pkg",
		SourcecodeFilesystem: testFilesystemShouldFail,
	})

	// Use the predefined Go symbols in the interpreter.
	if err := i.Use(GOSYMBOLSSANDBOX()); err != nil {
		t.Fatal(err)
	}

	// Evaluate the 'main.go' file in the virtual filesystem that should fail.
	_, err := i.EvalPath(`main.go`)
	assert.Error(t, err, "EvalPath should return an error for 'main.go' due to restricted operations")
}

// TestYaegi_RunWithExports checks if RunWithExports compiles and runs Go code correctly.
func TestYaegi_RunWithExports(t *testing.T) {
	c := &CompilerGolangYaegi{}
	body := `package main; func main() { println("Hello, world!") }`

	// Run the code with no additional exports.
	_, err := c.RunWithExports(body)
	assert.NoError(t, err, "RunWithExports should not return an error")

	// Check if the output buffer contains the expected output.
	expectedOutput := "Hello, world!\n"
	assert.Equal(t, expectedOutput, c.BOut.String(), "Output buffer should contain the expected output")
}

// TestYaegi_Run checks if Run compiles and runs Go code correctly with parameters.
func TestYaegi_Run(t *testing.T) {
	c := &CompilerGolangYaegi{}
	body := `package main; func main() { println("Hello, world!") }`

	// Run the code with no parameters.
	_, err := c.Run(body)
	assert.NoError(t, err, "Run should not return an error")

	// Check if the output buffer contains the expected output.
	expectedOutput := "Hello, world!\n"
	assert.Equal(t, expectedOutput, c.BOut.String(), "Output buffer should contain the expected output")
}

// TestCompileGolangYaegi checks if CompileGolangYaegi compiles Go code correctly.
func TestCompileGolangYaegi(t *testing.T) {
	src := `package main; func main() { println("Hello, world!") }`

	// Compile the source code.
	_, bOut, _, err := CompileGolangYaegi(src)
	assert.NoError(t, err, "CompileGolangYaegi should not return an error")

	// Check if the output buffer contains the expected output.
	expectedOutput := "Hello, world!\n"
	assert.Equal(t, expectedOutput, bOut.String(), "Output buffer should contain the expected output")
}

// TestCastToGolangYaegi checks if CastToGolangYaegi casts the target to an interp.Interpreter.
func TestCastToGolangYaegi(t *testing.T) {
	embed := interp.New(interp.Options{})
	target, err := CastToGolangYaegi(embed)
	assert.NoError(t, err, "CastToGolangYaegi should not return an error")
	assert.NotNil(t, target, "CastToGolangYaegi should return a non-nil interpreter")
}
