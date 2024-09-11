package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/LetsFocus/goLF/errors"

	"github.com/LetsFocus/template-service/constants"
	"github.com/LetsFocus/template-service/models"
	"github.com/LetsFocus/template-service/services"
)

type Handler struct {
	services.Invoices
}

func New(invoice services.Invoices) *Handler {
	return &Handler{Invoices: invoice}
}

func (h *Handler) Create(ctx *gin.Context) {
	var invoiceTemplate models.Template

	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	err = ctx.BindJSON(&invoiceTemplate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	invoiceTemplate.TenantID = tenantID

	resp, err := h.Invoices.Create(ctx, &invoiceTemplate)
	if err != nil {
		ctx.JSON(parseError(err))

		return
	}

	ctx.JSON(http.StatusOK, responseHandler(resp))
}

func (h *Handler) Get(ctx *gin.Context) {
	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	service := strings.ToLower(ctx.Query("service"))
	if strings.TrimSpace(service) == "" {
		ctx.JSON(http.StatusBadRequest, errorHandler(errors.MissingParam([]string{"service"})))

		return
	}

	//universal, err := strconv.ParseBool(strings.TrimSpace(ctx.Query("universal")))
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, errors.MissingParam([]string{"universal"}))
	//
	//	return
	//}

	pageSize, err := strconv.Atoi(ctx.Query(constants.PageSize))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	pageNumber, err := strconv.Atoi(ctx.Query(constants.PageNumber))
	if err != nil || pageNumber <= 0 {
		pageNumber = 1
	}

	searchKey := ctx.Query(constants.SearchKey)

	p := models.Pagination{PageNumber: pageNumber, PageSize: pageSize}

	f := models.Filters{Service: service, Pagination: p, SearchKey: searchKey}

	resp, pagination, err := h.Invoices.Get(ctx, tenantID, f)
	if err != nil {
		ctx.JSON(parseError(err))

		return
	}

	ctx.JSON(http.StatusOK, getAllResponseHandler(resp, pagination))
}

func (h *Handler) GetByID(ctx *gin.Context) {
	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	ID, err := models.ValidateUUID(ctx.Param(constants.ID), constants.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	resp, err := h.Invoices.GetByID(ctx, tenantID, ID)
	if err != nil {
		ctx.JSON(parseError(err))

		return
	}

	ctx.JSON(http.StatusOK, responseHandler(resp))
}

func (h *Handler) Delete(ctx *gin.Context) {
	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	ID, err := models.ValidateUUID(ctx.Param(constants.ID), constants.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	resp, err := h.Invoices.Delete(ctx, tenantID, ID)
	if err != nil {
		ctx.JSON(parseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, responseHandler(resp))
}

func (h *Handler) Patch(ctx *gin.Context) {
	var invoiceTemplate models.Template

	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	ID, err := models.ValidateUUID(ctx.Param(constants.ID), constants.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	err = ctx.BindJSON(&invoiceTemplate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))

		return
	}

	invoiceTemplate.ID = ID
	invoiceTemplate.TenantID = tenantID

	resp, err := h.Invoices.Patch(ctx, &invoiceTemplate)
	if err != nil {
		ctx.JSON(parseError(err))

		return
	}

	ctx.JSON(http.StatusOK, responseHandler(resp))
}

func parseError(err error) (int, interface{}) {
	if parsedErr, ok := err.(errors.Errors); ok {
		return parsedErr.StatusCode, errorHandler(err)
	}

	return http.StatusInternalServerError, errorHandler(err)
}

func responseHandler(data interface{}) interface{} {
	return map[string]interface{}{
		"data": data,
	}
}

func getAllResponseHandler(data interface{}, pagination models.Pagination) interface{} {
	return map[string]interface{}{
		"data":      data,
		"pagination": pagination,
	}
}

func errorHandler(err interface{}) interface{} {
	return map[string]interface{}{
		"errors": err,
	}
}
