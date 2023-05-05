package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	grpc "google.golang.org/grpc"
	// codes "google.golang.org/grpc/codes"
	// status "google.golang.org/grpc/status"

	// "golang.org/x/net/context"
	"golang.org/x/net/websocket"

	"github.com/soheilhy/cmux"
)


// // GRPC SERVR

// // The request message containing the user's name.
// type HelloRequest struct {
// 	state         protoimpl.MessageState
// 	sizeCache     protoimpl.SizeCache
// 	unknownFields protoimpl.UnknownFields

// 	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
// }

// func (x *HelloRequest) Reset() {
// 	*x = HelloRequest{}
// 	if protoimpl.UnsafeEnabled {
// 		mi := &file_examples_helloworld_helloworld_helloworld_proto_msgTypes[0]
// 		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
// 		ms.StoreMessageInfo(mi)
// 	}
// }

// func (x *HelloRequest) String() string {
// 	return protoimpl.X.MessageStringOf(x)
// }

// func (*HelloRequest) ProtoMessage() {}

// func (x *HelloRequest) ProtoReflect() protoreflect.Message {
// 	mi := &file_examples_helloworld_helloworld_helloworld_proto_msgTypes[0]
// 	if protoimpl.UnsafeEnabled && x != nil {
// 		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
// 		if ms.LoadMessageInfo() == nil {
// 			ms.StoreMessageInfo(mi)
// 		}
// 		return ms
// 	}
// 	return mi.MessageOf(x)
// }

// // Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
// func (*HelloRequest) Descriptor() ([]byte, []int) {
// 	return file_examples_helloworld_helloworld_helloworld_proto_rawDescGZIP(), []int{0}
// }

// func (x *HelloRequest) GetName() string {
// 	if x != nil {
// 		return x.Name
// 	}
// 	return ""
// }

// // The response message containing the greetings
// type HelloReply struct {
// 	state         protoimpl.MessageState
// 	sizeCache     protoimpl.SizeCache
// 	unknownFields protoimpl.UnknownFields

// 	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
// }

// func (x *HelloReply) Reset() {
// 	*x = HelloReply{}
// 	if protoimpl.UnsafeEnabled {
// 		mi := &file_examples_helloworld_helloworld_helloworld_proto_msgTypes[1]
// 		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
// 		ms.StoreMessageInfo(mi)
// 	}
// }

// func (x *HelloReply) String() string {
// 	return protoimpl.X.MessageStringOf(x)
// }

// func (*HelloReply) ProtoMessage() {}

// func (x *HelloReply) ProtoReflect() protoreflect.Message {
// 	mi := &file_examples_helloworld_helloworld_helloworld_proto_msgTypes[1]
// 	if protoimpl.UnsafeEnabled && x != nil {
// 		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
// 		if ms.LoadMessageInfo() == nil {
// 			ms.StoreMessageInfo(mi)
// 		}
// 		return ms
// 	}
// 	return mi.MessageOf(x)
// }

// // Deprecated: Use HelloReply.ProtoReflect.Descriptor instead.
// func (*HelloReply) Descriptor() ([]byte, []int) {
// 	return file_examples_helloworld_helloworld_helloworld_proto_rawDescGZIP(), []int{1}
// }

// // 
// func (x *HelloReply) GetMessage() string {
// 	if x != nil {
// 		return x.Message
// 	}
// 	return ""
// }


type serverHTTPHandler struct{}
type serverRPCRcvr struct{}

func (h *serverHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example http response")
}

func serveRpc() {
	// Create the main listener.
	l, err := net.Listen("tcp", ":23456")
	if err != nil {
		log.Fatal(err)
	}
	// Create a cmux.
	m := cmux.New(l)

	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	trpcL := m.Match(cmux.Any()) // Any means anything that is not yet matched.

	// Create your protocol servers.
	trpcS := rpc.NewServer()
	trpcS.Register(&serverRPCRcvr{})

	// Use the muxed listeners for your servers.
	go trpcS.Accept(trpcL)

	// Start serving!
	m.Serve()
}

func serveHTTP(l net.Listener) {
	s := &http.Server{
		Handler: &serverHTTPHandler{},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

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

type recursiveHTTPHandler struct{}

func (h *recursiveHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "example http response")
}

func recursiveServeHTTP(l net.Listener) {
	s := &http.Server{
		Handler: &recursiveHTTPHandler{},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

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

type RecursiveRPCRcvr struct{}

func (r *RecursiveRPCRcvr) Cube(i int, j *int) error {
	*j = i * i
	return nil
}

func recursiveServeRPC(l net.Listener) {
	s := rpc.NewServer()
	if err := s.Register(&RecursiveRPCRcvr{}); err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			if err != cmux.ErrListenerClosed {
				panic(err)
			}
			return
		}
		go s.ServeConn(conn)
	}
}

func EchoServer(ws *websocket.Conn) {
	if _, err := io.Copy(ws, ws); err != nil {
		panic(err)
	}
}

func serveWS(l net.Listener) {
	s := &http.Server{
		Handler: websocket.Handler(EchoServer),
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

type ServerRPCRcvr struct{}

func (r *ServerRPCRcvr) Cube(i int, j *int) error {
	*j = i * i
	return nil
}

// 
// func serveRPC(l net.Listener) {
// 	s := rpc.NewServer()
// 	if err := s.Register(&ServerRPCRcvr{}); err != nil {
// 		panic(err)
// 	}
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			if err != cmux.ErrListenerClosed {
// 				panic(err)
// 			}
// 			return
// 		}
// 		go s.ServeConn(conn)
// 	}
// }
// 

// UnimplementedGreeterServer must be embedded to have forward compatible implementations.
type UnimplementedGreeterServer struct {
}

type grpcServer struct {
	UnimplementedGreeterServer
}

// func (s *grpcServer) SayHello(ctx context.Context, in *grpchello.HelloRequest) (*grpchello.HelloReply, error) {
// 	return &grpchello.HelloReply{Message: "Hello " + in.Name + " from cmux"}, nil
// }

// 
// func serveGRPC(l net.Listener) {
// 	grpcs := grpc.NewServer()
// 	grpchello.RegisterGreeterServer(grpcs, &grpcServer{})
// 	if err := grpcs.Serve(l); err != cmux.ErrListenerClosed {
// 		panic(err)
// 	}
// }
//

// // This is an example for serving HTTP, HTTPS, and GoRPC/TLS on the same port.
func main() {
	// // Create the TCP listener.
	l, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Panic(err)
	}

	// // Create a mux.
	tcpm := cmux.New(l)

	// // We first match on HTTP 1.1 methods.
	httpl := tcpm.Match(cmux.HTTP1Fast())

	// // If not matched, we assume that its TLS.
	
	// If not matched, we assume that its TLS.
	//
	// // Note that you can take this listener, do TLS handshake and
	// // create another mux to multiplex the connections over TLS.
	// // tlsl := m.Match(cmux.Any())
	tlsl := tcpm.Match(cmux.Any())
	tlsl = tlsListener(tlsl)

	// // Now, we build another mux recursively to match HTTPS and GoRPC.
	// // You can use the same trick for SSH.
	tlsm := cmux.New(tlsl)

	// // We first match on HTTP 1.1 methods.
	// // httpl := m.Match(cmux.HTTP1Fast())
	httpsl := tlsm.Match(cmux.HTTP1Fast())
	gorpcl := tlsm.Match(cmux.Any())

	// // We first match the connection against HTTP2 fields. If matched, the
	// // connection will be sent through the "grpcl" listener.
	// grpcl := tlsm.Match(cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc"))
	// // Otherwise, we match it againts a websocket upgrade request.
	wsl := tcpm.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))

	// // Otherwise, we match it againts HTTP1 methods. If matched,
	// // it is sent through the "httpl" listener.
	// httpl := tcpm.Match(cmux.HTTP1Fast())
	// // If not matched by HTTP, we assume it is an RPC connection.
	// rpcl := tcpm.Match(cmux.Any())

	go recursiveServeHTTP(httpl)
	go recursiveServeHTTP(httpsl)
	go recursiveServeRPC(gorpcl)

	// // Then we used the muxed listeners.
	// go serveGRPC(grpcl)
	go serveWS(wsl)
	// go serveRPC(rpcl)

	go serveHTTP(httpl)
	go serveHTTPS(tlsl)

	// Listen for the process signal to trigger grceful shutdown When an interrupt or termination signal is sent
	go func() {
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
	}()

	go func() {
		if err := tlsm.Serve(); err != cmux.ErrListenerClosed {
			panic(err)
		}
	}()

	if err := tcpm.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}
