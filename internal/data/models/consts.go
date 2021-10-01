//
// @project geniusrabbit::corelib 2016 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2018
//

package models

// ApproveStatus type
type ApproveStatus uint

// Status approve
const (
	StatusPending      ApproveStatus = 0
	StatusApproved     ApproveStatus = 1
	StatusRejected     ApproveStatus = 2
	StatusPendingName                = `pending`
	StatusApprovedName               = `approved`
	StatusRejectedName               = `rejected`
)

// Name of the status
func (st ApproveStatus) Name() string {
	switch st {
	case StatusApproved:
		return StatusApprovedName
	case StatusRejected:
		return StatusRejectedName
	}
	return StatusPendingName
}

// ApproveNameToStatus name to const
func ApproveNameToStatus(name string) ApproveStatus {
	switch name {
	case StatusApprovedName:
		return StatusApproved
	case StatusRejectedName:
		return StatusRejected
	}
	return StatusPending
}

// Status active
const (
	StatusPause  = 0
	StatusActive = 1
)

// Status private
const (
	StatusPublic  = 0
	StatusPrivate = 1
)

// ActiveStatusName from const
func ActiveStatusName(status uint) string {
	if status == StatusActive {
		return `active`
	}
	return `pause`
}
