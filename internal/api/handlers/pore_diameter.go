package handlers

import (
	"github.com/gin-gonic/gin"
)

type PoreDiameterHandler struct {
	*BaseHandler
}

func NewPoreDiameterHandler(base *BaseHandler) *PoreDiameterHandler {
	return &PoreDiameterHandler{BaseHandler: base}
}

func (h *PoreDiameterHandler) Handle(c *gin.Context) {
	var params = make(map[string]interface{})

	// Parse form parameters
	if ha := c.PostForm("ha"); ha == "true" {
		params["ha"] = true
	}

	h.ProcessAnalysis(c, "pore_diameter", params)
}
