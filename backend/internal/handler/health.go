package handler

import "time"

type HealthRequest struct{}

type HealthResponse struct {
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Method() string  { return "GET" }
func (h *HealthHandler) Pattern() string { return "/health" }

func (h *HealthHandler) Validate(_ HealthRequest) error { return nil }

func (h *HealthHandler) Serve(_ HealthRequest) (HealthResponse, error) {
	return HealthResponse{
		Timestamp: time.Now().Unix(),
		Status:    "healthy",
	}, nil
}
