//
// @project geniusrabbit::rotator 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package models

func mapDef(m map[string]interface{}, key string, def interface{}) interface{} {
	if m == nil {
		return def
	}
	if v, _ := m[key]; nil != v {
		return v
	}
	return def
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
