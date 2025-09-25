package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

var (
	baseURL = "http://localhost:808"
)

type LoadBalancer struct {
	ReverseProxy httputil.ReverseProxy
}
type Endpoints struct {
	List []*url.URL
}

func (e *Endpoints) Shuffle() {
	temp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, temp)
}

func MakeLoadBalancer(amount int) {
	// Instantiate Objects
	var lb LoadBalancer
	var ep Endpoints

	// Server + Router
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	for i := 0; i < amount; i++ {
		ep.List = append(ep.List, createEndpoint(baseURL, i))
	}

	// HandlerFunctions
	router.HandleFunc("/loadBalancer", makeRequest(&lb, &ep))

	log.Fatal(server.ListenAndServe())
}

func makeRequest(lb *LoadBalancer, ep *Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for !testServer(ep.List[0].String()) {
			ep.Shuffle()
		}
		lb.ReverseProxy = *httputil.NewSingleHostReverseProxy(ep.List[0])
		ep.Shuffle()
		lb.ReverseProxy.ServeHTTP(w, r)
	}
}

func createEndpoint(endpoint string, idx int) *url.URL {
	link := endpoint + strconv.Itoa(idx)
	url, _ := url.Parse(link)
	return url
}

func testServer(endpoint string) bool {
	resp, err := http.Get(endpoint)
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}
