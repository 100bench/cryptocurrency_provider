package public

import (
	"encoding/json"
	"net/http"

	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/100bench/cryptocurrency_provider.git/pkg/dto"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
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
	s.setupRoutes()
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	s.router.Get("/rates", s.handleGetRates)
	s.router.Get("/currencies", s.handleGetAvailableCurrencies)
	// s.router.Get("/rates/{currency}/latest", s.handleGetLatestRate)
}

func (s *Server) handleGetRates(w http.ResponseWriter, r *http.Request) {
	var req dto.GetRatesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Currencies) == 0 {
		s.respondWithError(w, http.StatusBadRequest, "At least one currency is required")
		return
	}

	rates, err := s.service.GetRates(r.Context(), req)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve rates")
		return
	}

	s.respondWithJSON(w, http.StatusOK, dto.GetRatesResponse{Rates: rates})
}

func (s *Server) handleGetAvailableCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := s.service.GetAvailableCurrencies(r.Context())
	if err != nil {
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
