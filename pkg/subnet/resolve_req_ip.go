package subnet

import (
	"fmt"
	"net"
	"net/http"

	"github.com/aridae/go-metrics-store/pkg/logger"
	"go.uber.org/multierr"
)

type HostSelectorFunc func(req *http.Request) (string, error)

// ResolveRequestIP пытается извлечь IP адрес клиента из HTTP запроса, используя предоставленные функции-селекторы.
//
// Селектор - это функция, которая получает HTTP запрос и возвращает строку, содержащую представление IP адреса.
// Этот метод пробует каждый селектор до тех пор, пока не найдет валидный IP адрес.
// Если ни один из селекторов не вернул валидный IP адрес, возвращается ошибка.
//
// Параметры:
//
//	req - HTTP запрос, из которого нужно извлечь IP адрес.
//	hostSelectors - Набор функций-селекторов, которые пытаются извлечь IP адрес из запроса.
//
// Возвращаемые значения:
//
//	net.IP - Извлечённый IP адрес, если один из селекторов вернул валидный IP адрес.
//	error - Ошибка, если ни один из селекторов не смог вернуть валидный IP адрес.
//
// Примеры:
//
//	func getRemoteAddr(req *http.Request) (string, error) {
//	    return req.RemoteAddr, nil
//	}
//
//	func getForwardedFor(req *http.Request) (string, error) {
//	    return req.Header.Get("X-Forwarded-For"), nil
//	}
//
//	ip, err := ResolveRequestIP(req, getRemoteAddr, getForwardedFor)
//	if err != nil {
//	    log.Printf("Failed to resolve client IP: %v", err)
//	} else {
//	    log.Printf("Client IP: %s", ip.String())
//	}
//
// Ошибки:
//
//   - Если ни один из селекторов не вернул валидный IP адрес, возвращается ошибка с сообщением о неудаче.
//   - Если селектор возвращает пустую строку или некорректный IP адрес, также возвращается ошибка.
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
			err = fmt.Errorf("ParseIP unexpectedly returned nil IP for <host:%s>", host)

			logger.Errorf("[subnet.ResolveRequestIP] failed to parse IP: %v", err)

			multierror = multierr.Append(multierror, err)
			continue
		}

		break
	}

	if ip != nil {
		return ip, nil
	}

	return nil, multierror
}
