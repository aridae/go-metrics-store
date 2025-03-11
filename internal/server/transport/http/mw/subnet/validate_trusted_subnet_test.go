package subnet

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestValidateTrustedSubnetMiddlewareSuccess проверяет успешную валидацию.
func TestValidateTrustedSubnetMiddleware_SuccessCase(t *testing.T) {
	// Создание доверенной подсети.
	_, trustSubnet, _ := net.ParseCIDR("192.168.0.0/24")

	// Создание middleware.
	middleware := ValidateTrustedSubnetMiddleware(trustSubnet)

	// Создание тестового обработчика.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Test passed.")) // Сообщение успешного прохождения теста.
	})

	// Создание HTTP-запроса с доверенным IP-адресом.
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(realIPHeader, "192.168.0.123") // Доверенный IP адрес.

	// Создание записывающего ответа.
	recorder := httptest.NewRecorder()

	// Вызов middleware с тестовым обработчиком.
	handler := middleware(testHandler)
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа.
	if recorder.Code != http.StatusOK {
		t.Fatalf("Unexpected response status: %d, expected: %d", recorder.Code, http.StatusOK)
	}

	// Проверка тела ответа.
	expectedBody := "Test passed."
	if body := strings.TrimSpace(recorder.Body.String()); body != expectedBody {
		t.Fatalf("Unexpected response body: %q, expected: %q", body, expectedBody)
	}
}

// TestValidateTrustedSubnetMiddlewareError проверяет поведение middleware при возникновении ошибки парсинга IP клиента из хедера.
func TestValidateTrustedSubnetMiddlewareError_InvalidClientRealIP(t *testing.T) {
	// Создание доверенной подсети.
	_, trustSubnet, _ := net.ParseCIDR("192.168.0.0/24")

	// Создание middleware.
	middleware := ValidateTrustedSubnetMiddleware(trustSubnet)

	// Создание тестового обработчика.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Test passed."))
	})

	// Создание HTTP-запроса с неверным IP адресом.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(realIPHeader, "some-really-invalid-ip")

	recorder := httptest.NewRecorder()
	handler := middleware(testHandler)
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа.
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Unexpected response status: %d, expected: %d", recorder.Code, http.StatusInternalServerError)
	}
}

// TestValidateTrustedSubnetMiddlewareError проверяет негативный кейс валидации.
func TestValidateTrustedSubnetMiddlewareError_UntrustedClientRealIP(t *testing.T) {
	// Создание доверенной подсети.
	_, trustSubnet, _ := net.ParseCIDR("192.168.0.0/24")

	// Создание middleware.
	middleware := ValidateTrustedSubnetMiddleware(trustSubnet)

	// Создание тестового обработчика.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Test passed."))
	})

	// Создание HTTP-запроса с неверным IP адресом.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(realIPHeader, "172.18.19.20") // Недоверенный IP адрес.

	recorder := httptest.NewRecorder()
	handler := middleware(testHandler)
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа.
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("Unexpected response status: %d, expected: %d", recorder.Code, http.StatusForbidden)
	}

	// Проверка тела ответа.
	expectedBody := subnetNotTrustedMessage
	if body := strings.TrimSpace(recorder.Body.String()); body != expectedBody {
		t.Fatalf("Unexpected response body: %q, expected: %q", body, expectedBody)
	}
}
