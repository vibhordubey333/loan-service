package service

import (
	"context"
	"time"

	"vibhordubey333/loan-service/internal/domain"
	"vibhordubey333/loan-service/internal/repository"

	"github.com/google/uuid"
	"errors"
	"log"
)

type LoanService interface {
	CreateLoan(ctx context.Context, borrowerID string, amount, rate, roi float64) (*domain.Loan, error)
	GetLoan(ctx context.Context, id uuid.UUID) (*domain.Loan, error)
	ApproveLoan(ctx context.Context, id uuid.UUID, validatorID, proofImageURL string) error
	InvestInLoan(ctx context.Context, loanID, investorID uuid.UUID, amount float64) error
	DisburseLoan(ctx context.Context, id uuid.UUID, officerID, signedAgreementURL string) error
}

type loanService struct {
	repo         repository.LoanRepository
	emailService EmailService
	pdfService   PDFService
}

func NewLoanService(repo repository.LoanRepository, emailService EmailService, pdfService PDFService) LoanService {
	return &loanService{
		repo:         repo,
		emailService: emailService,
		pdfService:   pdfService,
	}
}

func (s *loanService) ApproveLoan(ctx context.Context, id uuid.UUID, validatorID, proofImageURL string) error {
	loan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !loan.CanApprove() {
		return errors.New("loan cannot be approved in current state")
	}

	loan.State = domain.LoanStateApproved
	loan.ApprovalDetails = &domain.ApprovalDetails{
		FieldValidatorID: validatorID,
		ProofImageURL:    proofImageURL,
		ApprovedAt:       time.Now(),
	}
	loan.UpdatedAt = time.Now()

	return s.repo.Update(ctx, loan)
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

func (s *loanService) InvestInLoan(ctx context.Context, loanID, investorID uuid.UUID, amount float64) error {
	loan, err := s.repo.GetByID(ctx, loanID)
	if err != nil {
		return err
	}

	if !loan.CanInvest() {
		return errors.New("loan is not available for investment")
	}

	if loan.TotalInvestedAmount()+amount > loan.PrincipalAmount {
		return errors.New("investment amount exceeds remaining principal")
	}

	investment := &domain.Investment{
		ID:         uuid.New(),
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     amount,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.AddInvestment(ctx, investment); err != nil {
		log.Println("Error while adding investment")
		return err
	}

	loan.Investments = append(loan.Investments, *investment)

	if loan.IsFullyInvested() {
		log.Println("Loan is fully invested")
		loan.State = domain.LoanStateInvested
		loan.UpdatedAt = time.Now()

		//Todo: Generate agreement letter

		//Todo: Send emails to all investors
		agreementURL := "test.com"

		// Send emails to all investors
		for _, inv := range loan.Investments {
			if err := s.emailService.SendInvestmentAgreement(inv.InvestorID, agreementURL); err != nil {
				//Todo: Implement retry mechanism
				log.Println("Failed to send email to investor", inv.InvestorID)
				continue
			}
		}
		// Persist the updated loan state
		if err := s.repo.Update(ctx, loan); err != nil {
			return err
		}
	}

	return nil
}

func (s *loanService) DisburseLoan(ctx context.Context, id uuid.UUID, officerID, signedAgreementURL string) error {
	loan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !loan.CanDisburse() {
		return errors.New("loan cannot be disbursed in current state")
	}

	loan.State = domain.LoanStateDisbursed
	loan.DisbursementDetails = &domain.DisbursementDetails{
		FieldOfficerID:     officerID,
		SignedAgreementURL: signedAgreementURL,
		DisbursedAt:        time.Now(),
	}
	loan.UpdatedAt = time.Now()

	return s.repo.Update(ctx, loan)
}

func (s *loanService) GetLoan(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	return s.repo.GetByID(ctx, id)
}
