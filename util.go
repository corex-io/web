package web

import (
	"net"
	"net/http"
	"os"
	"strings"
)

func contains(ip string, cidrs ...*net.IPNet) bool {
	addr := net.ParseIP(ip)
	for _, cidr := range cidrs {
		if cidr.Contains(addr) {
			return true
		}
	}
	return false
}

// If s starts with one of suffixs; return ture
func hasSuffixs(s string, suffixs ...string) bool {
	for _, suffix := range suffixs {
		if ok := strings.HasSuffix(s, suffix); ok {
			return true
		}
	}
	return false
}

func hasPrefixs(s string, prefixs ...string) bool {
	for _, prefix := range prefixs {
		if ok := strings.HasPrefix(s, prefix); ok {
			return true
		}
	}
	return false
}

func toHTTPError(err error) int {
	if os.IsNotExist(err) {
		return http.StatusNotFound
	}
	if os.IsPermission(err) {
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}
