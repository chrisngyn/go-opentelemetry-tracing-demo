package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"github.com/rs/cors"

	"github.com/chiennguyen196/go-opentelemetry-tracing-demo/pkg/tracing"
)

func main() {
	cleanFn := tracing.Init()
	defer func() {
		if err := cleanFn(); err != nil {
			log.Println("Clean tracing error", err)
		}
	}()

	httpServer := NewHttpServer()

	address := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	RunHTTPServerOnAddr(address, func(router chi.Router) http.Handler {
		return HandlerFromMux(httpServer, router)
	})
}

func RunHTTPServerOnAddr(addr string, createHandler func(router chi.Router) http.Handler) {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	// we are mounting all APIs under /api path
	rootRouter.Mount("/api", createHandler(apiRouter))

	log.Println("Starting HTTP server", addr)

	err := http.ListenAndServe(addr, rootRouter)
	if err != nil {
		log.Panicln("Unable to start HTTP server", err)
	}
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	router.Use(otelchi.Middleware("my-server", otelchi.WithChiRoutes(router)))
	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}
