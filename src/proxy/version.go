// simple verion
//
// @author darryl.west@ebay.com
// @created 2017-07-20 09:59:37

package hub

import "fmt"

const (
	major = 18
	minor = 2
	patch = 22
)

// Version - return the version number as a single string
func Version() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
