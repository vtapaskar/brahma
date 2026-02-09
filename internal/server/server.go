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
	"github.com/vtapaskar/brahma/internal/registry"
	"go.uber.org/zap"
)

type Server struct {
	config    config.ServerConfig
	collector *metrics.Collector
	registry  *registry.Registry
	logger    *zap.Logger
	server    *http.Server
}

func New(cfg config.ServerConfig, collector *metrics.Collector, reg *registry.Registry, logger *zap.Logger) *Server {
	s := &Server{
		config:    cfg,
		collector: collector,
		registry:  reg,
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
	api.HandleFunc("/devices/register", s.registerDeviceHandler).Methods("POST")
	api.HandleFunc("/devices", s.listDevicesHandler).Methods("GET")
	api.HandleFunc("/devices/{uid}", s.getDeviceHandler).Methods("GET")
	api.HandleFunc("/devices/{uid}", s.unregisterDeviceHandler).Methods("DELETE")

	api.HandleFunc("/metrics/{uid}/cpu", s.cpuStatsHandler).Methods("POST")
	api.HandleFunc("/metrics/{uid}/process", s.processStatsHandler).Methods("POST")
	api.HandleFunc("/metrics/{uid}/mgmt-network", s.mgmtNetworkStatsHandler).Methods("POST")
	api.HandleFunc("/metrics/{uid}/router-state", s.routerBaseStateHandler).Methods("POST")

	api.HandleFunc("/log/{uid}/crash", s.crashReportHandler).Methods("POST")
	api.HandleFunc("/log/{uid}/backtrace", s.backtraceHandler).Methods("POST")
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

func (s *Server) validateUID(uid string) bool {
	_, exists := s.registry.GetByUID(uid)
	return exists
}

func (s *Server) cpuStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if !s.validateUID(uid) {
		http.Error(w, "Device not registered", http.StatusNotFound)
		return
	}

	var stats metrics.CPUStats
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		s.logger.Warn("Failed to decode CPU stats", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stats.UID = uid
	if err := s.collector.CollectCPUStats(&stats); err != nil {
		s.logger.Error("Failed to collect CPU stats", zap.Error(err))
		http.Error(w, "Failed to process CPU stats", http.StatusInternalServerError)
		return
	}

	s.registry.UpdateLastSeen(uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"uid":    uid,
	})
}

func (s *Server) processStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if !s.validateUID(uid) {
		http.Error(w, "Device not registered", http.StatusNotFound)
		return
	}

	var stats metrics.ProcessStats
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		s.logger.Warn("Failed to decode process stats", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stats.UID = uid
	if err := s.collector.CollectProcessStats(&stats); err != nil {
		s.logger.Error("Failed to collect process stats", zap.Error(err))
		http.Error(w, "Failed to process stats", http.StatusInternalServerError)
		return
	}

	s.registry.UpdateLastSeen(uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"uid":    uid,
	})
}

func (s *Server) mgmtNetworkStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if !s.validateUID(uid) {
		http.Error(w, "Device not registered", http.StatusNotFound)
		return
	}

	var stats metrics.MgmtNetworkStats
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		s.logger.Warn("Failed to decode mgmt network stats", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stats.UID = uid
	if err := s.collector.CollectMgmtNetworkStats(&stats); err != nil {
		s.logger.Error("Failed to collect mgmt network stats", zap.Error(err))
		http.Error(w, "Failed to process mgmt network stats", http.StatusInternalServerError)
		return
	}

	s.registry.UpdateLastSeen(uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"uid":    uid,
	})
}

func (s *Server) routerBaseStateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if !s.validateUID(uid) {
		http.Error(w, "Device not registered", http.StatusNotFound)
		return
	}

	var state metrics.RouterBaseState
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		s.logger.Warn("Failed to decode router base state", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	state.UID = uid
	if err := s.collector.CollectRouterBaseState(&state); err != nil {
		s.logger.Error("Failed to collect router base state", zap.Error(err))
		http.Error(w, "Failed to process router base state", http.StatusInternalServerError)
		return
	}

	s.registry.UpdateLastSeen(uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"uid":    uid,
	})
}

func (s *Server) crashReportHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if !s.validateUID(uid) {
		http.Error(w, "Device not registered", http.StatusNotFound)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		s.logger.Warn("Failed to parse multipart form", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	processTag := r.FormValue("process_tag")
	version := r.FormValue("version")

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

	report := &metrics.LogReport{
		DeviceUID:  uid,
		ProcessTag: processTag,
		Version:    version,
		Content:    content,
		Filename:   header.Filename,
	}

	logID, err := s.collector.CollectCrashReport(report)
	if err != nil {
		s.logger.Error("Failed to store crash report", zap.Error(err))
		http.Error(w, "Failed to store crash report", http.StatusInternalServerError)
		return
	}

	s.registry.UpdateLastSeen(uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
		"log_id": logID,
		"uid":    uid,
		"s3_key": report.S3Key,
	})
}

func (s *Server) backtraceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if !s.validateUID(uid) {
		http.Error(w, "Device not registered", http.StatusNotFound)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		s.logger.Warn("Failed to parse multipart form", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	processTag := r.FormValue("process_tag")
	version := r.FormValue("version")

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

	report := &metrics.LogReport{
		DeviceUID:  uid,
		ProcessTag: processTag,
		Version:    version,
		Content:    content,
		Filename:   header.Filename,
	}

	logID, err := s.collector.CollectBacktrace(report)
	if err != nil {
		s.logger.Error("Failed to store backtrace", zap.Error(err))
		http.Error(w, "Failed to store backtrace", http.StatusInternalServerError)
		return
	}

	s.registry.UpdateLastSeen(uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
		"log_id": logID,
		"uid":    uid,
		"s3_key": report.S3Key,
	})
}

func (s *Server) registerDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var req registry.RegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warn("Failed to decode registration request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ForeignKey == "" {
		http.Error(w, "foreign_key is required", http.StatusBadRequest)
		return
	}

	device, err := s.registry.Register(req)
	if err != nil {
		s.logger.Error("Failed to register device", zap.Error(err))
		http.Error(w, "Failed to register device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

func (s *Server) listDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices := s.registry.ListDevices()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (s *Server) getDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	device, exists := s.registry.GetByUID(uid)
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (s *Server) unregisterDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if !s.registry.Unregister(uid) {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "unregistered",
		"uid":    uid,
	})
}
