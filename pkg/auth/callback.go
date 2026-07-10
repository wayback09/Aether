package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// StartCallbackServer starts a temporary HTTP server on a free port,
// calls onStart with the chosen port (so the browser can be opened),
// waits for the Microsoft OAuth2 redirect, and returns the authorization code.
// The server shuts itself down immediately after receiving the code.
func StartCallbackServer(onStart func(port int)) (code string, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("could not bind callback server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	// Notify caller that server is listening and they can open the browser
	if onStart != nil {
		onStart(port)
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	srv := &http.Server{Handler: mux}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error_description")
			if errMsg == "" {
				errMsg = r.URL.Query().Get("error")
			}
			http.Error(w, "Authorization failed: "+errMsg, http.StatusBadRequest)
			errCh <- fmt.Errorf("Microsoft auth denied: %s", errMsg)
			return
		}

		// Serve a friendly success page before the browser is closed
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Aether — Login Successful</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; background: #0e0e10; color: #e5e5e5;
         display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
  .card { text-align: center; padding: 48px; }
  h1 { font-size: 24px; margin-bottom: 8px; color: #fff; }
  p { color: #888; font-size: 14px; }
</style>
</head>
<body>
  <div class="card">
    <h1>✓ Signed in successfully</h1>
    <p>You can close this tab and return to Aether.</p>
  </div>
</body>
</html>`)

		codeCh <- code

		// Shut down the server gracefully in a goroutine
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}()
	})

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for either a code or an error, with a 5-minute timeout
	select {
	case c := <-codeCh:
		return c, nil
	case e := <-errCh:
		return "", e
	case <-time.After(5 * time.Minute):
		srv.Close()
		return "", fmt.Errorf("login timed out — the browser was not completed within 5 minutes")
	}
}
