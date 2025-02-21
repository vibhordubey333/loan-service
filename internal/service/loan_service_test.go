package service

import (
	"context"
	"testing"
	_ "time"

	"vibhordubey333/loan-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLoanRepository struct {
	mock.Mock
}

type MockPDFService struct {
	mock.Mock
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockLoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func (m *MockLoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) AddInvestment(ctx context.Context, investment *domain.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockEmailService) SendInvestmentAgreement(investorID uuid.UUID, agreementURL string) error {
	args := m.Called(investorID, agreementURL)
	return args.Error(0)
}

func (m *MockPDFService) GenerateAgreementLetter(loan *domain.Loan) (string, error) {
	args := m.Called(loan)
	return args.String(0), args.Error(1)
}

func TestCreateLoan(t *testing.T) {
	repo := new(MockLoanRepository)
	emailService := new(MockEmailService)
	pdfService := new(MockPDFService)
	service := NewLoanService(repo, emailService, pdfService)

	ctx := context.Background()
	borrowerID := "12345"
	amount := 1000.0
	rate := 5.0
	roi := 8.0

	repo.On("Create", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	loan, err := service.CreateLoan(ctx, borrowerID, amount, rate, roi)

	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, borrowerID, loan.BorrowerIDNumber)
	assert.Equal(t, amount, loan.PrincipalAmount)
	assert.Equal(t, rate, loan.Rate)
	assert.Equal(t, roi, loan.ROI)
	assert.Equal(t, domain.LoanStateProposed, loan.State)

	repo.AssertExpectations(t)
}

func TestApproveLoan(t *testing.T) {
	repo := new(MockLoanRepository)
	emailService := new(MockEmailService)
	pdfService := new(MockPDFService)
	service := NewLoanService(repo, emailService, pdfService)

	ctx := context.Background()
	loanID := uuid.New()
	validatorID := "V123"
	proofImageURL := "https://example.com/proof.jpg"

	existingLoan := &domain.Loan{
		ID:    loanID,
		State: domain.LoanStateProposed,
	}

	repo.On("GetByID", ctx, loanID).Return(existingLoan, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	err := service.ApproveLoan(ctx, loanID, validatorID, proofImageURL)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestInvestInLoan(t *testing.T) {
	repo := new(MockLoanRepository)
	emailService := new(MockEmailService)
	pdfService := new(MockPDFService)
	service := NewLoanService(repo, emailService, pdfService)

	ctx := context.Background()
	loanID := uuid.New()
	investorID := uuid.New()
	amount := 1000.0

	existingLoan := &domain.Loan{
		ID:              loanID,
		State:           domain.LoanStateApproved,
		PrincipalAmount: 1000.0,
	}

	repo.On("GetByID", ctx, loanID).Return(existingLoan, nil)
	repo.On("AddInvestment", ctx, mock.AnythingOfType("*domain.Investment")).Return(nil)
	repo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	pdfService.On("GenerateAgreementLetter", mock.AnythingOfType("*domain.Loan")).Return("https://example.com/agreement.pdf", nil)
	emailService.On("SendInvestmentAgreement", investorID, mock.AnythingOfType("string")).Return(nil)

	err := service.InvestInLoan(ctx, loanID, investorID, amount)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	emailService.AssertExpectations(t)
	pdfService.AssertExpectations(t)
}

func TestDisburseLoan(t *testing.T) {
	repo := new(MockLoanRepository)
	emailService := new(MockEmailService)
	pdfService := new(MockPDFService)
	service := NewLoanService(repo, emailService, pdfService)

	ctx := context.Background()
	loanID := uuid.New()
	officerID := "O123"
	signedAgreementURL := "https://example.com/signed.pdf"

	existingLoan := &domain.Loan{
		ID:    loanID,
		State: domain.LoanStateInvested,
	}

	repo.On("GetByID", ctx, loanID).Return(existingLoan, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	err := service.DisburseLoan(ctx, loanID, officerID, signedAgreementURL)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
