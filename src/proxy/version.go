// simple verion
//
// @author darryl.west@ebay.com
// @created 2018-02-27 16:22:25

package proxy

import "fmt"

const (
	major = 18
	minor = 3
	patch = 1
	logo  = `
 __  __         _     ___       _          ___                  
|  \/  |___  __| |__ |   \ __ _| |_ __ _  | _ \_ _ _____ ___  _ 
| |\/| / _ \/ _| / / | |) / _' |  _/ _' | |  _/ '_/ _ \ \ / || |
|_|  |_\___/\__|_\_\ |___/\__,_|\__\__,_| |_| |_| \___/_\_\\_, |
                                                           |__/ 
`
)

// Version - return the version number as a single string
func Version() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
