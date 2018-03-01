//
// client tests
//
// @author darryl.west <darwest@ebay.com>
// @created 2017-08-28 08:35:20
//

package unit

import (
	"fmt"
	"proxy"
	"testing"
	"time"

	. "github.com/franela/goblin"
)

func TestClient(t *testing.T) {
	g := Goblin(t)

	g.Describe("Client", func() {
		now := time.Now()
		proxy.CreateLogger()
		cfg := proxy.NewDefaultConfig()

		g.It("should create a context struct with defaults set", func() {
			client := proxy.NewClient(cfg)

			g.Assert(fmt.Sprintf("%T", client)).Equal("*proxy.Client")
			g.Assert(len(client.GetID())).Equal(26)
			g.Assert(client.GetCreatedAt().After(now)).IsTrue()
		})

		g.It("should parse a request content length", func() {
			sz := 5543
			line := fmt.Sprintf("Content-Length: %d\r", sz)

			client := proxy.NewClient(cfg)

			val, err := client.ParseContentLength(line)
			g.Assert(err).Equal(nil)
			g.Assert(val).Equal(sz)
		})

		g.It("should handle a mock request")
	})
}
