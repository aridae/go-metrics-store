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

// ValidateTrustedSubnetMiddleware создаёт HTTP middleware, который проверяет, принадлежит ли IP адрес клиента
// к одной из доверенных подсетей, указанных в белом списке.
//
// Мiddleware извлекает IP адрес клиента из заголовка "X-Real-IP". Если адрес отсутствует или не принадлежит ни одной
// из доверенных подсетей, возвращается ошибка с кодом статуса HTTP 403 Forbidden.
//
// Параметры:
//
//	whitelist - Список указателей на структуры net.IPNet, определяющие доверенные подсети.
//
// Возвращаемые значения:
//
//	Функция возвращает обработчик http.Handler, который оборачивает указанный next-обработчик.
//
// Примеры:
//
//	mw := ValidateTrustedSubnetMiddleware(
//		net.ParseCIDR("192.168.0.0/24"),
//		net.ParseCIDR("10.0.0.0/16"),
//	)
//	handler := mw(http.DefaultServeMux)
//
// Ошибки:
//
//   - В случае ошибки при разборе заголовка "X-Real-IP" будет возвращена ошибка с кодом статуса HTTP 500 Internal Server Error.
//   - Если IP адрес клиента не принадлежит ни одной из доверенных подсетей, будет возвращена ошибка с кодом статуса HTTP 403 Forbidden.
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
