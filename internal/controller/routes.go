package controller

import (
	pkghttp "github.com/JonnyShabli/EffectiveMobile/pkg/http"
	"github.com/go-chi/chi/v5"
)

func WithApiHandler(api SubsHandlerInterface) pkghttp.RouterOption {
	return func(r chi.Router) {
		r.Route("/api/subs", func(r chi.Router) {
			r.Post("/", api.InsertSub)
			r.Get("/", api.GetSub)
			r.Put("/", api.UpdateSub)
			r.Delete("/{sub_id}", api.DeleteSub)
			r.Get("/list", api.ListSub)
			r.Get("/sumPriceByDate", api.SumPriceByDate)
		})
	}
}
