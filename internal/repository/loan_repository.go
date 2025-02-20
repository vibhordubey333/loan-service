package repository

import (
	"context"
	"database/sql"
	"vibhordubey333/loan-service/internal/domain"
	"github.com/google/uuid"
	"encoding/json"
	"errors"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *domain.Loan) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error)
	Update(ctx context.Context, loan *domain.Loan) error
}

type loanRepository struct {
	db *sql.DB
}

func (r *loanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	//TODO implement me
	panic("implement me")
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO loans (
			id, borrower_id_number, principal_amount, rate, roi, 
			state, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		loan.ID, loan.BorrowerIDNumber, loan.PrincipalAmount,
		loan.Rate, loan.ROI, loan.State, loan.CreatedAt, loan.UpdatedAt,
	).Scan(&loan.ID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *loanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	approvalJSON, err := json.Marshal(loan.ApprovalDetails)
	if err != nil {
		return err
	}

	disbursementJSON, err := json.Marshal(loan.DisbursementDetails)
	if err != nil {
		return err
	}

	query := `
		UPDATE loans 
		SET state = $1,
			approval_details = $2,
			disbursement_details = $3,
			agreement_letter_url = $4,
			updated_at = $5
		WHERE id = $6`

	result, err := tx.ExecContext(ctx, query,
		loan.State,
		approvalJSON,
		disbursementJSON,
		loan.AgreementLetterURL,
		loan.UpdatedAt,
		loan.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("loan not found")
	}

	return tx.Commit()
}
