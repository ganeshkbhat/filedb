package netgrpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	grpc "google.golang.org/grpc"
)

// Setup gRPC server.
type grpcServer struct{}

func (s *grpcServer) Query(ctx context.Context, req *graph.Request) (*graph.Response, error) {
	return
}

// Handler function for http/https queries.

func queryHandler(w http.ResponseWriter, r *http.Request) {
	addCorsHeaders(w)

}

// Wrapper functions to start serving different services.
func serveGRPC(l net.Listener) {
	s := grpc.NewServer(grpc.CustomCodec(&query.Codec{}))
	graph.RegisterDgraphServer(s, &grpcServer{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("While serving gRpc request: %v", err)
	}
}

func serveHTTP(l net.Listener) {
	if err := http.Serve(l, nil); err != nil {
		log.Fatalf("While serving http request: %v", err)
	}
}

// // wrong struct and mod implementation

// Summary:
// This function creates a main listener for TCP on port 23456. It then creates a cmux to match connections in order, first for gRPC, then
// HTTP, and otherwise for Go RPC/TCP. The function then creates protocol servers and uses muxed listeners for these servers before serving.
// Args:
// This function does not take in any arguments.
// Returns:
// This function does not have a return statement.
//
// func serveRpc() {
// 	// Create the main listener.
// 	l, err := net.Listen("tcp", ":23456")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
//  // Create a cmux.
// 	m := cmux.New(l)
//
// 	// Match connections in order:
// 	// First grpc, then HTTP, and otherwise Go RPC/TCP.
//  // Any means anything that is not yet matched.
// 	trpcL := m.Match(cmux.Any())
// 	// Create your protocol servers.
//
// 	trpcS := grpc.NewServer()
// 	trpcS.Register(&serverRPCRcvr{})
//
// 	// Use the muxed listeners for your servers.
// 	go trpcS.Accept(trpcL)
// 	// Start serving!
// 	m.Serve()
//
// }

func setupServer() {
	// For internal communication.
	go worker.RunServer(*workerPort)
	// Create a listener at the desired port.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux object.
	tcpm := cmux.New(l)

	// Declare the match for different services required.
	httpl := tcpm.Match(cmux.HTTP1Fast())
	grpcl := tcpm.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	http2 := tcpm.Match(cmux.HTTP2())

	// Link the endpoint to the handler function.
	http.HandleFunc("/query", queryHandler)

	// Initialize the servers by passing in the custom listeners (sub-listeners).
	go serveGRPC(grpcl)
	go serveHTTP(httpl)
	go serveHTTP(http2)

	// Close the listener when done.
	go func() {
		<-closeCh
		// Stops listening further but already accepted connections are not closed.
		l.Close()
	}()

	log.Println("grpc server started.")
	log.Println("http server started.")
	log.Println("Server listening on port", *port)

	// Start cmux serving.
	if err := tcpm.Serve(); !strings.Contains(err.Error(),
		"use of closed network connection") {
		log.Fatal(err)
	}
}
