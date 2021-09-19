//
// @project GeniusRabbit rotator 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package models

// User set of flags
const (
	UserFlagCanHaveNegativeBalance = 1 << iota
)

// User model
type User struct {
	ID      uint64
	Balance int64
	Flags   uint64 // CanHaveNegativeBalance, IsSuperuser, Trusted
}
