package main

import (
	"fmt"
	"net/http"
	"time"

	loadbalancer "github.com/thealexcons/ServerLoadBalancer"
)

func main() {
	// Create the server group and add nodes to it
	sg := &loadbalancer.ServerGroup{}
	sg.AddNode("http://localhost:1232", 1)
	// this node has weight 2 (will receive twice the no. of requests compared to other nodes)
	sg.AddNode("http://localhost:4192", 2)
	sg.AddNode("http://localhost:4311", 1)

	sg.SetRetries(3)

	// Spin up the example nodes above
	runExampleNodeServers()

	// Create the server at port 8080 and handle requests using the
	// load balanacer
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(sg.LoadBalancer),
	}

	// Start health checking routine every 30 seconds, with a timeout of 3 seconds
	sg.StartHealthChecker(time.Second*30, time.Second*3)

	// Start the server
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Start some example nodes to simulate the load balancer
// Obviously, in practice, these nodes would be running separately.
func runExampleNodeServers() {
	server1 := &http.Server{
		Addr:    ":4311",
		Handler: http.HandlerFunc(exampleHandler),
	}
	go server1.ListenAndServe()

	server2 := &http.Server{
		Addr:    ":1232",
		Handler: http.HandlerFunc(exampleHandler),
	}
	go server2.ListenAndServe()

	server3 := &http.Server{
		Addr:    ":4192",
		Handler: http.HandlerFunc(exampleHandler),
	}
	go server3.ListenAndServe()
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at " + r.URL.Host)
	fmt.Fprintf(w, "Serve content at "+r.URL.Path)
}
