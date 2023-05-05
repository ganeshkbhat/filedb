package netserve

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"

	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	// "time"

	"github.com/soheilhy/cmux"
	// grpc "google.golang.org/grpc"
)

// Summary:
// The serverHTTPHandler struct represents an implementation of the http.Handler interface.
// Args:
// N/A
// Returns:
// N/A
type serverHTTPHandler struct{}

// Summary:
// `serverRPCRcvr` is a struct type in Go language used by the `net/rpc` package to receive and handle RPC messages sent to a server.
// Args:
// This struct type does not accept any arguments.
// Returns:
// None. This is a struct type and does not have a return statement.
type serverRPCRcvr struct{}

// Summary:
// This function handles incoming HTTP requests and returns a sample HTTP response.
// Args:
// - w: An object of the type http.ResponseWriter used to write HTTP response.
// - r: An object of the type http.Request used to read HTTP request.
// Returns:
// This function does not return any values, but it writes "example http response" to the HTTP response object.
func (h *serverHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example http response")
}


// Summary:
// The serveHTTP function creates an HTTP server and sets up a handler for it to use, then serves incoming requests on the provided listener.
// Args:
// - l: A net.Listener representing the listener to serve requests on.
// Returns:
// None (void function).
func serveHTTP(l net.Listener) {
	s := &http.Server{
		Handler: &serverHTTPHandler{},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

// Summary:
// This function serves HTTPS over a specified listener using SSL/TLS certificates loaded from disk.
// Args:
// - l (net.Listener): The listener to serve HTTPS on.
// - listenport (bool): A boolean value indicating whether to create a new listener on the default HTTPS port (443).
// Returns:
// No return statement, as the function only performs actions.
func serveHTTPS(l net.Listener, listenport bool) {
	// Load certificates.
	certificate, err := tls.LoadX509KeyPair("./certs/ssl.cert", "./certs/ssl.key")
	if err != nil {
		log.Panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
	}

	// Create TLS listener.
	tlsl := tls.NewListener(l, config)

	if listenport {
		go func() {
			// // Create custom listener
			ln, err := tls.Listen("tcp", ":443", config)
			if err != nil {
				log.Fatal(err)
			}
			// Start server with https/ssl enabled on http://localhost:443
			log.Fatal(ln)
		}()
	}

	// Serve HTTP over TLS.
	serveHTTP(tlsl)
}

// This is an example for serving HTTP and HTTPS on the same port.

// Summary:
// Netserve is a function that creates a TCP listener using the specified protocol and IP address and serves HTTP and HTTPS protocols using cmux package.
// It also listens for interrupt or termination signal to trigger graceful shutdown.
// Args:
// - netprotocol (string): network protocol to be used for creating TCP listener.
// - netipport (string): IP address and port number of the listener.
// Returns:
// void
func Netserve(netprotocol string, netipport string) {
	// Create the TCP listener.
	l, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Panic(err)
	}

	// Create a mux.
	m := cmux.New(l)

	// We first match on HTTP 1.1 methods.
	httpl := m.Match(cmux.HTTP1Fast())

	// If not matched, we assume that its TLS.
	//
	// Note that you can take this listener, do TLS handshake and
	// create another mux to multiplex the connections over TLS.
	tlsl := m.Match(cmux.Any())

	go serveHTTP(httpl)
	go serveHTTPS(tlsl, false)

	// http2 := tcpm.Match(cmux.HTTP2())

	// Listen for the process signal to trigger grceful shutdown When an interrupt or termination signal is sent
	go func() {
		// sig := make(chan os.Signal, 1)
		// signal.Notify(sig, syscall.SIGHUP)
		// for range sig {
		// 	upg.Upgrade()
		// }
		c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
		signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

		_ = <-c // This blocks the main thread until an interrupt is received
		fmt.Println("Gracefully shutting down.")
		_ = l.Close()
		fmt.Println("Running cleanup tasks.")

		// Your cleanup tasks go here
		// db.Close()
		// redisConn.Close()
		fmt.Println("Fiber was successful shutdown.")
	}()

	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}
