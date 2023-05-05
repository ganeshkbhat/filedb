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
	"github.com/soheilhy/cmux"
	"golang.org/x/net/websocket"
)

// Summary:
// This is a struct declaration that defines a type named `serverHTTPHandler`.
// Args:
// There are no arguments for this struct declaration.
// Returns:
// There is no return statement for this struct declaration.
type serverHTTPHandler struct{}

// Summary:
// The `ServeHTTP` method handles incoming HTTP requests and responds with an example HTTP response.
// Args:
// - `w http.ResponseWriter`: a response writer that is used to send an HTTP response to the client.
// - `r *http.Request`: a pointer to the incoming HTTP request.
// Returns:
// This method doesn't return anything.
func (h *serverHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example http response")
}

// Summary:
// `serveHTTP` is a function that serves HTTP requests using the specified `net.Listener`.
// Args:
// - `l net.Listener`: an object of type `net.Listener` that represents the network listener.
// Returns:
// This function does not return any values, but it may panic if an error occurs during server operation.
func serveHTTP(l net.Listener) {
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

	// Serve HTTP over TLS.
	serveHTTP(tlsl)
}

// Summary:
// tlsListener is a function that creates a TLS listener using the provided net.Listener.
// Args:
// - l: A net.Listener representing the listener to be used for the TLS listener creation.
// Returns:
// - A net.Listener representing the TLS listener.
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

// // This is an example for serving HTTP, HTTPS, and GoRPC/TLS on the same port.
// Summary:
// Serve function starts a server for gRPC, HTTP, HTTPS, and WebSocket connections. It creates
//
//	a TCP listener, matches it against HTTP1.1 headers, WebSocket upgrade request, and any TLS connection. It also
//	listens for any interrupt or termination signal to trigger the graceful shutdown and cleanup tasks.
//
// Args:
// - netprotocol (string): The network protocol to be used, e.g. tcp, udp, etc.
// - ipaddressport (string): The IP address and port number for the listener to listen on.
//
// Returns:
// No return value.
func Serve(netprotocol string, ipaddressport string) {
	if netprotocol == "" || netprotocol == "null" {
		netprotocol = "tcp"
	}
	if ipaddressport == "" || ipaddressport == "null" {
		ipaddressport = "127.0.0.1:50051"
	}

	// // Create the TCP listener.
	l, err := net.Listen(netprotocol, ipaddressport)
	if err != nil {
		log.Panic(err)
	}

	// // Create a mux.
	tcpm := cmux.New(l)

	// // We first match on HTTP 1.1 methods.
	httpl := tcpm.Match(cmux.HTTP1Fast())

	// If not matched, we assume that its TLS.
	//
	// // Note that you can take this listener, do TLS handshake and
	// // create another mux to multiplex the connections over TLS.
	// // tlsl := m.Match(cmux.Any())
	tlsl := tcpm.Match(cmux.Any())
	tlsl = tlsListener(tlsl)

	// // // Otherwise, we match it againts a websocket upgrade request.
	wsl := tcpm.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))

	// // Now, we build another mux recursively to match HTTPS and GoRPC.
	// // You can use the same trick for SSH.
	// tlsm := cmux.New(tlsl)
	tlsm := cmux.New(l)

	// // // We first match the connection against HTTP2 fields. If matched, the
	// // // connection will be sent through the "grpcl" listener.
	// // grpcl := tlsm.Match(cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc"))
	// // // Otherwise, we match it againts a websocket upgrade request.
	wsl2 := tlsm.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))

	// // // Otherwise, we match it againts HTTP1 methods. If matched,
	// // // it is sent through the "httpl" listener.
	// // httpl := tcpm.Match(cmux.HTTP1Fast())
	// // // If not matched by HTTP, we assume it is an RPC connection.
	// // rpcl := tcpm.Match(cmux.Any())

	// // Then we used the muxed listeners.
	go serveWS(wsl)
	go serveWS(wsl2)

	go serveHTTP(httpl)
	go serveHTTPS(tlsl)

	// Summary:
	// Listen for the process signal to trigger grceful shutdown When an interrupt or termination signal is sent.
	// This function starts a goroutine for graceful shutdown of the server. It creates a channel to notify on receipt of interrupts.
	// It waits on the channel and when a signal is received it gracefully shuts down the server and runs cleanup tasks.
	// Args:
	// - l (net.Listener): an instance of net.Listener
	// Returns:
	// Void
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

	// Summary:
	// This is a go anonymous function that launches a TLS server.
	// Args:
	// None
	// Returns:
	// None, it is a goroutine launched inside the main thread. However, it may panic if an error occurs while serving.
	go func() {
		if err := tlsm.Serve(); err != cmux.ErrListenerClosed {
			panic(err)
		}
	}()

	if err := tcpm.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}

}

// Summary:
// The `main()` function creates a TCP listener and serves a gRPC service on the specified address.
//	This is an example for serving HTTP, HTTPS, and GoRPC/TLS on the same port.
//
// Args:
// None
// Returns:
// None
func main() {
	// // Create the TCP listener.
	Serve("tcp", "127.0.0.1:50051")
}
