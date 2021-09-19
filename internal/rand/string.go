//
// @project trafficstars.com 2015
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015
//

package rand

import (
  "crypto/rand"
)

func Str(strSize int) string {
  return StrFromDict(strSize, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*_+-=~")
}

func StrUrlSafe(strSize int) string {
  return StrFromDict(strSize, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@$_")
}

func StrAlphaNum(strSize int) string {
  return StrFromDict(strSize, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
}

func StrAlpha(strSize int) string {
  return StrFromDict(strSize, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
}

func StrNum(strSize int) string {
  return StrFromDict(strSize, "0123456789")
}

func StrFromDict(strSize int, dict string) string {
  var bytes = make([]byte, strSize)
  rand.Read(bytes)
  for k, v := range bytes {
    bytes[k] = dict[v%byte(len(dict))]
  }
  return string(bytes)
}
