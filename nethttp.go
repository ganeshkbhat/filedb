package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	// "net/rpc"

	"os"
	"os/signal"
	"strings"
	"syscall"

	// "google.golang.org/grpc"
	// grpc "google.golang.org/grpc"
	// codes "google.golang.org/grpc/codes"
	// status "google.golang.org/grpc/status"

	// "golang.org/x/net/context"
	"golang.org/x/net/websocket"

	"github.com/soheilhy/cmux"
)

// Summary:
// This is a struct declaration that defines a type named `serverHTTPHandler`.
// Args:
// There are no arguments for this struct declaration.
// Returns:
// There is no return statement for this struct declaration.
type serverHTTPHandler struct{}

func (h *serverHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example http response")
}

// Summary:
// EchoServer echoes back any message sent over a websocket connection.
// Args:
// - ws (*websocket.Conn): a websocket connection object
// Returns:
// None. The function writes the message back to the client.
func EchoServer(ws *websocket.Conn) {
	if _, err := io.Copy(ws, ws); err != nil {
		panic(err)
	}
}

// Summary:
// The `ServeHTTP` method handles incoming HTTP requests and responds with an example HTTP response.
// Args:
// - `w http.ResponseWriter`: a response writer that is used to send an HTTP response to the client.
// - `r *http.Request`: a pointer to the incoming HTTP request.
// Returns:
// This method doesn't return anything.
func serveHTTP1(l net.Listener) {
	s := &http.Server{
		Handler: &serverHTTPHandler{},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

// Summary:
// This function creates a TLS listener and serves HTTP over it.
// Args:
// - `l` (net.Listener): The Listener object for the HTTPS connection.
// Returns:
// This function does not return anything.
func serveHTTPS(l net.Listener) {
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

	// Serve HTTP over TLS.
	serveHTTP1(tlsl)
}

// Summary:
// This function creates an HTTP server and configures it to handle WebSocket connections using the EchoServer function.
// Args:
// - l: A net.Listener instance that represents the listening socket.
// Returns:
// This function does not return a value. It serves the HTTP server until the listener is closed. If an error occurs, it panics.
func serveWS(l net.Listener) {
	s := &http.Server{
		Handler: websocket.Handler(EchoServer),
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

// Summary:
// This is a function called 'serveWSS' that serves a WebSocket connection using the input 'net.Listener' object.
// Args:
// - l: A 'net.Listener' object that represents the listener to use when serving the WebSocket connection.
// Returns:
// This function does not return any value. If an error occurs during the WebSocket server initialization and it is not related to
// the listener being closed, the function will panic.
func serveWSS(l net.Listener) {
	s := &http.Server{
		Handler: websocket.Handler(EchoServer),
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

// Summary:
//
//	tlsListener is a function that creates a TLS listener using the provided net.Listener.
//
// Args:
//   - l: A net.Listener representing the listener to be used for the TLS listener creation.
//
// Returns:
//   - A net.Listener representing the TLS listener.
func tlsListener(l net.Listener) net.Listener {
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
	return tlsl
}

// Summary:
//
// This is an example for serving HTTP and HTTPS on the same port.
// tlsListener is a function that creates a TLS listener using the provided net.Listener.
//
// Args:
// - l: A net.Listener representing the listener to be used for the TLS listener creation.
//
// Returns:
// - A net.Listener representing the TLS listener.
func ServeHTTPAndHTTPS(netprotocol string, ipaddressport string, enablehttps bool, enablews bool, enablewss bool) {
	if netprotocol == "" || netprotocol == "null" {
		netprotocol = "tcp"
	}
	if ipaddressport == "" || ipaddressport == "null" {
		ipaddressport = "127.0.0.1:50051"
	}
	// Create the TCP listener.
	l, err := net.Listen(netprotocol, ipaddressport)
	if err != nil {
		log.Panic(err)
	}

	// Create a mux.
	m := cmux.New(l)

	// We first match on HTTP 1.1 methods.
	httpl := m.Match(cmux.HTTP1Fast())
	go serveHTTP1(httpl)

	if enablehttps {
		//
		// If not matched, we assume that its TLS.
		// Note that you can take this listener, do TLS handshake and
		// create another mux to multiplex the connections over TLS.
		//
		tlsl := m.Match(cmux.Any())
		go serveHTTPS(tlsl)
	}

	if enablews {
		// // Otherwise, we match it againts a websocket upgrade request.
		wsl := m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))
		go serveWS(wsl)
	}

	if enablewss {
		// // Otherwise, we match it againts a websocket upgrade request.
		// wssl := m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))
		// go serveWSS(wssl)
	}

	// Summary:
	// Listen for the process signal to trigger grceful shutdown When an interrupt or termination signal is sent.
	// 		This function starts a goroutine for graceful shutdown of the server. It creates a channel to notify on receipt of interrupts.
	// 		It waits on the channel and when a signal is received it gracefully shuts down the server and runs cleanup tasks.
	// Args:
	// 		- l (net.Listener): an instance of net.Listener
	// Returns:
	// 		Void
	go func(l net.Listener) {
		// sig := make(chan os.Signal, 1)
		// signal.Notify(sig, syscall.SIGHUP)
		// for range sig {
		// 	upg.Upgrade()
		// }

		// // Create channel to signify a signal being sent
		c := make(chan os.Signal, 1)
		// // When an interrupt or termination signal is sent, notify the channel
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		// // This blocks the main thread until an interrupt is received
		_ = <-c
		fmt.Println("Gracefully shutting down.")
		_ = l.Close()
		fmt.Println("Running cleanup tasks.")

		// // Your cleanup tasks go here
		// db.Close()
		// redisConn.Close()
		fmt.Println("Fiber was successful shutdown.")
	}(l)

	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}

	// //
	// // Ways of Invoking
	// //
	// // Summary:
	// // 	This is a go anonymous function that launches a TLS server.
	// // Args:
	// // 	None
	// // Returns:
	// // 	None, it is a goroutine launched inside the main thread. However, it may panic if an error occurs while serving.
	// //
	//
	// go func() {
	// 	if err := tlsm.Serve(); err != cmux.ErrListenerClosed {
	// 		panic(err)
	// 	}
	// }()
	//
	// if err := tcpm.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
	// 		panic(err)
	// }
	//
}

// Summary:
// The `main()` function creates a TCP listener and serves a gRPC service on the specified address.
//
//	This is an example for serving HTTP, HTTPS, and GoRPC/TLS on the same port.
//
// Args:
// None
// Returns:
// None
func main() {
	// // Create the TCP listener.
	ServeHTTPAndHTTPS("tcp", "127.0.0.1:50051", true, true, false)
}
