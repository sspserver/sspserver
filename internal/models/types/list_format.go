//
// @Company GeniusRabbit::rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package types

// ListFormat array
type ListFormat []*Format

// HasSuitsCompare in the list
func (l ListFormat) HasSuitsCompare(format *Format) bool {
	for _, f := range l {
		if format.SuitsCompare(f) == 0 {
			return true
		}
	}
	return false
}
