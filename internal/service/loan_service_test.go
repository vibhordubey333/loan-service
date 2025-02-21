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

func (m *MockLoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockLoanRepository) AddInvestment(ctx context.Context, investment *domain.Investment) error {
	//TODO implement me
	panic("implement me")
}

type MockEmailService struct {
	mock.Mock
}

func (m MockEmailService) SendInvestmentAgreement(investorID uuid.UUID, agreementURL string) error {
	//TODO implement me
	panic("implement me")
}

type MockPDFService struct {
	mock.Mock
}

func (m MockPDFService) GenerateAgreementLetter(loan *domain.Loan) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockLoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Loan), args.Error(1)
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
