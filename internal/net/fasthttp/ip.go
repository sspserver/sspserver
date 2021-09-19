//
// @project geniusrabbit::corelib 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package fasthttp

import (
	"bytes"
	"net"
	"strings"

	"github.com/valyala/fasthttp"
)

//ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) <= 0 {
		return true
	}
	return false
}

// TODO: Use masks instead ranges
var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

// IPAdressByRequest conext
func IPAdressByRequest(ctx *fasthttp.RequestCtx) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		if ip := getIPAdress(string(ctx.Request.Header.Peek(h))); ip != "" {
			return ip
		}
	}
	return ""
}

// IPAdressByRequestCF context detection
func IPAdressByRequestCF(ctx *fasthttp.RequestCtx) (ip string) {
	if ip = string(ctx.Request.Header.Peek("Cf-Connecting-Ip")); len(ip) > 1 {
		realIP := net.ParseIP(ip)
		if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
			return IPAdressByRequest(ctx)
		}
		return
	}
	return IPAdressByRequest(ctx)
}

func getIPAdress(ips string) (res string) {
	for _, ip := range strings.Split(ips, ",") {
		// header can contain spaces too, strip those out.
		ip = strings.TrimSpace(ip)
		realIP := net.ParseIP(ip)

		if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
			// bad address, go to next
			continue
		}

		res = ip
		break
	}
	return
}
