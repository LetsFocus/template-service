package services

import (
	"github.com/LetsFocus/template-service/models"
	"github.com/LetsFocus/template-service/stores"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	stores.Invoices
}

func New(invoices stores.Invoices) *Service {
	return &Service{Invoices: invoices}
}

func (s *Service) Create(ctx *gin.Context, invoice *models.Template) (models.Template, error) {
	err := invoice.Validate()
	if err != nil {
		return models.Template{}, err
	}

	invoice.Universal = false

	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	invoice.ID = uuid.New()

	resp, err := s.Invoices.Create(ctx, invoice)
	if err != nil {
		return models.Template{}, err
	}

	return resp, nil
}

func (s *Service) Get(ctx *gin.Context, tenantId uuid.UUID, f models.Filters) ([]models.Template, models.Pagination, error) {
	f.Offset = (f.PageNumber - 1) * f.PageSize
	f.Limit = f.PageSize
	resp, pagination, err := s.Invoices.Get(ctx, tenantId, f)
	if err != nil {
		return nil, pagination, err
	}

	return resp, pagination, nil
}

func (s *Service) Patch(ctx *gin.Context, invoice *models.Template) (models.Template, error) {
	err := invoice.ValidatePatch()
	if err != nil {
		return models.Template{}, err
	}
	resp, err := s.Invoices.Patch(ctx, invoice)
	if err != nil {
		return models.Template{}, err
	}
	return resp, nil
}

func (s *Service) GetByID(ctx *gin.Context, tenantId, id uuid.UUID) (models.Template, error) {
	resp, err := s.Invoices.GetByID(ctx, tenantId, id)
	if err != nil {
		return models.Template{}, err
	}

	return resp, nil
}

func (s *Service) Delete(ctx *gin.Context, tenantId, id uuid.UUID) (models.Template, error) {
	resp, err := s.Invoices.Delete(ctx, tenantId, id)
	if err != nil {
		return models.Template{}, err
	}

	return resp, nil
}
