package extensions

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"

	"Aether/pkg/fs"
)

// Server handles serving extension UI files locally for iframes
type Server struct {
	port int
}

// NewServer creates a new local extension server
func NewServer() *Server {
	return &Server{}
}

// Start launches the HTTP server on a random available port and returns the base URL
func (s *Server) Start() (string, error) {
	extDir := filepath.Join(fs.GetDataDir(), "extensions")
	fsHandler := http.FileServer(http.Dir(extDir))

	mux := http.NewServeMux()
	mux.Handle("/", enableCORS(fsHandler))

	// Bind to port 0 to let the OS assign a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("failed to start local server: %w", err)
	}

	s.port = listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d", s.port)

	go func() {
		fmt.Printf("[Extensions] UI Server listening at %s\n", url)
		if err := http.Serve(listener, mux); err != nil {
			fmt.Printf("[Extensions] Server error: %v\n", err)
		}
	}()

	return url, nil
}

// GetPort returns the currently bound port
func (s *Server) GetPort() int {
	return s.port
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow our Wails frontend (usually wails:// or http://localhost:...) to load these assets
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
