package web

import (
	"net"
	"testing"
)

func TestIPNet(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("127.0.0.1/32")
	if err != nil {
		t.Fatal(err)
	}
	ok := contains("127.0.0.1", ipNet)
	t.Log(ok)
}
