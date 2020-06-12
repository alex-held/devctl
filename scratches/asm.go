package main

import (
	. "fmt"
	"reflect"
)

type Test struct {
	TestInterface
	Args []interface{}
	Name string
}

type ExecFunc = func(TestInterface, string) (interface{}, error)

func (test Test) WithParam(msg string) (interface{}, error) {
	Println("Name:   " + test.Name)
	Println("WithParam(msg string)")
	Printf("[%d]: %v\n", 1, msg)
	return test, nil
}

func (test Test) WithMultiParam(b bool, msg string) (interface{}, error) {
	Println("WithMultiParam(i int, msg string)")
	Printf("[%d]: %v\n", 1, b)
	Printf("[%d]: %v\n", 2, msg)
	return test, nil
}

type TestSuite struct {
	Tests
	action reflect.Value
}

func (t TestSuite) Call(args ...interface{}) (interface{}, error) {
	var values []reflect.Value
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}
	result := t.action.Call(values)
	r0 := result[0]
	r1 := result[1]
	res := r0.Elem().Interface()
	er := r1.Elem()
	err := er.Interface().(error)
	return res, err
}

type Tests = []Test

func main() {
	tests := TestSuite{
		Tests: Tests{
			{Args: []interface{}{"held"}, Name: "Alex"},
			{Args: []interface{}{"papa"}, Name: "Mama"},
		},
	}
	//    tests.Run(func(ti TestInterface, args []interface{}) (interface{}, error) {return ti.WithParam(args[0].(string))})
	tests.Run(TestInterface.WithParam)

	/*  tests.Run(func(testInterface TestInterface, args ...interface{}) func() (interface{}, error) {
	    return func() (interface{}, error) { return testInterface.WithParam(args[0].(string)) }
	})*/

}
func get(in []reflect.Value, retType reflect.Type) (r []reflect.Value) {
	m := in[0]
	key := in[1]
	result := m.MapIndex(key)
	ok := false

	if result.IsValid() == false {
		r = []reflect.Value{reflect.Zero(retType), reflect.ValueOf(ok)}
		return r
	}
	resultval := result.Interface()
	if retType != reflect.TypeOf(resultval) {
		r = []reflect.Value{reflect.Zero(retType), reflect.ValueOf(ok)}
		return r
	}
	ok = true
	r = []reflect.Value{reflect.Zero(retType), reflect.ValueOf(ok)}
	return r
}

func stubFunc(typ reflect.Type, fn func(args []reflect.Value, retType reflect.Type) (results []reflect.Value)) reflect.Value {
	retType := typ.Out(0)
	f := reflect.MakeFunc(typ, func(args []reflect.Value) (results []reflect.Value) {
		return fn(args, retType)
	})
	return f
}

func (t TestSuite) makeAction(actionUnderTestP *interface{}) {

	fn := reflect.ValueOf(actionUnderTestP).Elem()
	fn.Set(stubFunc(fn.Type(), get))
	t.action = fn

}

func (t TestSuite) Run(execute interface{}) {
	executeValue := reflect.ValueOf(execute)
	Println(executeValue.Type().String())

	for _, test := range t.Tests {
		t.makeAction(&execute)
		Printf("RESULT: NAME=%v", test.Name)
	}
}

type TestInterface interface {
	WithParam(msg string) (interface{}, error)
	WithMultiParam(b bool, msg string) (interface{}, error)
}
