package public

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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
	s.router.Get("/rates", s.handleGetMin)
	s.router.Get("/currencies", s.handleGetMax)
	s.router.Get("/rates/{currency}/latest", s.handleGetAvg)
	s.router.Get("/rates/{currency}/max", s.handleGetLast)
}

// @Summary Get minimum rates for currencies
// @Description Get minimum rates for a list of currencies
// @Tags rates
// @Accept json
// @Produce json
// @Param currencies query string true "Comma-separated list of currency symbols (e.g., BTC,ETH)"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /rates [get]
func (s *Server) handleGetMin(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	currencies := q.Get("currencies")
	if currencies == "" {
		s.respondWithError(w, http.StatusBadRequest, "missing 'currencies' query parameter")
		return
	}
	currencyList := strings.Split(currencies, ",")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	rates, err := s.service.GetMinRate(ctx, currencyList)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get rates")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rates)
}

// @Summary Get maximum rates for currencies
// @Description Get maximum rates for a list of currencies
// @Tags rates
// @Accept json
// @Produce json
// @Param currencies query string true "Comma-separated list of currency symbols (e.g., BTC,ETH)"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /currencies [get]
func (s *Server) handleGetMax(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	currencies := q.Get("currencies")
	if currencies == "" {
		s.respondWithError(w, http.StatusBadRequest, "missing 'currencies' query parameter")
		return
	}
	currencyList := strings.Split(currencies, ",")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	rates, err := s.service.GetMaxRate(ctx, currencyList)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get rates")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rates)
}

// @Summary Get average rate for a currency
// @Description Get the average rate for a specific currency
// @Tags rates
// @Accept json
// @Produce json
// @Param currency path string true "Currency symbol (e.g., BTC)"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /rates/{currency}/latest [get]
func (s *Server) handleGetAvg(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	if currency == "" {
		s.respondWithError(w, http.StatusBadRequest, "missing 'currency' path parameter")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	rates, err := s.service.GetAvgRate(ctx, []string{currency})
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get rates")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rates)
}

// @Summary Get last rate for a currency
// @Description Get the last recorded rate for a specific currency
// @Tags rates
// @Accept json
// @Produce json
// @Param currency path string true "Currency symbol (e.g., BTC)"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /rates/{currency}/max [get]
func (s *Server) handleGetLast(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	if currency == "" {
		s.respondWithError(w, http.StatusBadRequest, "missing 'currency' path parameter")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	rates, err := s.service.GetLast(ctx, []string{currency})
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get rates")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rates)
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
