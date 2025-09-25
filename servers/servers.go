package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	Ports []int
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
	port := s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}

func RunServers(amount int) {
	// Serverlist object
	var myServerList ServerList
	myServerList.Populate(amount)

	// Waitgroup
	var wg sync.WaitGroup
	wg.Add(amount)
	defer wg.Wait()

	for x := 0; x < amount; x++ {
		go makeServers(&myServerList, &wg)
	}
}

func makeServers(sl *ServerList, wg *sync.WaitGroup) {
	// Router
	r := http.NewServeMux()
	defer wg.Done()
	// Port
	port := sl.Pop()

	server := http.Server{
		Addr:    fmt.Sprintf(":808%d", port),
		Handler: r,
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server running on port: %d", port)
	})

	r.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("400 - Server Shut Down!"))
		server.Shutdown(context.Background())
	})

	server.ListenAndServe()
}
