package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"

	"github.com/josephlbailey/alert-service/internal/db/domain"
)

type CreateAlertReq struct {
	Message string `json:"message" binding:"required"`
}

type UpdateAlertReq struct {
	Message string `json:"message" binding:"required"`
}

type AlertRes struct {
	ExternalID uuid.UUID `json:"externalId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Message    string    `json:"message"`
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (req *CreateAlertReq) Bind(c *gin.Context, p *domain.CreateAlertParams) error {
	if err := c.ShouldBindJSON(req); err != nil {
		var (
			ve validator.ValidationErrors
			je *json.UnmarshalTypeError
		)
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		} else if errors.As(err, &je) {
			out := make([]ErrorMsg, 1)
			out[0] = ErrorMsg{je.Field, "invalid type for field"}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return err
	}

	p.ExternalID = uuid.Must(uuid.NewV4())
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Message = req.Message
	return nil
}

func (req *UpdateAlertReq) Bind(c *gin.Context, p *domain.UpdateAlertByIDParams) error {
	if err := c.ShouldBindJSON(req); err != nil {
		var (
			ve validator.ValidationErrors
			je *json.UnmarshalTypeError
		)
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		} else if errors.As(err, &je) {
			out := make([]ErrorMsg, 1)
			out[0] = ErrorMsg{je.Field, "invalid type for field"}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return err
	}

	p.UpdatedAt = time.Now()
	p.Message = req.Message
	return nil
}

func NewAlertResponse(alert *domain.Alert) *AlertRes {
	resp := new(AlertRes)
	resp.CreatedAt = alert.CreatedAt
	resp.UpdatedAt = alert.UpdatedAt
	resp.ExternalID = alert.ExternalID
	resp.Message = alert.Message
	return resp
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "gte":
		return "should be greater than " + fe.Param()
	}
	return "unknown error"
}
