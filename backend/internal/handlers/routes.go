package handlers

import (
	"brickwise/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(
	r chi.Router,
	propertyService *service.PropertyService,
	timeEntryService *service.TimeEntryService,
	residenceService *service.ResidenceService,
	mortgageService *service.MortgageService,
	dashboardService *service.DashboardService,
	qboService *service.QBOService,
) {
	propertyHandler := NewPropertyHandler(propertyService)
	timeEntryHandler := NewTimeEntryHandler(timeEntryService)
	residenceHandler := NewResidenceHandler(residenceService)
	mortgageHandler := NewMortgageHandler(mortgageService)
	dashboardHandler := NewDashboardHandler(dashboardService)
	qboHandler := NewQBOHandler(qboService)

	r.Get("/dashboard/summary", dashboardHandler.GetSummary)

	r.Route("/integrations/qbo", func(r chi.Router) {
		r.Get("/status", qboHandler.GetStatus)
		r.Get("/connect", qboHandler.Connect)
		r.Get("/callback", qboHandler.Callback)
		r.Delete("/disconnect", qboHandler.Disconnect)
		r.Get("/classes", qboHandler.ListClasses)
		r.Put("/mappings", qboHandler.UpdateMappings)
		r.Post("/sync", qboHandler.Sync)
	})

	r.Route("/properties", func(r chi.Router) {
		r.Get("/", propertyHandler.ListProperties)
		r.Post("/", propertyHandler.CreateProperty)
		r.Get("/{id}", propertyHandler.GetProperty)
		r.Put("/{id}", propertyHandler.UpdateProperty)
		r.Put("/{id}/value", propertyHandler.UpdatePropertyValue)
		r.Get("/{id}/mortgage", mortgageHandler.GetMortgage)
		r.Put("/{id}/mortgage", mortgageHandler.UpsertMortgage)
		r.Delete("/{id}/mortgage", mortgageHandler.DeleteMortgage)
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

	r.Route("/residences", func(r chi.Router) {
		r.Get("/", residenceHandler.ListResidences)
		r.Post("/", residenceHandler.CreateResidence)
		r.Get("/{id}", residenceHandler.GetResidence)
		r.Put("/{id}", residenceHandler.UpdateResidence)
	})
}
