package api

import (
	"context"
	"golang-exercise/api/codes"
	"log"
	"net"
	"testing"

	companyMessage "github.com/mateja97/golang-exercise/protobuf/golang/client/message/company/v1"
	companyService "github.com/mateja97/golang-exercise/protobuf/golang/client/service/company/v1"
	"google.golang.org/grpc"
)

type MockCompanyService struct {
	companyService.UnimplementedCompanyServiceServer
	companies map[string]*companyMessage.Company
}

func NewMockCompanyService() *MockCompanyService {
	return &MockCompanyService{
		companies: make(map[string]*companyMessage.Company),
	}
}

func (s *MockCompanyService) Create(ctx context.Context, req *companyService.CreateRequest) (*companyService.CreateResponse, error) {
	id := "company-id-1" // Generate or mock ID
	company := &companyMessage.Company{
		Id:                id,
		Name:              req.Name,
		Description:       req.Description,
		AmountOfEmployees: req.AmountOfEmployees,
		Registered:        req.Registered,
		Type:              req.Type,
	}
	s.companies[id] = company
	return &companyService.CreateResponse{Company: company}, nil
}

func (s *MockCompanyService) Patch(ctx context.Context, req *companyService.PatchRequest) (*companyService.PatchResponse, error) {
	company, exists := s.companies[req.Id]
	if !exists {
		return nil, codes.ErrNotFound
	}

	if req.Description != nil {
		company.Description = req.GetDescription()
	}
	if req.AmountOfEmployees != nil {
		company.AmountOfEmployees = req.GetAmountOfEmployees()
	}
	if req.Registered != nil {
		company.Registered = req.GetRegistered()
	}
	if req.Type != nil {
		company.Type = req.GetType()
	}
	return &companyService.PatchResponse{Company: company}, nil
}

func (s *MockCompanyService) Delete(ctx context.Context, req *companyService.DeleteRequest) (*companyService.DeleteResponse, error) {
	delete(s.companies, req.Id)
	return &companyService.DeleteResponse{}, nil
}

func (s *MockCompanyService) Get(ctx context.Context, req *companyService.GetRequest) (*companyService.GetResponse, error) {
	company, exists := s.companies[req.Id]
	if !exists {
		return nil, codes.ErrNotFound
	}
	return &companyService.GetResponse{Company: company}, nil
}

func TestCompanyService(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	server := grpc.NewServer()
	mockService := NewMockCompanyService()
	companyService.RegisterCompanyServiceServer(server, mockService)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	defer server.Stop()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := companyService.NewCompanyServiceClient(conn)

	t.Run("Create", func(t *testing.T) {
		req := &companyService.CreateRequest{
			Name:              "Test Company",
			Description:       "A test company",
			AmountOfEmployees: 100,
			Registered:        true,
			Type:              companyMessage.CompanyType_COMPANY_TYPE_CORPORATIONS,
		}
		resp, err := client.Create(context.Background(), req)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if resp.Company.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Company.Name)
		}
	})

	t.Run("Patch", func(t *testing.T) {
		description := "Updated description"
		amountOfEmployees := uint32(120)
		req := &companyService.PatchRequest{
			Id:                "company-id-1",
			Description:       &description,
			AmountOfEmployees: &amountOfEmployees,
		}
		resp, err := client.Patch(context.Background(), req)
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}
		if resp.Company.Description != "Updated description" {
			t.Errorf("Expected description 'Updated description', got %s", resp.Company.Description)
		}
	})

	t.Run("Get", func(t *testing.T) {
		req := &companyService.GetRequest{Id: "company-id-1"}
		resp, err := client.Get(context.Background(), req)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if resp.Company.Id != "company-id-1" {
			t.Errorf("Expected ID 'company-id-1', got %s", resp.Company.Id)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		req := &companyService.DeleteRequest{Id: "company-id-1"}
		_, err := client.Delete(context.Background(), req)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
	})
}
