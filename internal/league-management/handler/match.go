package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/azmanabdlh/ayo-example/internal/league-management/dtos"
	"github.com/azmanabdlh/ayo-example/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SubstitutePlayer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var param dtos.SubstitutePlayerParam

	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Info("error substitue player param json: %v", err)

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	err := h.svc.SubstitutePlayer(c.Request.Context(), id, param)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(201, gin.H{
		"message": "success to substitute player",
	})

}

func (h *Handler) AssignMatchPlayerLineup(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var param dtos.MatchPlayerParam

	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Info("error match player lineup param json: %v", err)

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	err := h.svc.AssignMatchPlayerLineup(c.Request.Context(), id, param)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "success to assign match player lineup",
	})
}

func (h *Handler) Finish(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.Finish(c.Request.Context(), id)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success set match to finish",
	})
}

func (h *Handler) AddRecordGoal(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var param dtos.GoalParam
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Info("error goal param json: %v", err)

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	err := h.svc.AddRecordGoal(c.Request.Context(), id, param)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(201, gin.H{
		"message": fmt.Sprintf("success add goal info to match_id: %d", id),
	})

}

func (h *Handler) CreateMatch(c *gin.Context) {
	var param dtos.MatchParam
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Info("error match param json: %v", err)

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	err := h.svc.CreateMatch(c.Request.Context(), param)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "success to add match",
	})
}

func (h *Handler) ModifyMatch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id < 1 {
		c.JSON(400, gin.H{
			"message": "invalid match_id",
		})

		return
	}

	var param dtos.MatchParam
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Info("error match param json: %v", err)

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = h.svc.ModifyMatch(c.Request.Context(), id, param)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "success to add match",
	})
}

func (h *Handler) FindMatchHighlight(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id < 1 {
		c.JSON(400, gin.H{
			"message": "invalid match_id",
		})

		return
	}

	matchHighlight, err := h.svc.FindMatchHighlight(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"data": matchHighlight,
	})
}
