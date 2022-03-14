package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jasanchez1/Dpricing/db"
	"github.com/jasanchez1/Dpricing/models"
)

var vendorIDKey = "vendorID"

func vendors(router chi.Router) {
	router.Get("/", getAllVendors)
	router.Post("/", createVendor)
	router.Route("/{vendorId}", func(router chi.Router) {
		router.Use(VendorContext)
		router.Get("/", getVendor)
		router.Put("/", updateVendor)
		router.Delete("/", deleteVendor)
	})
}
func VendorContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vendorId := chi.URLParam(r, "vendorId")
		if vendorId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("vendor ID is required")))
			return
		}
		id, err := strconv.Atoi(vendorId)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid vendor ID")))
		}
		ctx := context.WithValue(r.Context(), vendorIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createVendor(w http.ResponseWriter, r *http.Request) {
	vendor := &models.Vendor{}
	if err := render.Bind(r, vendor); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddVendor(vendor); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, vendor); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getAllVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := dbInstance.GetAllVendors()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, vendors); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getVendor(w http.ResponseWriter, r *http.Request) {
	vendorID := r.Context().Value(vendorIDKey).(int)
	vendor, err := dbInstance.GetVendorById(vendorID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &vendor); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteVendor(w http.ResponseWriter, r *http.Request) {
	vendorId := r.Context().Value(vendorIDKey).(int)
	err := dbInstance.DeleteVendor(vendorId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateVendor(w http.ResponseWriter, r *http.Request) {
	vendorId := r.Context().Value(vendorIDKey).(int)
	vendorData := models.Vendor{}
	if err := render.Bind(r, &vendorData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	vendor, err := dbInstance.UpdateVendor(vendorId, vendorData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &vendor); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
