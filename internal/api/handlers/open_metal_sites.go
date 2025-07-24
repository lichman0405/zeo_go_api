package handlers

import (
	"github.com/gin-gonic/gin"
)

type OpenMetalSitesHandler struct {
	*BaseHandler
}

func NewOpenMetalSitesHandler(base *BaseHandler) *OpenMetalSitesHandler {
	return &OpenMetalSitesHandler{BaseHandler: base}
}

func (h *OpenMetalSitesHandler) Handle(c *gin.Context) {
	var params = make(map[string]interface{})

	// Parse form parameters
	if ha := c.PostForm("ha"); ha == "true" {
		params["ha"] = true
	}

	h.ProcessAnalysis(c, "open_metal_sites", params)
}
