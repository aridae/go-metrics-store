package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"sync"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	port := 8082

	srv := NewServer(port, nil)
	assert.NotNil(t, srv, "Server should not be nil")
	assert.Equal(t, port, srv.port, "Port should match")
	assert.IsType(t, &grpc.Server{}, srv.server, "server should be of type *grpc.Server")
}

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fakeServer := &fakeGrpcServer{}
	port := 8082
	srv := &Server{
		server: fakeServer,
		port:   port,
	}

	err := srv.Run(ctx)
	assert.NoError(t, err, "Running server should succeed")

	cancel()
	time.Sleep(time.Second) // to avoid race

	assert.True(t, fakeServer.WasStartCalled(), "Serve method should have been called")
	assert.True(t, fakeServer.WasStopCalled(), "Stop method should have been called after context cancellation")
}

func TestServer_Run_PortAlreadyInUse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fakeServer := &fakeGrpcServer{}
	port := 7012
	srv := &Server{
		server: fakeServer,
		port:   port,
	}
	err := srv.Run(ctx)
	assert.NoError(t, err, "Running server should succeed")

	// Call one more time over the same port; gotta get error - port already in use
	err = srv.Run(ctx)
	assert.Error(t, err, "Running server on busy port should return error")
}

// Fake implementation of grpc.Server interface
type fakeGrpcServer struct {
	startCalled bool
	stopCalled  bool
	mux         sync.RWMutex
}

func (f *fakeGrpcServer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	// No-op for testing purposes
}

func (f *fakeGrpcServer) Serve(_ net.Listener) error {
	f.mux.Lock()
	defer f.mux.Unlock()

	f.startCalled = true
	return nil
}

func (f *fakeGrpcServer) Stop() {
	f.mux.Lock()
	defer f.mux.Unlock()

	f.stopCalled = true
}

func (f *fakeGrpcServer) WasStartCalled() bool {
	f.mux.RLock()
	defer f.mux.RUnlock()

	return f.startCalled
}

func (f *fakeGrpcServer) WasStopCalled() bool {
	f.mux.RLock()
	defer f.mux.RUnlock()

	return f.stopCalled
}
