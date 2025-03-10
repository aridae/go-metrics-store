package subnet

import (
	"fmt"
	"net"
	"net/http"

	"github.com/aridae/go-metrics-store/pkg/slice"
	"github.com/aridae/go-metrics-store/pkg/subnet"
)

const (
	realIPHeader            = "X-Real-IP"
	subnetNotTrustedMessage = "The IP address in the request does not belong to a trusted subnet"
)

func ValidateTrustedSubnetMiddleware(whitelist ...*net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, err := subnet.ResolveRequestIP(r, realIPHeaderHostSelector)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			isTrusted := slice.ContainsByFunc(whitelist, func(trustedNet *net.IPNet) bool {
				return trustedNet.Contains(ip)
			})
			if !isTrusted {
				http.Error(w, subnetNotTrustedMessage, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func realIPHeaderHostSelector(req *http.Request) (string, error) {
	if host := req.Header.Get(realIPHeader); host != "" {
		return host, nil
	}

	return "", fmt.Errorf("%s header is not set, but expected to be", realIPHeader)
}
