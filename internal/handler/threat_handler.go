package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"threat-intelligence/internal/domain"
	"threat-intelligence/internal/service"
)

type ThreatHandler struct {
	service service.ThreatService
}

func NewThreatHandler(s service.ThreatService) *ThreatHandler {
	return &ThreatHandler{service: s}
}

func (h *ThreatHandler) GetByMD5(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	md5 := r.PathValue("md5")
	threat, err := h.service.GetByMD5(ctx, md5)
	if err != nil {
		if errors.Is(err, domain.ErrThreatNotFound) {
			http.Error(w, "threat not found", http.StatusNotFound)
			return
		} else if errors.Is(err, domain.ErrBadMD5) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(threat)
}
