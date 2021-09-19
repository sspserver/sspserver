//
// @project Geniusrabbit::corelib 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package infostructs

func intO(v int) (vv *int) {
	if v != 0 {
		vv = new(int)
		*vv = v
	}
	return
}
