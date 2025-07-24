package handlers

import (
	"context"
	"fmt"
	"net/http"

	"zeo-api/internal/config"
	"zeo-api/internal/core/cache"
	"zeo-api/internal/core/parser"
	"zeo-api/internal/core/runner"
	"zeo-api/internal/utils/file"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	zeoRunner *runner.ZeoRunner
	cache     *cache.Cache
	config    *config.Config
}

func NewBaseHandler(zeoRunner *runner.ZeoRunner, cacheInstance *cache.Cache, cfg *config.Config) *BaseHandler {
	return &BaseHandler{
		zeoRunner: zeoRunner,
		cache:     cacheInstance,
		config:    cfg,
	}
}

func (h *BaseHandler) ProcessAnalysis(c *gin.Context, analysisType string, params map[string]interface{}) {
	// Get uploaded file
	fileHeader, err := c.FormFile("structure_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "structure_file is required",
		})
		return
	}

	// Validate file extension
	if !file.IsValidStructureFile(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid file format. Supported: .cif, .cssr, .v1, .arc",
		})
		return
	}

	// Save uploaded file
	savedPath, err := file.SaveUploadedFile(fileHeader, analysisType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("failed to save file: %v", err),
		})
		return
	}
	defer file.CleanupFile(savedPath)

	// Build Zeo++ arguments
	zeoArgs, err := runner.BuildZeoArgs(analysisType, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("invalid parameters: %v", err),
		})
		return
	}
	outputFiles := getOutputFiles(analysisType)

	// Generate cache key
	cacheKey := cache.GenerateCacheKey(savedPath, zeoArgs)

	// Check cache
	if h.config.Cache.Enabled {
		if cachedData, found := h.cache.Get(cacheKey); found {
			if parsed, ok := cachedData[outputFiles[0]]; ok {
				result, err := parser.ParseOutputFile(analysisType, string(parsed))
				if err == nil {
					c.JSON(http.StatusOK, gin.H{
						"success": true,
						"data":    result,
						"cached":  true,
					})
					return
				}
			}
		}
	}

	// Execute Zeo++ analysis
	ctx, cancel := context.WithTimeout(context.Background(), h.config.Zeo.Timeout)
	defer cancel()

	result, err := h.zeoRunner.RunCommand(ctx, savedPath, zeoArgs, outputFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Zeo++ execution failed: %v", err),
		})
		return
	}

	if !result.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Zeo++ error: %s", result.Stderr),
			"stdout":  result.Stdout,
		})
		return
	}

	// Parse and cache results
	mainOutput := outputFiles[0]
	if outputData, exists := result.OutputFiles[mainOutput]; exists {
		parsedResult, err := parser.ParseOutputFile(analysisType, string(outputData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("failed to parse results: %v", err),
			})
			return
		}

		// Cache the results
		if h.config.Cache.Enabled {
			cacheData := map[string][]byte{
				mainOutput: outputData,
			}
			h.cache.Set(cacheKey, cacheData)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    parsedResult,
			"cached":  false,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "no output generated from Zeo++",
		})
	}
}

func getOutputFiles(analysisType string) []string {
	switch analysisType {
	case "pore_diameter":
		return []string{"output.res"}
	case "surface_area":
		return []string{"output.sa"}
	case "accessible_volume":
		return []string{"output.vol"}
	case "probe_volume":
		return []string{"output.volpo"}
	case "channel_analysis":
		return []string{"output.chan"}
	case "framework_info":
		return []string{"output.strinfo"}
	case "blocking_spheres":
		return []string{"output.block"}
	case "open_metal_sites":
		return []string{"output.oms"}
	case "pore_size_dist":
		return []string{"output.psd"}
	default:
		return []string{"output"}
	}
}

func (h *BaseHandler) ProcessFileDownload(c *gin.Context, analysisType string, params map[string]interface{}) {
	// Get uploaded file
	fileHeader, err := c.FormFile("structure_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "structure_file is required",
		})
		return
	}

	// Validate file extension
	if !file.IsValidStructureFile(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid file format. Supported: .cif, .cssr, .v1, .arc",
		})
		return
	}

	// Save uploaded file
	savedPath, err := file.SaveUploadedFile(fileHeader, analysisType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("failed to save file: %v", err),
		})
		return
	}
	defer file.CleanupFile(savedPath)

	// Build Zeo++ arguments
	zeoArgs, err := runner.BuildZeoArgs(analysisType, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("invalid parameters: %v", err),
		})
		return
	}
	outputFiles := getOutputFiles(analysisType)

	// Execute Zeo++ analysis
	ctx, cancel := context.WithTimeout(context.Background(), h.config.Zeo.Timeout)
	defer cancel()

	result, err := h.zeoRunner.RunCommand(ctx, savedPath, zeoArgs, outputFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Zeo++ execution failed: %v", err),
		})
		return
	}

	if !result.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Zeo++ error: %s", result.Stderr),
			"stdout":  result.Stdout,
		})
		return
	}

	// Serve the file
	mainOutput := outputFiles[0]
	if outputData, exists := result.OutputFiles[mainOutput]; exists {
		filename := fmt.Sprintf("%s_%s", analysisType, mainOutput)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		c.Data(http.StatusOK, "application/octet-stream", outputData)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "no output generated from Zeo++",
		})
	}
}
