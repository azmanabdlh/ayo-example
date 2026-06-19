package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/azmanabdlh/ayo-example/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AssignPlayerToTeam(c *gin.Context) {
	req := struct {
		PlayerID   int64 `json:"player_id"`
		TeamID     int64 `json:"team_id"`
		BackNumber int   `json:"back_number"`
	}{}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Info("error parse param assign player to team json: %v", err)

		c.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})

		return
	}

	err := h.svc.AssignPlayerTeam(c.Request.Context(), req.PlayerID, req.TeamID, req.BackNumber)
	if err != nil {
		logger.Info("error assign player to team json: %v", req)

		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(201, gin.H{
		"code":    201,
		"message": fmt.Sprintf("successfully assign player:%d  to team id: %d", req.PlayerID, req.TeamID),
	})

}

func (h *Handler) FindPlayer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	player, err := h.svc.FindPlayer(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": player,
	})
}

func (h *Handler) RemovePlayer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.RemovePlayer(c.Request.Context(), id)
	if err != nil {
		logger.Info("error remove player id: %d", id)

		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("successfully remove player by id: %d", id),
	})
}
