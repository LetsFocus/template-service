package models

import (
	"github.com/LetsFocus/goLF/errors"
	"github.com/LetsFocus/template-service/constants"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type Template struct {
	TenantID    uuid.UUID `json:"tenantId"`
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Service     string    `json:"service"`
	Universal   bool      `json:"universal"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Filters struct {
	Universal bool
	Service   string
	Limit     int
	Offset    int
	Pagination
	SearchKey string
}

type Pagination struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
	Count      int `json:"count"`
}

func (t *Template) Validate() error {
	if t.Name == "" {
		return errors.InvalidParam([]string{constants.Name})
	}

	if t.Description == "" {
		return errors.InvalidParam([]string{constants.Description})
	}

	if t.Content == "" {
		return errors.InvalidParam([]string{constants.Content})
	}

	if t.Service == "" {
		return errors.InvalidParam([]string{"service"})
	}
	t.Service = strings.ToLower(t.Service)

	return nil
}

func (t *Template) ValidatePatch() error {
	var invalidParams []string

	if t.Name == "" {
		invalidParams = append(invalidParams, "Name")
	}
	if t.Description == "" {
		invalidParams = append(invalidParams, "Description")
	}
	if t.Content == "" {
		invalidParams = append(invalidParams, "Content")
	}

	if len(invalidParams) == 3 {
		return errors.Errors{StatusCode: http.StatusBadRequest, Reason: "One of name , description or content is required", Code: http.StatusText(http.StatusBadRequest)}
	}

	return nil
}
