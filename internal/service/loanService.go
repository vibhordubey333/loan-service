package service

import (
	"context"
	"time"

	"vibhordubey333/loan-service/internal/domain"
	"vibhordubey333/loan-service/internal/repository"

	"github.com/google/uuid"
)

type LoanService interface {
	CreateLoan(ctx context.Context, borrowerID string, amount, rate, roi float64) (*domain.Loan, error)
	GetLoan(ctx context.Context, id uuid.UUID) (*domain.Loan, error)
}

type loanService struct {
	repo repository.LoanRepository
}

func NewLoanService(repo repository.LoanRepository) LoanService {
	return &loanService{
		repo: repo,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, borrowerID string, amount, rate, roi float64) (*domain.Loan, error) {
	loan := &domain.Loan{
		ID:               uuid.New(),
		BorrowerIDNumber: borrowerID,
		PrincipalAmount:  amount,
		Rate:             rate,
		ROI:              roi,
		State:            domain.LoanStateProposed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.Create(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *loanService) GetLoan(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	return s.repo.GetByID(ctx, id)
}
