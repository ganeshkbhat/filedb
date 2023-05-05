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

type serverHTTPHandler struct{}
type serverRPCRcvr struct{}

func (h *serverHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example http response")
}

// wrong struct and mod implementation
// func serveRpc() {
// 	// Create the main listener.
// 	l, err := net.Listen("tcp", ":23456")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Create a cmux.
// 	m := cmux.New(l)

// 	// Match connections in order:
// 	// First grpc, then HTTP, and otherwise Go RPC/TCP.
// 	trpcL := m.Match(cmux.Any()) // Any means anything that is not yet matched.

// 	// Create your protocol servers.
// 	trpcS := grpc.NewServer()
// 	trpcS.Register(&serverRPCRcvr{})

// 	// Use the muxed listeners for your servers.
// 	go trpcS.Accept(trpcL)

// 	// Start serving!
// 	m.Serve()
// }

func serveHTTP(l net.Listener) {
	s := &http.Server{
		Handler: &serverHTTPHandler{},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

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
			// log.Fatal(app.Listener(ln))
			log.Fatal(ln)
		}()
	}

	// Serve HTTP over TLS.
	serveHTTP(tlsl)
}

// This is an example for serving HTTP and HTTPS on the same port.
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

	// // Create a cmux object.
	// tcpm := cmux.New(l)

	// // Declare the match for different services required.
	// httpl := tcpm.Match(cmux.HTTP1Fast())
	// grpcl := tcpm.MatchWithWriters(
	// 	cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
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
