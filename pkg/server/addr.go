package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func buildAddr(ip net.IP, port int) string {
	if strings.Contains(ip.String(), ".") {
		return fmt.Sprintf("%s:%d", ip, port)
	}
	return fmt.Sprintf("[%s]:%d", ip, port)
}

func remoteAddr(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if strings.Contains(IPAddress, " ") {
		IPAddress = strings.Split(IPAddress, " ")[0]
	}
	if IPAddress == "" || net.ParseIP(IPAddress) == nil {
		IPAddress = r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
	}
	return IPAddress
}
