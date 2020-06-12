package registry

import (
    . "fmt"
    "reflect"
    "testing"

    "github.com/stretchr/testify/assert"
)

type Test struct {
    RegistryAPI
    Args []interface{}
    Expected interface{}
}
type Tests = []Test

type TestSuite struct {
    Tests
    Testing *testing.T
    resultType reflect.Type
    receiver reflect.Type
    assert func(t *testing.T, expected interface{}, actual interface{})
}


func NewTest(expected interface{}, args ...interface{}) Test {
    r := Test{
        Args: args,
        Expected: expected,
    }
    r.RegistryAPI = r
    return r
}

func TestSuiteWorks(t *testing.T) {

    tests := TestSuite{
        Tests: Tests{
            NewTest([]string{"3.1.00", "3.1.202"}),
            NewTest([]string{}),
        },
        resultType: reflect.ValueOf([]string{}).Type(),
        assert: func(t *testing.T, expected interface{}, actual interface{}) {
            Printf("[Expected] %v\n", expected)
            Printf("[Actual] %v\n", actual)
            assert.Contains(t, expected, actual)
        },
        Testing: t,
    }

    tests.Run(RegistryAPI.GetSDKs)
}

func createInArgs(test Test) []reflect.Value {
    var inArgs []reflect.Value
    inArgs = append(inArgs, reflect.ValueOf(test.RegistryAPI))

    for i := 0; i < len(test.Args); i++ {
        v := test.Args[i]
        value := reflect.ValueOf(v)
        inArgs = append(inArgs, value)
    }

    return inArgs
}

func (t TestSuite) Run(f interface{}) {
    for i, test := range t.Tests {
        Printf("Run %d Start\n", i)
        inArgs := createInArgs(test)

        callResult := reflect.ValueOf(f).Call(inArgs)
        result := callResult[0]

        expectedValue:= reflect.ValueOf(t.Tests[0].Expected)
        converted :=  result.Convert(expectedValue.Type()).Interface()

        assert.Contains(t.Testing, converted, t.Tests[0].Expected)
        t.assert(t.Testing, converted, t.Tests[0].Expected)
        Printf("Run %d End\n--------\n\n", i)
    }
}
