package public

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pkg/errors"

	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/100bench/cryptocurrency_provider.git/pkg/dto"
)

type Server struct {
	router  *chi.Mux
	service PublicService
}

func NewServer(service PublicService) (*Server, error) {
	if service == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "public server service")
	}

	s := &Server{
		service: service,
	}
	s.router = chi.NewRouter()
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.setupRoutes()
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	s.router.Get("/rates", s.handleGetRates)
	s.router.Get("/currencies", s.handleGetAvailableCurrencies)
	s.router.Get("/rates/{currency}/latest", s.handleGetLatestRate)
	s.router.Get("/rates/{currency}/max", s.handleGetMaxRate)
	s.router.Get("/rates/{currency}/min", s.handleGetMinRate)
	s.router.Get("/rates/{currency}/avg", s.handleGetAvgRate)
}

func (s *Server) handleGetLatestRate(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	if currency == "" {
		s.respondWithError(w, http.StatusBadRequest, "Currency parameter is required")
		return
	}

	rate, err := s.service.GetLatestRate(r.Context(), currency)
	if err != nil {
		log.Printf("Error getting latest rate for %s: %v", currency, err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve latest rate")
		return
	}
	// GetLatestRate возвращает *dto.RateItem, поэтому нужно проверить на nil
	if rate == nil {
		s.respondWithError(w, http.StatusNotFound, fmt.Sprintf("Latest rate for %s not found", currency))
		return
	}
	s.respondWithJSON(w, http.StatusOK, rate)
}

func (s *Server) handleGetMaxRate(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	if currency == "" {
		s.respondWithError(w, http.StatusBadRequest, "Currency parameter is required")
		return
	}

	rate, err := s.service.GetMaxRate(r.Context(), currency)
	if err != nil {
		log.Printf("Error getting max rate for %s: %v", currency, err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve max rate")
		return
	}
	if rate == nil {
		s.respondWithError(w, http.StatusNotFound, fmt.Sprintf("Max rate for %s not found", currency))
		return
	}
	s.respondWithJSON(w, http.StatusOK, rate)
}

func (s *Server) handleGetMinRate(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	if currency == "" {
		s.respondWithError(w, http.StatusBadRequest, "Currency parameter is required")
		return
	}

	rate, err := s.service.GetMinRate(r.Context(), currency)
	if err != nil {
		log.Printf("Error getting min rate for %s: %v", currency, err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve min rate")
		return
	}
	if rate == nil {
		s.respondWithError(w, http.StatusNotFound, fmt.Sprintf("Min rate for %s not found", currency))
		return
	}
	s.respondWithJSON(w, http.StatusOK, rate)
}

func (s *Server) handleGetAvgRate(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	if currency == "" {
		s.respondWithError(w, http.StatusBadRequest, "Currency parameter is required")
		return
	}

	rate, err := s.service.GetAvgRate(r.Context(), currency)
	if err != nil {
		log.Printf("Error getting avg rate for %s: %v", currency, err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve avg rate")
		return
	}
	if rate == nil {
		s.respondWithError(w, http.StatusNotFound, fmt.Sprintf("Avg rate for %s not found", currency))
		return
	}
	s.respondWithJSON(w, http.StatusOK, rate)
}

func (s *Server) handleGetRates(w http.ResponseWriter, r *http.Request) {
	var req dto.GetRatesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request payload: %v", err)
		s.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Currencies) == 0 {
		s.respondWithError(w, http.StatusBadRequest, "At least one currency is required")
		return
	}

	rates, err := s.service.GetRates(r.Context(), req)
	if err != nil {
		log.Printf("Error getting rates: %v", err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve rates")
		return
	}

	s.respondWithJSON(w, http.StatusOK, dto.GetRatesResponse{Rates: rates})
}

func (s *Server) handleGetAvailableCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := s.service.GetAvailableCurrencies(r.Context())
	if err != nil {
		log.Printf("Error getting available currencies: %v", err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve available currencies")
		return
	}

	s.respondWithJSON(w, http.StatusOK, currencies)
}

func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, dto.ErrorResponse{Error: message})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
