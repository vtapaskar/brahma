package registry

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vtapaskar/brahma/internal/splunk"
	"go.uber.org/zap"
)

type DeviceRegistration struct {
	UID          string            `json:"uid"`
	ForeignKey   string            `json:"foreign_key"`
	Hostname     string            `json:"hostname"`
	IPAddress    string            `json:"ip_address"`
	DeviceType   string            `json:"device_type"`
	Platform     string            `json:"platform"`
	Version      string            `json:"version"`
	Labels       map[string]string `json:"labels,omitempty"`
	RegisteredAt time.Time         `json:"registered_at"`
	LastSeen     time.Time         `json:"last_seen"`
}

type RegistrationRequest struct {
	ForeignKey string            `json:"foreign_key"`
	Hostname   string            `json:"hostname"`
	IPAddress  string            `json:"ip_address"`
	DeviceType string            `json:"device_type"`
	Platform   string            `json:"platform"`
	Version    string            `json:"version"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type Registry struct {
	devices      map[string]*DeviceRegistration
	byForeignKey map[string]string
	mu           sync.RWMutex
	splunkClient *splunk.Client
	logger       *zap.Logger
}

func NewRegistry(splunkClient *splunk.Client, logger *zap.Logger) *Registry {
	return &Registry{
		devices:      make(map[string]*DeviceRegistration),
		byForeignKey: make(map[string]string),
		splunkClient: splunkClient,
		logger:       logger,
	}
}

func (r *Registry) Register(req RegistrationRequest) (*DeviceRegistration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existingUID, exists := r.byForeignKey[req.ForeignKey]; exists {
		device := r.devices[existingUID]
		device.Hostname = req.Hostname
		device.IPAddress = req.IPAddress
		device.DeviceType = req.DeviceType
		device.Platform = req.Platform
		device.Version = req.Version
		device.Labels = req.Labels
		device.LastSeen = time.Now()

		r.sendRegistrationEvent(device, "device_updated")

		r.logger.Info("Device updated",
			zap.String("uid", device.UID),
			zap.String("foreign_key", device.ForeignKey),
		)

		return device, nil
	}

	device := &DeviceRegistration{
		UID:          uuid.New().String(),
		ForeignKey:   req.ForeignKey,
		Hostname:     req.Hostname,
		IPAddress:    req.IPAddress,
		DeviceType:   req.DeviceType,
		Platform:     req.Platform,
		Version:      req.Version,
		Labels:       req.Labels,
		RegisteredAt: time.Now(),
		LastSeen:     time.Now(),
	}

	r.devices[device.UID] = device
	r.byForeignKey[req.ForeignKey] = device.UID

	r.sendRegistrationEvent(device, "device_registered")

	r.logger.Info("Device registered",
		zap.String("uid", device.UID),
		zap.String("foreign_key", device.ForeignKey),
		zap.String("hostname", device.Hostname),
	)

	return device, nil
}

func (r *Registry) GetByUID(uid string) (*DeviceRegistration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[uid]
	return device, exists
}

func (r *Registry) GetByForeignKey(foreignKey string) (*DeviceRegistration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uid, exists := r.byForeignKey[foreignKey]
	if !exists {
		return nil, false
	}

	device, exists := r.devices[uid]
	return device, exists
}

func (r *Registry) UpdateLastSeen(uid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if device, exists := r.devices[uid]; exists {
		device.LastSeen = time.Now()
	}
}

func (r *Registry) ListDevices() []*DeviceRegistration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*DeviceRegistration, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}
	return devices
}

func (r *Registry) Unregister(uid string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	device, exists := r.devices[uid]
	if !exists {
		return false
	}

	delete(r.byForeignKey, device.ForeignKey)
	delete(r.devices, uid)

	r.sendRegistrationEvent(device, "device_unregistered")

	r.logger.Info("Device unregistered",
		zap.String("uid", uid),
		zap.String("foreign_key", device.ForeignKey),
	)

	return true
}

func (r *Registry) sendRegistrationEvent(device *DeviceRegistration, eventType string) {
	eventData := map[string]interface{}{
		"uid":           device.UID,
		"foreign_key":   device.ForeignKey,
		"hostname":      device.Hostname,
		"ip_address":    device.IPAddress,
		"device_type":   device.DeviceType,
		"platform":      device.Platform,
		"version":       device.Version,
		"labels":        device.Labels,
		"registered_at": device.RegisteredAt,
		"last_seen":     device.LastSeen,
		"event_type":    eventType,
	}

	if err := r.splunkClient.SendEvent(eventType, eventData); err != nil {
		r.logger.Warn("Failed to send registration event to Splunk",
			zap.String("uid", device.UID),
			zap.String("event_type", eventType),
			zap.Error(err),
		)
	}
}
