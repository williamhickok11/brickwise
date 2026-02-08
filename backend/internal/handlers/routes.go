package handlers

import (
	"brickwise/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(r chi.Router, propertyService *service.PropertyService, timeEntryService *service.TimeEntryService) {
	propertyHandler := NewPropertyHandler(propertyService)
	timeEntryHandler := NewTimeEntryHandler(timeEntryService)

	r.Route("/properties", func(r chi.Router) {
		r.Get("/", propertyHandler.ListProperties)
		r.Post("/", propertyHandler.CreateProperty)
		r.Get("/{id}", propertyHandler.GetProperty)
		r.Put("/{id}", propertyHandler.UpdateProperty)
		r.Delete("/{id}", propertyHandler.DeleteProperty)
	})

	r.Route("/time-entries", func(r chi.Router) {
		r.Get("/", timeEntryHandler.ListTimeEntries)
		r.Post("/", timeEntryHandler.CreateTimeEntry)
		r.Get("/export", timeEntryHandler.ExportTimeEntries)
		r.Get("/{id}", timeEntryHandler.GetTimeEntry)
		r.Put("/{id}", timeEntryHandler.UpdateTimeEntry)
		r.Delete("/{id}", timeEntryHandler.DeleteTimeEntry)
	})
}
