package services

import (
	"github.com/LetsFocus/template-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Invoices interface {
	Create(ctx *gin.Context, invoice *models.Template) (models.Template, error)
	Get(ctx *gin.Context, tenantId uuid.UUID, f models.Filters) ([]models.Template, models.Pagination, error)
	Patch(ctx *gin.Context, invoice *models.Template) (models.Template, error)
	GetByID(ctx *gin.Context, tenantId, id uuid.UUID) (models.Template, error)
	Delete(ctx *gin.Context, tenantId, id uuid.UUID) (models.Template, error)
}
