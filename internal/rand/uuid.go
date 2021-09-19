//
// @project trafficstars corelib 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package rand

import (
	"github.com/twinj/uuid"
)

func init() {
	uuid.SwitchFormat(uuid.FormatCanonical)
}

// UUID generated
func UUID() string {
	return uuid.NewV4().String()
}
