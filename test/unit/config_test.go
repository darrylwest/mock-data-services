//
// config tests
//
// @author darryl.west <darwest@ebay.com>
// @created 2017-08-27 08:35:20
//

package unit

import (
	"fmt"
	"proxy"
	"testing"

	. "github.com/franela/goblin"
)

func TestConfig(t *testing.T) {
	g := Goblin(t)

	g.Describe("Config", func() {
		proxy.CreateLogger()

		g.It("should create a context struct with defaults set", func() {
			cfg := proxy.NewDefaultConfig()
			g.Assert(fmt.Sprintf("%T", cfg)).Equal("*proxy.Config")
			g.Assert(cfg.Port).Equal(3300)
			g.Assert(cfg.BufSize).Equal(64)
			g.Assert(cfg.LogLevel > 1).IsTrue()
			g.Assert(cfg.DbFilename).Equal("data/proxy.db")
		})

		g.It("should parse an empty command line and return default config", func() {
			cfg := proxy.ParseArgs()
			g.Assert(cfg != nil).IsTrue()
		})
	})
}
