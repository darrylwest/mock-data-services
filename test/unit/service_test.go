//
// service
//
// @author darryl.west@ebay.com
// @created 2018-02-27 16:22:25
//

package unit

import (
	"fmt"
	"proxy"
	"testing"

	. "github.com/franela/goblin"
)

func TestService(t *testing.T) {
	g := Goblin(t)

	g.Describe("Service", func() {
		log := proxy.CreateLogger()
		log.SetLevel(4)
		cfg := proxy.NewDefaultConfig()

		g.It("should create a service struct", func() {
			service, err := proxy.NewService(cfg)
			g.Assert(err).Equal(nil)
			g.Assert(fmt.Sprintf("%T", service)).Equal("*proxy.Service")
		})
	})
}
