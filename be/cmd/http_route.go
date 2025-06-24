package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emicklei/go-restful/v3"

	"github.com/rentiansheng/dashboard/app/metrics/schema"
	"github.com/rentiansheng/dashboard/app/metrics/service"
	"github.com/rentiansheng/dashboard/pkg/config"
)

func initServer() error {

	// Configure logging
	if err := config.Config.ConfigureLogging(); err != nil {
		return err
	}

	container := initRoutes()

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + config.Config.Server.Port,
		Handler:      container,
		ReadTimeout:  config.Config.Server.ReadTimeout,
		WriteTimeout: config.Config.Server.WriteTimeout,
	}

	stop := make(chan os.Signal, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("failed to start server: %w", err)
		}
		stop <- syscall.SIGABRT
		return
	}()

	// 捕获退出信号
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // 阻塞直到退出

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}

// Router interface defines the contract for route registration
type Router interface {
	Routes() *restful.WebService
}

func initRoutes() *restful.Container {
	// Initialize services
	dataSourceSvc := service.NewDataSourceService()
	dataGroupSvc := service.NewGroupKeyService()

	// Initialize schemas
	dataSourceSchema := schema.NewDataSourceSchema(dataSourceSvc, dataGroupSvc)
	groupKeySchema := schema.NewGroupKeySchema(dataGroupSvc)

	// Create container
	container := restful.NewContainer()

	// Add CORS support
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedDomains: []string{"*"},
		CookiesAllowed: true,
		Container:      container,
	}
	container.Filter(cors.Filter)

	// Add logging filter
	container.Filter(logFilter)

	// Enable debugging
	restful.EnableTracing(true)
	restful.DefaultContainer.EnableContentEncoding(true)

	// Register routes
	registerRoutes(container, []Router{
		dataSourceSchema,
		groupKeySchema,
	})

	// Print registered routes for debugging
	for _, ws := range container.RegisteredWebServices() {
		for _, route := range ws.Routes() {
			log.Printf("Registered route: %s %s", route.Method, route.Path)
		}
	}

	return container
}

func registerRoutes(container *restful.Container, routers []Router) {
	for _, router := range routers {
		ws := router.Routes()
		container.Add(ws)
	}
}

func logFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()

	// Log request details
	log.Printf("[REQUEST] %s %s %s", req.Request.Method, req.Request.URL.Path, req.Request.RemoteAddr)

	chain.ProcessFilter(req, resp)

	duration := time.Since(start)
	log.Printf("[RESPONSE] %s %s %s %v", req.Request.Method, req.Request.URL.Path, req.Request.RemoteAddr, duration)
}
