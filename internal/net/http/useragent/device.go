//
// @project RevolvEads
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015
//

package useragent

type Device int

const (
  DEVICE_UNDEFINED Device = iota
  DEVICE_WEB
  DEVICE_MOBILE
  DEVICE_TABLET
  DEVICE_TV
)

var DeviceNames = []string{
  "Undefined",
  "PC",
  "Mobile",
  "Tablet",
  "TV",
}

func (d Device) String() string {
  if d < 0 || d >= Device(len(DeviceNames)) {
    return DeviceNames[0]
  }
  return DeviceNames[d]
}
