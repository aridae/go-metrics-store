package subnet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// ExampleResolveRequestIP демонстрирует использование метода ResolveRequestIP.
func ExampleResolveRequestIP() {
	// Определение селекторов для извлечения IP адреса.
	getRemoteAddr := func(req *http.Request) (string, error) {
		return req.RemoteAddr, nil
	}

	getForwardedFor := func(req *http.Request) (string, error) {
		return req.Header.Get("X-Forwarded-For"), nil
	}

	// Создание HTTP запроса.
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "192.168.0.1")

	// Попытка извлечь IP адрес из запроса.
	ip, err := ResolveRequestIP(req, getRemoteAddr, getForwardedFor)
	if err != nil {
		fmt.Println("Failed to resolve client IP:", err)
	} else {
		fmt.Println("Client IP:", ip.String())
	}

	// Output:
	// Client IP: 192.168.0.1
}
