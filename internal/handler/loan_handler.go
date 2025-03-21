package handler

import (
	"encoding/json"
	"net/http"

	"vibhordubey333/loan-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LoanHandler struct {
	service  service.LoanService
	validate *validator.Validate
}

func NewLoanHandler(service service.LoanService) *LoanHandler {
	return &LoanHandler{
		service:  service,
		validate: validator.New(),
	}
}

type CreateLoanRequest struct {
	BorrowerIDNumber string  `json:"borrower_id_number" validate:"required"`
	PrincipalAmount  float64 `json:"principal_amount" validate:"required,gt=0"`
	Rate             float64 `json:"rate" validate:"required,gt=0"`
	ROI              float64 `json:"roi" validate:"required,gt=0"`
}

type ApproveLoanRequest struct {
	FieldValidatorID string `json:"field_validator_id" validate:"required"`
	ProofImageURL    string `json:"proof_image_url" validate:"required,url"`
}

type InvestmentRequest struct {
	InvestorID uuid.UUID `json:"investor_id" validate:"required"`
	Amount     float64   `json:"amount" validate:"required,gt=0"`
}

type DisbursementRequest struct {
	FieldOfficerID     string `json:"field_officer_id" validate:"required"`
	SignedAgreementURL string `json:"signed_agreement_url" validate:"required,url"`
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loan, err := h.service.CreateLoan(r.Context(), req.BorrowerIDNumber,
		req.PrincipalAmount, req.Rate, req.ROI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loan)
}

func (h *LoanHandler) GetLoan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	loan, err := h.service.GetLoan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loan)
}

func (h *LoanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	var req ApproveLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.ApproveLoan(r.Context(), id, req.FieldValidatorID, req.ProofImageURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LoanHandler) InvestInLoan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	var req InvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.InvestInLoan(r.Context(), id, req.InvestorID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LoanHandler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	var req DisbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DisburseLoan(r.Context(), id, req.FieldOfficerID, req.SignedAgreementURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
