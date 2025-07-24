package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type SurfaceAreaHandler struct {
	*BaseHandler
}

func NewSurfaceAreaHandler(base *BaseHandler) *SurfaceAreaHandler {
	return &SurfaceAreaHandler{BaseHandler: base}
}

func (h *SurfaceAreaHandler) Handle(c *gin.Context) {
	var params = make(map[string]interface{})

	// Parse form parameters
	if ha := c.PostForm("ha"); ha == "true" {
		params["ha"] = true
	}

	if probeRadius := c.PostForm("probe_radius"); probeRadius != "" {
		if val, err := strconv.ParseFloat(probeRadius, 64); err == nil {
			params["probe_radius"] = val
		}
	}

	if samples := c.PostForm("samples"); samples != "" {
		if val, err := strconv.Atoi(samples); err == nil {
			params["samples"] = val
		}
	} else {
		params["samples"] = 2000
	}

	h.ProcessAnalysis(c, "surface_area", params)
}
