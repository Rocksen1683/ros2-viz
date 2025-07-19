package ros2viz 

import (
	"log"
	"net/http"
	"time"
)

// Server holds the dependencies for the web server
type Server struct {
	addr      string
	hub       *Hub             
	inspector DataProvider     
}

// NewServer creates a new configured Server
func NewServer(addr string, h *Hub, i DataProvider) *Server {
	return &Server{
		addr:      addr,
		hub:       h,
		inspector: i,
	}
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(s.hub, w, r) 
	})
	return http.ListenAndServe(s.addr, mux)
}

// PollROSGraph periodically fetches ROS data and broadcasts it
func (s *Server) PollROSGraph(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		if s.hub.HasClients() {
			data, err := s.inspector.GetROSGraphData()
			if err != nil {
				log.Println("Could not get ROS graph data, skipping broadcast.")
				continue
			}
			s.hub.Broadcast <- data
		}
	}
}