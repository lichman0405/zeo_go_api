package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlockingSpheresHandler struct {
	*BaseHandler
}

func NewBlockingSpheresHandler(base *BaseHandler) *BlockingSpheresHandler {
	return &BlockingSpheresHandler{BaseHandler: base}
}

func (h *BlockingSpheresHandler) Handle(c *gin.Context) {
	var params = make(map[string]interface{})

	// Parse form parameters
	if ha := c.PostForm("ha"); ha == "true" {
		params["ha"] = true
	}

	if probeRadius := c.PostForm("probe_radius"); probeRadius != "" {
		if val, err := strconv.ParseFloat(probeRadius, 64); err == nil {
			params["probe_radius"] = val
		}
	} else {
		params["probe_radius"] = 1.86
	}

	h.ProcessAnalysis(c, "blocking_spheres", params)
}
