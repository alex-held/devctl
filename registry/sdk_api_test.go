package registry

import (
	. "fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	RegistryAPI
	ApiTestBase
	Client   TestRegistryApiClient
	Args     []interface{}
	Expected interface{}
}
type Tests = []Test

type TestSuite struct {
	Tests
	Testing    *testing.T
	resultType reflect.Type
	receiver   reflect.Type
	assert     func(t *testing.T, expected interface{}, actual interface{})
}

func NewTest(t *testing.T, expected interface{}, path string, body string, args ...interface{}) Test {
	testBase := ApiTestBase{
		Body:         body,
		ExpectedPath: path,
	}

	client := CreateTestRegistry(testBase, t)
	r := Test{
		RegistryAPI: client.client,
		ApiTestBase: ApiTestBase{
			Body:         body,
			ExpectedPath: path,
		},
		Args:     args,
		Expected: expected,
		Client:   client,
	}
	return r
}

//noinspection ALL
func TestSuiteWorks(t *testing.T) {

	tests := TestSuite{
		Tests: Tests{
			NewTest(t, []string{"dotnet"}, "/repos/alex-held/dev-env-registry/contents/sdk", "[{\"name\":\"dotnet\",\"path\":\"sdk/dotnet\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\"}}]"),
			NewTest(t, []string{"dotnet", "java"}, "/repos/alex-held/dev-env-registry/contents/sdk", "[{\"name\":\"dotnet\",\"path\":\"sdk/dotnet\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/dotnet?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc91\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/dotnet\"}},{\"name\":\"java\",\"path\":\"sdk/java\",\"sha\":\"859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"size\":0,\"url\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/java?ref=master\",\"html_url\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/java\",\"git_url\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"download_url\":null,\"type\":\"dir\",\"_links\":{\"self\":\"https://api.github.com/repos/alex-held/dev-env-registry/contents/sdk/java?ref=master\",\"git\":\"https://api.github.com/repos/alex-held/dev-env-registry/git/trees/859e4a060e287c06f777da09fbf8fe51dc4afc92\",\"html\":\"https://github.com/alex-held/dev-env-registry/tree/master/sdk/java\"}}]"),
		},
		assert: func(t *testing.T, expected interface{}, actual interface{}) {
			e := expected.([]string)
			a := actual.([]string)
			Printf("[Expected]\t\t%v\n", e)
			Printf("[Actual]\t\t%v\n", a)
			assert.Subset(t, a, e)
			assert.Subset(t, e, a)
			assert.Equal(t, a, e)
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
		i += 1
		Printf("-------- Run #%d --------\n", i)
		inArgs := createInArgs(test)

		callResult := reflect.ValueOf(f).Call(inArgs)
		result := callResult[0]

		expectedValue := reflect.ValueOf(test.Expected)
		converted := result.Convert(expectedValue.Type()).Interface()
		t.assert(t.Testing, test.Expected, converted)
		Printf("--------\n\n")
	}
}
