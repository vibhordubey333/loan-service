package service

import (
	"fmt"

	"vibhordubey333/loan-service/internal/domain"
)

type PDFService interface {
	GenerateAgreementLetter(loan *domain.Loan) (string, error)
}

type pdfService struct{}

func NewPDFService() PDFService {
	return &pdfService{}
}

func (s *pdfService) GenerateAgreementLetter(loan *domain.Loan) (string, error) {
	//TODO: Implement wkhtmltopdf

	//Returning mock
	return fmt.Sprintf("https://vibhordubey.com/agreements/%s.pdf", loan.ID), nil
}
