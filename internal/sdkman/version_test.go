package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/alex-held/devctl/internal/system"
	"github.com/alex-held/devctl/internal/testutils"
)

func TestVersionService_All(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("VersionService", func() {
		g.Describe("All", func() {
			var client *Client
			var _ *logrus.Logger
			var mux *http.ServeMux
			var _ bytes.Buffer
			var teardown testutils.Teardown
			var ctx context.Context

			g.JustBeforeEach(func() {
				client, _, mux, _, teardown = setup()
				ctx = context.Background()
			})

			g.AfterEach(func() {
				teardown()
			})

			g.It("WHEN listing all versions", func() {
				expected := strings.Split("2.11.0,2.11.1,2.11.11,2.11.12,2.11.2,2.11.3,2.11.4,2.11.5,2.11.6,2.11.7,2.11.8,2.12.0,2.12.1,2.12.10,2.12.11,2.12.12,2.12.13,2.12.2,2.12.3,2.12.4,2.12.5,2.12.6,2.12.7,2.12.8,2.12.9,2.13.0,2.13.1,2.13.2,2.13.3,2.13.4,3.0.0-M1,3.0.0-M2,3.0.0-M3", ",") //nolint: lll

				mux.HandleFunc("/candidates/scala/darwinx64/versions/all", func(w http.ResponseWriter, request *http.Request) {
					_, _ = w.Write([]byte("2.11.0,2.11.1,2.11.11,2.11.12,2.11.2,2.11.3,2.11.4,2.11.5,2.11.6,2.11.7,2.11.8,2.12.0,2.12.1,2.12.10,2.12.11,2.12.12,2.12.13,2.12.2,2.12.3,2.12.4,2.12.5,2.12.6,2.12.7,2.12.8,2.12.9,2.13.0,2.13.1,2.13.2,2.13.3,2.13.4,3.0.0-M1,3.0.0-M2,3.0.0-M3")) //nolint: lll
				})

				all, err := client.Version.All(ctx, "scala", system.MacOsx64)
				fmt.Println(all)
				Expect(err).To(BeNil())
				Expect(all).To(Equal(expected))
			})
		})
	})
}
