package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	Ports []int
	mu    sync.Mutex
}

func (s *ServerList) Populate(amount int) {
	if amount >= 10 {
		log.Fatal("Amount of Ports cant exceed 10")
	}
	for x := 0; x < amount; x++ {
		s.Ports = append(s.Ports, x)
	}
}

func (s *ServerList) Pop() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Ports) == 0 {
		log.Fatal("No more ports available")
	}
	port := s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}

func RunServers(amount int) {
	var myServerList ServerList
	myServerList.Populate(amount)

	var wg sync.WaitGroup
	wg.Add(amount)
	defer wg.Wait()

	for x := 0; x < amount; x++ {
		go makeServer(&myServerList, &wg)
	}
}

func makeServer(sl *ServerList, wg *sync.WaitGroup) {
	defer wg.Done()

	port := sl.Pop()
	addr := fmt.Sprintf(":808%d", port)

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server running on port: %d", port)
	})

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Server on port %s failed: %v", addr, err)
	}
}
