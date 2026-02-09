package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vtapaskar/brahma/internal/config"
	"github.com/vtapaskar/brahma/internal/metrics"
	"go.uber.org/zap"
)

type Server struct {
	config    config.ServerConfig
	collector *metrics.Collector
	logger    *zap.Logger
	server    *http.Server
}

func New(cfg config.ServerConfig, collector *metrics.Collector, logger *zap.Logger) *Server {
	s := &Server{
		config:    cfg,
		collector: collector,
		logger:    logger,
	}

	router := mux.NewRouter()
	s.setupRoutes(router)

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) setupRoutes(router *mux.Router) {
	router.HandleFunc("/health", s.healthHandler).Methods("GET")

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/metrics", s.metricsHandler).Methods("POST")
	api.HandleFunc("/crash-report", s.crashReportHandler).Methods("POST")
	api.HandleFunc("/backtrace", s.backtraceHandler).Methods("POST")
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.collector.Stop()
	return s.server.Shutdown(ctx)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "brahma",
	})
}

func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.DeviceMetric

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		s.logger.Warn("Failed to decode metric payload", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.collector.CollectMetric(metric); err != nil {
		s.logger.Error("Failed to collect metric", zap.Error(err))
		http.Error(w, "Failed to process metric", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
	})
}

func (s *Server) crashReportHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		s.logger.Warn("Failed to parse multipart form", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	deviceID := r.FormValue("device_id")
	hostname := r.FormValue("hostname")
	reportType := r.FormValue("report_type")

	if deviceID == "" {
		http.Error(w, "device_id is required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		s.logger.Warn("Failed to get file from form", zap.Error(err))
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("Failed to read file content", zap.Error(err))
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	report := metrics.CrashReport{
		DeviceID:   deviceID,
		Hostname:   hostname,
		ReportType: reportType,
		Content:    content,
		Filename:   header.Filename,
	}

	reportID, err := s.collector.CollectCrashReport(report)
	if err != nil {
		s.logger.Error("Failed to store crash report", zap.Error(err))
		http.Error(w, "Failed to store crash report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "created",
		"report_id": reportID,
	})
}

func (s *Server) backtraceHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		s.logger.Warn("Failed to parse multipart form", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	deviceID := r.FormValue("device_id")
	hostname := r.FormValue("hostname")

	if deviceID == "" {
		http.Error(w, "device_id is required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		s.logger.Warn("Failed to get file from form", zap.Error(err))
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("Failed to read file content", zap.Error(err))
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	report := metrics.CrashReport{
		DeviceID:   deviceID,
		Hostname:   hostname,
		ReportType: "backtrace",
		Content:    content,
		Filename:   header.Filename,
	}

	reportID, err := s.collector.CollectBacktrace(report)
	if err != nil {
		s.logger.Error("Failed to store backtrace", zap.Error(err))
		http.Error(w, "Failed to store backtrace", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "created",
		"report_id": reportID,
	})
}
