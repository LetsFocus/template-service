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
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	err = ctx.BindJSON(&invoiceTemplate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	invoiceTemplate.TenantID = tenantID

	resp, err := h.Invoices.Create(ctx, &invoiceTemplate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) Get(ctx *gin.Context) {
	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	service := ctx.Query("service")
	if strings.TrimSpace(service) == "" {
		ctx.JSON(http.StatusBadRequest, errors.MissingParam([]string{"service"}))

		return
	}

	universal, err := strconv.ParseBool(strings.TrimSpace(ctx.Query("universal")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.MissingParam([]string{"universal"}))

		return
	}

	f := models.Filters{Service: service, Universal: universal}

	resp, err := h.Invoices.Get(ctx, tenantID, f)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(ctx *gin.Context) {
	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ID, err := models.ValidateUUID(ctx.Param(constants.ID), constants.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	resp, err := h.Invoices.GetByID(ctx, tenantID, ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) Delete(ctx *gin.Context) {
	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ID, err := models.ValidateUUID(ctx.Param(constants.ID), constants.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	resp, err := h.Invoices.Delete(ctx, tenantID, ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ctx.JSON(http.StatusNoContent, resp)
}

func (h *Handler) Patch(ctx *gin.Context) {
	var invoiceTemplate models.Template

	tenantID, err := models.ValidateUUID(ctx.GetHeader(constants.TenantID), constants.TenantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ID, err := models.ValidateUUID(ctx.Param(constants.ID), constants.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	err = ctx.BindJSON(&invoiceTemplate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	invoiceTemplate.ID = ID
	invoiceTemplate.TenantID = tenantID

	resp, err := h.Invoices.Patch(ctx, &invoiceTemplate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	ctx.JSON(http.StatusOK, resp)
}
