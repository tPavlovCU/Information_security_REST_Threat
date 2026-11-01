package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"threat-intelligence/internal/domain"
	"threat-intelligence/internal/service"
)

type CreateResult struct {
	MD5  string `json:"md5"`
	Name string `json:"name"`
}

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

func (h *ThreatHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data domain.Threat
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.service.Create(ctx, &data)
	if err != nil {
		if errors.Is(err, domain.ErrBadMD5) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.ErrBadConfidence) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.ErrBadSeverity) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.ErrBadType) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.ErrEmptyName) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.ErrThreatAlreadyExists) {
			http.Error(w, "Bad request", http.StatusConflict)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res CreateResult

	res.MD5 = data.MD5
	res.Name = data.Name

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
