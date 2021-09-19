//
// @project GeniusRabbit::corelib 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package debugtool

import (
	"fmt"
	"runtime/debug"
)

// Trace debug stack
func Trace() {
	if recover() != nil {
		fmt.Println(string(debug.Stack()))
	}
}
