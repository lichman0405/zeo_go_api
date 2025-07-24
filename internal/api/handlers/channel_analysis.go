package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChannelAnalysisHandler struct {
	*BaseHandler
}

func NewChannelAnalysisHandler(base *BaseHandler) *ChannelAnalysisHandler {
	return &ChannelAnalysisHandler{BaseHandler: base}
}

func (h *ChannelAnalysisHandler) Handle(c *gin.Context) {
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
		params["probe_radius"] = 1.21
	}

	h.ProcessAnalysis(c, "channel_analysis", params)
}
