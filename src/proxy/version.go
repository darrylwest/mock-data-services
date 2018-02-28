// simple verion
//
// @author darryl.west@ebay.com
// @created 2018-02-27 16:22:25

package hub

import "fmt"

const (
	major = 18
	minor = 2
	patch = 27
)

// Version - return the version number as a single string
func Version() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
