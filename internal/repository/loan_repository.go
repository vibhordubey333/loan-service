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
	AddInvestment(ctx context.Context, investment *domain.Investment) error
}

type loanRepository struct {
	db *sql.DB
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

func (r *loanRepository) AddInvestment(ctx context.Context, investment *domain.Investment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO investments (
			id, loan_id, investor_id, amount, created_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		investment.ID, investment.LoanID, investment.InvestorID,
		investment.Amount, investment.CreatedAt,
	).Scan(&investment.ID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *loanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	loan := &domain.Loan{}

	query := `
		SELECT 
			l.id, l.borrower_id_number, l.principal_amount, l.rate, 
			l.roi, l.state, l.created_at, l.updated_at,
			l.approval_details, l.disbursement_details, l.agreement_letter_url
		FROM loans l
		WHERE l.id = $1`

	var approvalJSON, disbursementJSON sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&loan.ID, &loan.BorrowerIDNumber, &loan.PrincipalAmount,
		&loan.Rate, &loan.ROI, &loan.State, &loan.CreatedAt, &loan.UpdatedAt,
		&approvalJSON, &disbursementJSON, &loan.AgreementLetterURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("loan not found")
		}
		return nil, err
	}

	// Unmarshal approval details if present
	if approvalJSON.Valid {
		var approval domain.ApprovalDetails
		if err := json.Unmarshal([]byte(approvalJSON.String), &approval); err != nil {
			return nil, err
		}
		loan.ApprovalDetails = &approval
	}

	// Unmarshal disbursement details if present
	if disbursementJSON.Valid {
		var disbursement domain.DisbursementDetails
		if err := json.Unmarshal([]byte(disbursementJSON.String), &disbursement); err != nil {
			return nil, err
		}
		loan.DisbursementDetails = &disbursement
	}

	// Get investments
	investmentsQuery := `
		SELECT id, loan_id, investor_id, amount, created_at
		FROM investments
		WHERE loan_id = $1`

	rows, err := r.db.QueryContext(ctx, investmentsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var inv domain.Investment
		if err := rows.Scan(&inv.ID, &inv.LoanID, &inv.InvestorID, &inv.Amount, &inv.CreatedAt); err != nil {
			return nil, err
		}
		loan.Investments = append(loan.Investments, inv)
	}

	return loan, nil
}
