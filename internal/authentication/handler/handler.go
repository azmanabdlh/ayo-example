package handler

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

type Service interface {
	Login(ctx context.Context, email, passwordPlain string) (string, error)
}

func New(svc Service) *Handler {
	h := new(Handler)
	h.svc = svc

	return h
}

func (h *Handler) Login(c *gin.Context) {
	req := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})

		return
	}

	token, err := h.svc.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": token,
	})
}
