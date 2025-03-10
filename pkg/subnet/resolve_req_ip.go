package subnet

import (
	"net"
	"net/http"

	"github.com/aridae/go-metrics-store/pkg/logger"
	"go.uber.org/multierr"
)

type HostSelectorFunc func(req *http.Request) (string, error)

func ResolveRequestIP(req *http.Request, hostSelectors ...HostSelectorFunc) (net.IP, error) {
	var ip net.IP
	var multierror error

	for _, hostSelector := range hostSelectors {
		host, err := hostSelector(req)
		if err != nil {
			logger.Errorf("[subnet.ResolveRequestIP] host selector error: %v", err)

			multierror = multierr.Append(multierror, err)
			continue
		}

		ip = net.ParseIP(host)
		if ip == nil {
			logger.Errorf("[subnet.ResolveRequestIP] ParseIP unexpectedly returned nil IP for <host:%s>", host)
			continue
		}

		break
	}

	if ip != nil {
		return ip, nil
	}

	return nil, multierror
}
