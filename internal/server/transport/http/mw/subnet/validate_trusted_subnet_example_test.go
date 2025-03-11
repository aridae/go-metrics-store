package subnet

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
)

// ExampleValidateTrustedSubnetMiddleware демонстрирует использование middleware.
func ExampleValidateTrustedSubnetMiddleware() {
	// Определяем доверенные подсети.
	_, trustSubnet1, _ := net.ParseCIDR("192.168.0.0/24")
	_, trustSubnet2, _ := net.ParseCIDR("10.0.0.0/16")

	// Создаем новый middleware.
	middleware := ValidateTrustedSubnetMiddleware(trustSubnet1, trustSubnet2)

	// Оборачиваем тестовый обработчик middleware.
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	}))

	// Создаем HTTP-запросы с разными IP-адресами.
	reqWithTrustedIP := httptest.NewRequest(http.MethodGet, "/", nil)
	reqWithUntrustedIP := httptest.NewRequest(http.MethodGet, "/", nil)

	// Устанавливаем заголовок X-Real-IP с доверенным IP-адресом.
	reqWithTrustedIP.Header.Set(realIPHeader, "192.168.0.1")

	// Устанавливаем заголовок X-Real-IP с недоверенным IP-адресом.
	reqWithUntrustedIP.Header.Set(realIPHeader, "172.20.30.40")

	// Выполняем запросы через middleware.
	recorderTrusted := httptest.NewRecorder()
	recorderUntrusted := httptest.NewRecorder()

	handler.ServeHTTP(recorderTrusted, reqWithTrustedIP)
	handler.ServeHTTP(recorderUntrusted, reqWithUntrustedIP)

	// Проверка результата
	respBody, _ := io.ReadAll(recorderTrusted.Body)
	fmt.Println("Ответ от обработчика запроса из доверенной подсети:", string(respBody))

	respBody, _ = io.ReadAll(recorderUntrusted.Body)
	fmt.Println("Ответ от обработчика запроса не из доверенной подсети:", string(respBody))

	// Output:
	// Ответ от обработчика запроса из доверенной подсети: Hello, World!
	// Ответ от обработчика запроса не из доверенной подсети: The IP address in the request does not belong to a trusted subnet
}
