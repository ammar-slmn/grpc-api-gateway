package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

type LoadBalancer struct {
	ReverseProxy httputil.ReverseProxy
	mu           sync.Mutex
	baseURL      string
}

type Endpoints struct {
	List []*url.URL
}

func (e *Endpoints) Shuffle() {
	temp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, temp)
}

func NewLoadBalancer(amount int, baseURL string) (*LoadBalancer, *Endpoints) {
	if baseURL == "" {
		baseURL = "http://localhost:808"
	}

	lb := &LoadBalancer{
		baseURL: baseURL,
	}
	ep := &Endpoints{}

	for i := 0; i < amount; i++ {
		ep.List = append(ep.List, createEndpoint(baseURL, i))
	}

	return lb, ep
}

func MakeHandler(lb *LoadBalancer, ep *Endpoints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lb.mu.Lock()
		defer lb.mu.Unlock()

		for !checkServerHealth(ep.List[0].String()) {
			ep.Shuffle()
		}
		lb.ReverseProxy = *httputil.NewSingleHostReverseProxy(ep.List[0])
		ep.Shuffle()
		lb.ReverseProxy.ServeHTTP(w, r)
	}
}

func createEndpoint(endpoint string, idx int) *url.URL {
	link := endpoint + strconv.Itoa(idx)
	url, err := url.Parse(link)
	if err != nil {
		return nil
	}
	return url
}

func checkServerHealth(endpoint string) bool {
	resp, err := http.Get(endpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
