package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func DefaultTechOptions() RouterOption {
	return RouterOptions(
		WithRecover(),
		WithDebugHandler(),
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

func WithApiHandler() RouterOption {
	return func(r chi.Router) {

	}
}

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
