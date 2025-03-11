package subnet

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestResolveRequestIP проверяет работу метода ResolveRequestIP.
func TestResolveRequestIP(t *testing.T) {
	type args struct {
		req           *http.Request
		hostSelectors []HostSelectorFunc
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		wantErr bool
	}{
		{
			name: "Valid IP from first selector",
			args: args{
				req: httptest.NewRequest("GET", "/path", nil),
				hostSelectors: []HostSelectorFunc{
					func(req *http.Request) (string, error) {
						return "127.0.0.1", nil
					},
					func(req *http.Request) (string, error) {
						return "192.168.0.1", nil
					},
				},
			},
			want:    net.ParseIP("127.0.0.1"),
			wantErr: false,
		},
		{
			name: "Valid IP from second selector",
			args: args{
				req: httptest.NewRequest("GET", "/path", nil),
				hostSelectors: []HostSelectorFunc{
					func(req *http.Request) (string, error) {
						return "", fmt.Errorf("empty string")
					},
					func(req *http.Request) (string, error) {
						return "192.168.0.1", nil
					},
				},
			},
			want:    net.ParseIP("192.168.0.1"),
			wantErr: false,
		},
		{
			name: "All selectors fail",
			args: args{
				req: httptest.NewRequest("GET", "/path", nil),
				hostSelectors: []HostSelectorFunc{
					func(req *http.Request) (string, error) {
						return "", fmt.Errorf("first selector fails")
					},
					func(req *http.Request) (string, error) {
						return "", fmt.Errorf("second selector fails")
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveRequestIP(tt.args.req, tt.args.hostSelectors...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveRequestIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveRequestIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}
