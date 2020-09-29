package main

import (
    "fmt"
	"net"
	"sort"
	"flag"
)

func main() {
	host, max_ports := getArgs()
	// gets the host and ports to scan

	results := make(chan int)
	ports := make(chan int, 100)
	// scan 100 ports at a time

	var open_ports []int
	// list to store all the open ports

	for i := 0; i < cap(ports); i++ {
		go scanner(ports, results, host)
		// concurrently scanning ports
		// scannign ports before channel reaches 100
	}

	go func(){
		for j := 1; j < max_ports; j++ {
			ports <- j
			// continually adding ports to be used
     	}
	}()

	for k := 0; k  < max_ports; k++ {
		port := <- results
		// port is sent from results channel

		if port != 0 {
			open_ports = append(open_ports, port)
		}
	} 

	close(ports)
	close(results)
	// closes the channels

	sort.Ints(open_ports)
	// sorts into asc order

	for _, port := range open_ports {
		fmt.Printf("port: %d OPEN \n", port)
	}
}

func getArgs() (string, int) {
	
	h := flag.String("host", "127.0.0.1", "the network whose ports are to be scanned")
	mp := flag.Int("max_ports", 1024, "starting from 1 ports to be scanned")
	flag.Parse()
	// gets flags from cmd

	host := *h
	max_ports := *mp

	return host, max_ports
}

func scanner(ports chan int, results chan int, host string) {
	for port := range ports {
		address := fmt.Sprintf("%s:%d", host, port)
		// formats address to be used by net.Dial

		connection, err := net.Dial("tcp", address)
		// scans the port 

		if err != nil {
			results <- 0
			// sends negative result if it does not work
			continue
		}
		connection.Close()
		// shuts down once scanned
		results <- port 
	}
}