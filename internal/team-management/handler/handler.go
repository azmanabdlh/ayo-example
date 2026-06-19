package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/azmanabdlh/ayo-example/internal/team-management/dtos"
	"github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"github.com/azmanabdlh/ayo-example/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll(ctx context.Context, page, limit int) ([]model.Team, error)
	Find(ctx context.Context, id int64) (model.Team, error)
	Modify(ctx context.Context, teamID int64, param dtos.TeamParam) error
	Create(ctx context.Context, param dtos.TeamParam) error
	Remove(ctx context.Context, teamID int64) error

	// player
	AssignPlayerTeam(ctx context.Context, playerID, teamID int64, newBackNumber int) error
	FindPlayer(ctx context.Context, playerID int64) (model.Player, error)
	RemovePlayer(ctx context.Context, playerID int64) error
}

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	h := new(Handler)
	h.svc = svc

	return h
}

func (h *Handler) FindAllTeam(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	teams, err := h.svc.FindAll(c.Request.Context(), offset, limit)
	if err != nil {
		logger.Info("error findAll team.")

		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": teams,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *Handler) FindTeam(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	team, err := h.svc.Find(c.Request.Context(), id)
	if err != nil {
		logger.Info("error find team id: %d", id)

		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": team,
	})
}

func (h *Handler) ModifyTeam(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var team dtos.TeamParam

	if err := c.ShouldBindJSON(&team); err != nil {
		logger.Info("error team param json: %v", err)

		c.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})

		return
	}

	err := h.svc.Modify(c.Request.Context(), id, team)
	if err != nil {
		logger.Info("error modify team id: %d", id)

		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("successfully modify team by id: %d", id),
	})
}

func (h *Handler) CreateTeam(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var team dtos.TeamParam

	if err := c.ShouldBindJSON(&team); err != nil {
		logger.Info("error team param json: %v", err)

		c.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})

		return
	}

	err := h.svc.Create(c.Request.Context(), team)
	if err != nil {
		logger.Info("error find team id: %d", id)

		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(201, gin.H{
		"code":    201,
		"message": "successfully add team",
	})
}

func (h *Handler) RemoveTeam(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.Remove(c.Request.Context(), id)
	if err != nil {
		logger.Info("error find team id: %d", id)

		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("successfully remove team by id: %d", id),
	})
}
