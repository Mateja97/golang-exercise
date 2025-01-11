package api

import (
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"golang-exercise/writer"

	"google.golang.org/grpc"
)

type handler struct {
	companyService.UnimplementedCompanyServiceServer
	writer *writer.Writer
}

var h *handler

func Init(srv grpc.ServiceRegistrar, opts ...func(*handler)) {
	h = &handler{}
	for _, o := range opts {
		o(h)
	}
	companyService.RegisterCompanyServiceServer(srv, h)
}
