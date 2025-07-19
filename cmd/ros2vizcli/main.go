package main

import (
	"flag"
	"log"
	"time"

	"ros2viz/src/ros2viz" /
)

func main() {
	// --- Command-Line Flag Parsing ---
	addr := flag.String("addr", ":8080", "HTTP network address")
	scriptPath := flag.String("script", "scripts/inspect_ros_graph.py", "Path to the Python introspection script")
	pollInterval := flag.Duration("poll", 2*time.Second, "How often to poll the ROS graph")
	flag.Parse()

	log.Println("--- ROS 2 Visualizer Backend ---")

	// --- Dependency Injection ---
	rosInspector := ros2viz.NewROSInspector(*scriptPath)
	wsHub := ros2viz.NewHub()
	appServer := ros2viz.NewServer(*addr, wsHub, rosInspector)

	// --- Start Application Components ---
	log.Println("Starting WebSocket hub...")
	go wsHub.Run()

	log.Println("Starting ROS graph polling...")
	go appServer.PollROSGraph(*pollInterval)

	log.Printf("Starting server on %s", *addr)
	if err := appServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
