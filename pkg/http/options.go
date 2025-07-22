package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func DefaultTechOptions() RouterOption {
	return RouterOptions(
		WithRecover(),
		WithDebugHandler(),
		WithSwagger(),
	)
}

func RouterOptions(options ...RouterOption) func(chi.Router) {
	return func(r chi.Router) {
		for _, option := range options {
			option(r)
		}
	}
}

type RouterOption func(chi.Router)

func WithDebugHandler() RouterOption {
	return func(r chi.Router) {
		r.Mount("/debug", middleware.Profiler())
	}
}

// WithRecover adds recover middleware, which can catch panics from handlers.
func WithRecover() RouterOption {
	return func(r chi.Router) {
		r.Use(middleware.Recoverer)
	}
}

func WithSwagger() RouterOption {
	return func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		))
	}
}
