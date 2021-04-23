package matchers

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func TestMatcherOpensFile(t *testing.T) {
	RegisterFailHandlerWithT(t, func(m string, _ ...int) { t.Fatal(m) })

	type testCase struct {
		Fs              interface{}
		Setup           func(fs interface{}) error
		ExpectedSuccess bool
	}
	const filename = "/test"
	tests := []testCase{
		{
			Fs: afero.NewMemMapFs(),
			Setup: func(fs interface{}) error {
				_, err := fs.(afero.Fs).Create(filename)
				return err
			},
			ExpectedSuccess: true,
		},
		{
			Fs: afero.NewMemMapFs(),
			Setup: func(fs interface{}) error {
				_, err := fs.(afero.Fs).Create("does not assert this file name")
				return err
			},
			ExpectedSuccess: false,
		},
	}

	for _, tc := range tests {
		err := tc.Setup(tc.Fs)
		if err != nil {
			t.Fatal(err)
		}
		if tc.ExpectedSuccess {
			Expect(filename).Should(BeAnExistingFileFs(tc.Fs))
			continue
		}
		Expect(filename).Should(Not(BeAnExistingFileFs(tc.Fs)))

	}
}
