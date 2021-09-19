//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package platform

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	platforms = map[string]Factory{}
)

// Register platform factory
func Register(fact Factory) error {
	protocol := fact.Info().Protocol
	if _, ok := platforms[protocol]; ok {
		return fmt.Errorf("platform [%s] already exists", protocol)
	}
	log.WithField("module", "ssp:platform").Infof("Register platform: %s", protocol)
	platforms[protocol] = fact
	return nil
}

// Each all platforms
func Each(fn func(protocol string, platFact Factory)) {
	for protocol, platFact := range platforms {
		fn(protocol, platFact)
	}
}

// ByProtocol returns factory generator
func ByProtocol(protocol string) Factory {
	plat, _ := platforms[protocol]
	return plat
}
