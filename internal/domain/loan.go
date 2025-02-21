package domain

import (
	"time"

	"github.com/google/uuid"
	"database/sql"
)

type LoanState string

const (
	LoanStateProposed  LoanState = "PROPOSED"
	LoanStateApproved  LoanState = "APPROVED"
	LoanStateInvested  LoanState = "INVESTED"
	LoanStateDisbursed LoanState = "DISBURSED"
)

type Loan struct {
	ID               uuid.UUID `json:"id"`
	BorrowerIDNumber string    `json:"borrower_id_number"`
	PrincipalAmount  float64   `json:"principal_amount"`
	Rate             float64   `json:"rate"`
	ROI              float64   `json:"roi"`
	State            LoanState `json:"state"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	ApprovalDetails *ApprovalDetails `json:"approval_details,omitempty"`
	Investments     []Investment     `json:"investments,omitempty"`

	DisbursementDetails *DisbursementDetails `json:"disbursement_details,omitempty"`
	AgreementLetterURL  sql.NullString       `json:"agreement_letter_url,omitempty"`
}

type ApprovalDetails struct {
	FieldValidatorID string    `json:"field_validator_id"`
	ProofImageURL    string    `json:"proof_image_url"`
	ApprovedAt       time.Time `json:"approved_at"`
}

type Investment struct {
	ID         uuid.UUID `json:"id"`
	LoanID     uuid.UUID `json:"loan_id"`
	InvestorID uuid.UUID `json:"investor_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type DisbursementDetails struct {
	FieldOfficerID     string    `json:"field_officer_id"`
	SignedAgreementURL string    `json:"signed_agreement_url"`
	DisbursedAt        time.Time `json:"disbursed_at"`
}

func (l *Loan) CanApprove() bool {
	return l.State == LoanStateProposed
}

func (l *Loan) CanInvest() bool {
	return l.State == LoanStateApproved
}

func (l *Loan) CanDisburse() bool {
	return l.State == LoanStateInvested
}

func (l *Loan) TotalInvestedAmount() float64 {
	var total float64
	for _, inv := range l.Investments {
		total += inv.Amount
	}
	return total
}

func (l *Loan) IsFullyInvested() bool {
	return l.TotalInvestedAmount() >= l.PrincipalAmount
}
