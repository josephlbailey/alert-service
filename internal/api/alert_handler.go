package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"

	"github.com/josephlbailey/alert-service/internal/api/models"
	"github.com/josephlbailey/alert-service/internal/db"
	"github.com/josephlbailey/alert-service/internal/db/domain"
)

func (s *Server) CreateAlert(c *gin.Context) {
	var (
		req models.CreateAlertReq
		p   domain.CreateAlertParams
	)

	err := req.Bind(c, &p)

	if err != nil {
		return
	}

	s.logger.Info("creating alert...")
	alert, err := s.store.CreateAlertTX(c, p)
	if err != nil {
		s.logger.Error("error creating alert entity", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	s.logger.Info("created alert.", zap.String("externalId", alert.ExternalID.String()))
	c.JSON(http.StatusCreated, models.NewAlertResponse(alert))
}

func (s *Server) GetAlertByExternalID(c *gin.Context) {
	var externalID uuid.UUID
	err := externalID.Parse(c.Param("externalID"))
	if err != nil {
		s.logger.Warn("invalid identifier format, returning 400")
		c.JSON(http.StatusBadRequest, NewError(errors.New("invalid identifier format")))
		return
	}
	s.logger.Info("getting alert...", zap.String("externalId", externalID.String()))
	alert, err := s.store.GetAlertByExternalID(c, externalID)
	if err != nil {

		if errors.Is(err, db.ErrAlertNotExists) {
			s.logger.Warn("alert not found, returning 404")
			c.JSON(http.StatusNotFound, NewError(errors.New("alert not found")))
			return
		}

		s.logger.Error("error getting alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewError(errors.New("error occurred while getting alert")))
		return
	}

	s.logger.Info("returning alert.", zap.String("externalId", alert.ExternalID.String()))
	c.JSON(http.StatusOK, models.NewAlertResponse(alert))
}

func (s *Server) UpdateAlertByExternalID(c *gin.Context) {
	var (
		externalID uuid.UUID
		req        models.UpdateAlertReq
		p          domain.UpdateAlertByIDParams
	)

	err := externalID.Parse(c.Param("externalID"))
	if err != nil {
		s.logger.Warn("invalid identifier format, returning 400")
		c.JSON(http.StatusBadRequest, NewError(errors.New("invalid identifier format")))
		return
	}

	err = req.Bind(c, &p)

	if err != nil {
		return
	}

	alert, err := s.store.GetAlertByExternalID(c, externalID)
	if err != nil {

		if errors.Is(err, db.ErrAlertNotExists) {
			s.logger.Warn("alert not found, returning 404")
			c.JSON(http.StatusNotFound, NewError(errors.New("alert not found")))
			return
		}

		s.logger.Error("error getting alert entity to update", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	p.ID = alert.ID

	s.logger.Info("updating alert...", zap.String("externalID", externalID.String()))
	alert, err = s.store.UpdateAlertByIDTX(c, p)

	if err != nil {
		s.logger.Error("error updating alert entity", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	s.logger.Info("updated alert.", zap.String("externalId", alert.ExternalID.String()))
	c.JSON(http.StatusOK, models.NewAlertResponse(alert))

}

func (s *Server) DeleteAlertByExternalID(c *gin.Context) {
	var externalID uuid.UUID

	err := externalID.Parse(c.Param("externalID"))
	if err != nil {
		s.logger.Warn("invalid identifier format, returning 400")
		c.JSON(http.StatusBadRequest, NewError(errors.New("invalid identifier format")))
		return
	}

	alert, err := s.store.GetAlertByExternalID(c, externalID)
	if err != nil {

		if errors.Is(err, db.ErrAlertNotExists) {
			s.logger.Warn("alert not found, returning 404")
			c.JSON(http.StatusNotFound, NewError(errors.New("alert not found")))
			return
		}

		s.logger.Error("error getting alert entity to delete", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	s.logger.Info("deleting alert...", zap.String("externalID", externalID.String()))
	err = s.store.DeleteAlertByIDTX(c, alert.ID)

	if err != nil {
		s.logger.Error("error deleting alert entity", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	s.logger.Info("deleted alert.", zap.String("externalId", alert.ExternalID.String()))
	c.JSON(http.StatusOK, models.NewAlertResponse(alert))

}
