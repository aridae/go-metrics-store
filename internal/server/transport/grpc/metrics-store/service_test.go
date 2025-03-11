package metricsstore

import (
	"github.com/aridae/go-metrics-store/internal/server/transport/grpc/metrics-store/_mock"
	"go.uber.org/mock/gomock"
	"testing"
)

type testKit struct {
	useCasesController *_mock.MockuseCasesController
	service            *Implementation
}

func setUpTestKit(t *testing.T) *testKit {
	ctrl := gomock.NewController(t)

	useCasesControllerMock := _mock.NewMockuseCasesController(ctrl)

	service := &Implementation{useCasesController: useCasesControllerMock}

	return &testKit{
		useCasesController: useCasesControllerMock,
		service:            service,
	}
}
