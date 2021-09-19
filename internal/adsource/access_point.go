//
// @project GeniusRabbit rotator 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package adsource

// AccessPoint is the DSP source
type AccessPoint interface {
	// ID of source
	ID() uint64
}
