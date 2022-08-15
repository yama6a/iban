package http

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockController struct {
	t                *testing.T
	SetupRoutesFunc  func()
	setupRoutesCalls int
}

func (ctrl *mockController) SetupRoutes() {
	if ctrl.SetupRoutesFunc == nil {
		ctrl.t.Fatalf("mockController.SetupRoutesFunc: method is nil but Controller.SetupRoutes was just called")
	}
	ctrl.SetupRoutesFunc()
	ctrl.setupRoutesCalls++
}

func (ctrl *mockController) GetSetupRoutesCalls() int {
	return ctrl.setupRoutesCalls
}

func TestHttpServer_setupRoutes(t *testing.T) {
	tests := []struct {
		name          string
		controllers   []Controller
		expectedCalls int
	}{
		{
			name:          "no controllers succeeds",
			controllers:   []Controller{},
			expectedCalls: 0,
		},
		{
			name: "adding three controllers results in three calls",
			controllers: []Controller{
				&mockController{
					t:               t,
					SetupRoutesFunc: func() {},
				},
				&mockController{
					t:               t,
					SetupRoutesFunc: func() {},
				},
				&mockController{
					t:               t,
					SetupRoutesFunc: func() {},
				},
			},
			expectedCalls: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt // shadow tt for parallel execution
			t.Parallel()

			// there should be no calls before we set up the routes
			calls := 0
			for _, controller := range tt.controllers {
				calls += controller.(*mockController).GetSetupRoutesCalls()
			}
			require.Equal(t, 0, calls)

			s := &Server{logger: zap.NewNop(), controllers: tt.controllers}
			s.setupRoutes()

			// There should be tt.expectedCalls calls after we set up the routes
			// and every controller should have been called exactly once
			calls = 0
			for _, controller := range tt.controllers {
				ctrl := controller.(*mockController)
				calls += ctrl.GetSetupRoutesCalls()
				require.Equal(t, 1, ctrl.GetSetupRoutesCalls())
			}
			require.Equal(t, tt.expectedCalls, calls)
		})
	}
}
