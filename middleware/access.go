package middleware

import (
	"net"
	"net/http"

	"github.com/corex-io/web"
)

// AccessIP access
func AccessIP(cidrs ...string) func(*web.Context) {
	var ipnet []*net.IPNet
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		ipnet = append(ipnet, ipNet)
	}
	return func(ctx *web.Context) {
		addr := net.ParseIP(ctx.Remote())
		for _, cidr := range ipnet {
			if cidr.Contains(addr) {
				return
			}
		}
		ctx.Error(http.StatusForbidden)
	}
}
