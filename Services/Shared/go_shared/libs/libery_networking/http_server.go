package libery_networking

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

type ServerConfig struct {
	Port        string
	ServiceName string
}

type Server interface {
	Config() *ServerConfig
}

type Broker struct {
	config                  *ServerConfig
	router                  *patriot_router.Router
	http_server             *http.Server
	grpc_server             GrpcServer
	after_shutdown_callback func()
}

func (broker *Broker) Config() *ServerConfig {
	return broker.config
}

func CorsAllowAll(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, AuthType")
		handler(response, request)
	}
}

func (broker *Broker) AwaitGracefulShutdown() (http_err, grpc_err error) {
	quite := make(chan os.Signal, 1)

	signal.Notify(quite, os.Interrupt, os.Kill, syscall.SIGTERM)

	signal_received := <-quite

	echo.Echo(echo.PurpleBG, fmt.Sprintf("Signal received: %s", signal_received.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := broker.Shutdown(ctx); err != nil {
		http_err = err
	}

	if err := broker.ShutdownGrpcServer(); err != nil {
		grpc_err = err
	}

	return
}

func (broker *Broker) callAfterShutdown() {
	if broker.after_shutdown_callback == nil {
		return
	}

	broker.after_shutdown_callback()
}

func (broker *Broker) Run(binder func(server Server, router *patriot_router.Router)) {
	binder(broker, broker.router)	

	echo.Echo(echo.GreenBG, fmt.Sprintf("Starting %s on port %s", broker.config.ServiceName, broker.config.Port))
	if err := broker.http_server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		echo.EchoFatal(err)
	} else {
		echo.Echo(echo.GreenFG, "Server stopped")
	}
}

func (broker *Broker) RunGrpcServer() {
	echo.Echo(echo.PinkBG, "Starting GRPC Servers")
	go broker.grpc_server.Connect()
}

func (broker *Broker) Shutdown(ctx context.Context) (err error) {
	err = broker.http_server.Shutdown(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (broker *Broker) ShutdownGrpcServer() error {
	if broker.grpc_server == nil {
		return nil
	}

	err := broker.grpc_server.Shutdown()

	return err
}

func (broker *Broker) SetGrpcServer(grpc_server GrpcServer) {
	broker.grpc_server = grpc_server
}

// Sets a function to be called after StartServer shuts down the http and grpc servers
func (broker *Broker) OnAfterShutdown(callback func()) {
	broker.after_shutdown_callback = callback
}

// Starts the http server and the grpc server if it was set. await for a signal to shutdown and call the after_shutdown_callback
func (broker *Broker) StartServer(binder func(server Server, router *patriot_router.Router)) {
	go broker.Run(binder)

	if broker.grpc_server != nil {
		go broker.RunGrpcServer()
	}

	http_err, grpc_err := broker.AwaitGracefulShutdown()

	if http_err != nil {
		echo.EchoErr(http_err)
	}

	if grpc_err != nil {
		echo.EchoErr(grpc_err)
	}

	broker.callAfterShutdown()
}

func NewBroker(ctx context.Context, config *ServerConfig) (*Broker, error) {
	if config.Port == "" {
		return nil, fmt.Errorf("port is required")
	}

	var new_broker *Broker = new(Broker)
	new_broker.config = config
	new_broker.router = patriot_router.CreateRouter()
	new_broker.router.SetCorsHandler(CorsAllowAll)

	new_broker.http_server = &http.Server{
		Addr:    config.Port,
		Handler: new_broker.router,
	}

	return new_broker, nil
}
