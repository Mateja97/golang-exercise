package api

import (
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"

	"google.golang.org/grpc"
)

type handler struct {
	companyService.UnimplementedCompanyServiceServer
}

var h *handler

func Init(srv grpc.ServiceRegistrar, opts ...func(*handler)) {
	h = &handler{}
	for _, o := range opts {
		o(h)
	}
	companyService.RegisterCompanyServiceServer(srv, h)
}
