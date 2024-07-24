package stores

import (
	"database/sql"
	"github.com/LetsFocus/goLF/errors"
	"github.com/LetsFocus/template-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) Create(ctx *gin.Context, invoice *models.Template) (models.Template, error) {
	rows, err := s.DB.ExecContext(ctx, CREATEQUERY,
		invoice.TenantID, invoice.ID, invoice.Name, invoice.Description, invoice.Content, invoice.Service,
		invoice.Universal, invoice.CreatedAt, invoice.UpdatedAt)
	if err != nil {
		return models.Template{}, errors.Errors{StatusCode: http.StatusInternalServerError, Reason: err.Error(), Code: "DB Error"}
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return models.Template{}, errors.RowsAffectedError(err)
	}

	if rowsAffected == 0 {
		return models.Template{}, errors.Errors{StatusCode: http.StatusInternalServerError,
			Code: http.StatusText(http.StatusInternalServerError), Reason: "Data not inserted"}
	}

	return *invoice, nil
}

func (s *Store) Get(ctx *gin.Context, tenantId uuid.UUID, f models.Filters) ([]models.Template, error) {
	// TODO: Consider TenantID as well
	rows, err := s.DB.QueryContext(ctx, GETQUERY, f.Service)
	if err != nil {
		return nil, errors.Errors{StatusCode: http.StatusInternalServerError, Reason: err.Error(), Code: "DB Error"}
	}

	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var template models.Template

		err = rows.Scan(&template.TenantID, &template.ID, &template.Name, &template.Description, &template.Content,
			&template.Service, &template.Universal, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, errors.InternalServerError(err)
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func (s *Store) Patch(ctx *gin.Context, invoice *models.Template) (models.Template, error) {
	dynamicQuery, values := generateQuery(invoice)

	query, values := buildQuery(dynamicQuery, values, invoice.TenantID, invoice.ID)
	row := s.DB.QueryRowContext(ctx, query, values...)

	var updatedTemplate models.Template
	err := row.Scan(
		&updatedTemplate.TenantID, &updatedTemplate.ID, &updatedTemplate.Name, &updatedTemplate.Description,
		&updatedTemplate.Content, &updatedTemplate.Service, &updatedTemplate.Universal,
		&updatedTemplate.CreatedAt, &updatedTemplate.UpdatedAt,
	)
	if err != nil {
		return models.Template{}, errors.InternalServerError(err)
	}

	return updatedTemplate, nil
}

func generateQuery(invoice *models.Template) (string, []interface{}) {
	var updateQuery []string
	var updateValues []interface{}

	addField := func(field string, value interface{}) {
		updateQuery = append(updateQuery, field+"=$"+strconv.Itoa(len(updateValues)+1))
		updateValues = append(updateValues, value)
	}

	if invoice.Name != "" {
		addField("name", invoice.Name)
	}
	if invoice.Description != "" {
		addField("description", invoice.Description)
	}
	if invoice.Content != "" {
		addField("content", invoice.Content)
	}

	query := strings.Join(updateQuery, ", ")
	return query, updateValues
}

func buildQuery(dynamicQuery string, values []interface{}, tenantID, id uuid.UUID) (string, []interface{}) {
	query := "UPDATE templates SET " + dynamicQuery + ", updated_at = $" + strconv.Itoa(len(values)+1) +
		" WHERE tenant_id = $" + strconv.Itoa(len(values)+2) + " AND id = $" + strconv.Itoa(len(values)+3) + " Returning tenant_id, id, name, description, content, service, universal, created_at, updated_at"
	values = append(values, time.Now(), tenantID, id)
	return query, values
}
func (s *Store) GetByID(ctx *gin.Context, tenantId, id uuid.UUID) (models.Template, error) {
	var template models.Template

	rows := s.DB.QueryRowContext(ctx, GETBYIDQUERY, id, tenantId)

	err := rows.Scan(&template.TenantID, &template.ID, &template.Name, &template.Description, &template.Content,
		&template.Service, &template.Universal, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return models.Template{}, errors.InternalServerError(err)
	}

	return template, nil
}

// TODO: refactor this to not to return anything in the response
func (s *Store) Delete(ctx *gin.Context, tenantId, id uuid.UUID) (models.Template, error) {
	rows, err := s.DB.ExecContext(ctx, DELETEQUERY, id, tenantId)
	if err != nil {
		return models.Template{}, errors.Errors{StatusCode: http.StatusInternalServerError, Reason: err.Error(), Code: "DB Error"}
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return models.Template{}, errors.RowsAffectedError(err)
	}

	if rowsAffected == 0 {
		return models.Template{}, errors.Errors{StatusCode: http.StatusInternalServerError,
			Code: http.StatusText(http.StatusInternalServerError), Reason: "Data not Deleted"}
	}

	return models.Template{ID: id, TenantID: tenantId}, nil
}
