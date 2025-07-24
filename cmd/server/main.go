package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zeo-api/internal/api/handlers"
	"zeo-api/internal/api/middleware"
	"zeo-api/internal/config"
	"zeo-api/internal/core/cache"
	"zeo-api/internal/core/runner"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Printf("Failed to load config from file: %v, using defaults", err)
		cfg, err = config.LoadDefaultConfig()
		if err != nil {
			log.Fatalf("Failed to load default config: %v", err)
		}
	}

	// Initialize Zeo++ runner
	zeoRunner := runner.NewZeoRunner(&cfg.Zeo)
	if err := zeoRunner.ValidateZeoExecutable(); err != nil {
		log.Fatalf("Zeo++ executable not found: %v", err)
	}

	// Initialize cache
	cacheInstance := cache.NewCache(&cfg.Cache, cfg.Zeo.Workdir)

	// Initialize Gin
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(rate.Limit(cfg.Concurrency.RateLimitPerIP), cfg.Concurrency.RateLimitPerIP*2)
	go rateLimiter.Cleanup()

	// Global semaphore for concurrent requests
	globalLimiter := middleware.NewGlobalSemaphore(cfg.Concurrency.MaxConcurrentUploads)

	// Initialize base handler
	baseHandler := handlers.NewBaseHandler(zeoRunner, cacheInstance, cfg)

	// Initialize specific handlers
	poreDiameterHandler := handlers.NewPoreDiameterHandler(baseHandler)
	surfaceAreaHandler := handlers.NewSurfaceAreaHandler(baseHandler)
	accessibleVolumeHandler := handlers.NewAccessibleVolumeHandler(baseHandler)
	probeVolumeHandler := handlers.NewProbeVolumeHandler(baseHandler)
	channelAnalysisHandler := handlers.NewChannelAnalysisHandler(baseHandler)
	frameworkInfoHandler := handlers.NewFrameworkInfoHandler(baseHandler)
	blockingSpheresHandler := handlers.NewBlockingSpheresHandler(baseHandler)
	openMetalSitesHandler := handlers.NewOpenMetalSitesHandler(baseHandler)
	poreSizeDistHandler := handlers.NewPoreSizeDistHandler(baseHandler)

	// API routes
	api := router.Group("/api")
	{
		// Apply middleware
		api.Use(rateLimiter.RateLimit())
		api.Use(globalLimiter.Middleware())

		// Analysis endpoints
		api.POST("/pore_diameter", poreDiameterHandler.Handle)
		api.POST("/surface_area", surfaceAreaHandler.Handle)
		api.POST("/accessible_volume", accessibleVolumeHandler.Handle)
		api.POST("/probe_volume", probeVolumeHandler.Handle)
		api.POST("/channel_analysis", channelAnalysisHandler.Handle)
		api.POST("/framework_info", frameworkInfoHandler.Handle)
		api.POST("/blocking_spheres", blockingSpheresHandler.Handle)
		api.POST("/open_metal_sites", openMetalSitesHandler.Handle)
		api.POST("/pore_size_dist/download", poreSizeDistHandler.Handle)
	}

	// Health check endpoint
	router.Any("/health", func(c *gin.Context) {
		if c.Request.Method == http.MethodHead {
			c.Status(http.StatusOK)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Zeo++ Analysis API",
			"version": "1.0.0",
			"endpoints": []string{
				"POST /api/pore_diameter",
				"POST /api/surface_area",
				"POST /api/accessible_volume",
				"POST /api/probe_volume",
				"POST /api/channel_analysis",
				"POST /api/framework_info",
				"POST /api/pore_size_dist/download",
				"POST /api/blocking_spheres",
				"POST /api/open_metal_sites",
			},
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
