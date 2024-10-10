package libery_networking

type GrpcServer interface {
	Connect() error
	Shutdown() error
}
