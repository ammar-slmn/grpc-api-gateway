package loadbalancer_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"grpc-api-gateway/pkg/loadbalancer"
)

func TestLoadBalancer(t *testing.T) {
	t.Run("Basic Handler", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		lb, ep := loadbalancer.NewLoadBalancer(1, ts.URL)
		req := httptest.NewRequest("GET", "/loadBalancer", nil)
		w := httptest.NewRecorder()

		handler := loadbalancer.MakeHandler(lb, ep)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestEndpointsShuffle(t *testing.T) {
	testURL := "http://test"
	_, ep := loadbalancer.NewLoadBalancer(3, testURL)
	firstURL := ep.List[0].String()
	ep.Shuffle()

	if ep.List[len(ep.List)-1].String() != firstURL {
		t.Error("Shuffle did not move first element to end of list")
	}
}

func TestHealthCheck(t *testing.T) {
	t.Run("Healthy Server", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		lb, ep := loadbalancer.NewLoadBalancer(1, ts.URL)
		handler := loadbalancer.MakeHandler(lb, ep)

		req := httptest.NewRequest("GET", "/loadBalancer", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Error("Health check failed for healthy server")
		}
	})

	t.Run("Unhealthy Server", func(t *testing.T) {
		lb, ep := loadbalancer.NewLoadBalancer(1, "http://localhost:1234")
		handler := loadbalancer.MakeHandler(lb, ep)

		req := httptest.NewRequest("GET", "/loadBalancer", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Error("Health check passed for unhealthy server")
		}
	})
}
