package handlers

import (
	"github.com/gin-gonic/gin"
)

type FrameworkInfoHandler struct {
	*BaseHandler
}

func NewFrameworkInfoHandler(base *BaseHandler) *FrameworkInfoHandler {
	return &FrameworkInfoHandler{BaseHandler: base}
}

func (h *FrameworkInfoHandler) Handle(c *gin.Context) {
	var params = make(map[string]interface{})

	// Parse form parameters
	if ha := c.PostForm("ha"); ha == "true" {
		params["ha"] = true
	}

	h.ProcessAnalysis(c, "framework_info", params)
}
