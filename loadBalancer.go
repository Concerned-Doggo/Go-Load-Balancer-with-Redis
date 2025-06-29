package main

import (
	"net/http"
	"sync"
)

type LoadBalancer struct {
	Port          string
	RoundRobinCnt int
	Servers       []Server
	mu            sync.Mutex
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		Port:          port,
		RoundRobinCnt: 0,
		Servers:       servers,
	}
}

func (lb *LoadBalancer) GetNextAvailableServer() Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.Servers) == 0 {
		return nil
	}

	server := lb.Servers[lb.RoundRobinCnt%len(lb.Servers)]
	lb.RoundRobinCnt++

	return server
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.GetNextAvailableServer()
	if targetServer == nil {
		http.Error(rw, "No Server Available", http.StatusServiceUnavailable)
		return
	}
	targetServer.Serve(rw, r)
}

func (lb *LoadBalancer) ChartServerProxy(rw http.ResponseWriter, r *http.Request) {

	targetServer := lb.GetNextAvailableServer()
	if targetServer == nil {
		http.Error(rw, "No Server Available", http.StatusServiceUnavailable)
		return
	}
	targetServer.ServeChart(rw, r)
}


