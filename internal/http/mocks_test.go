package http

import "testing"

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
