package service

import (
	"fmt"

	"gopkg.in/gomail.v2"
	"github.com/google/uuid"
)

type EmailService interface {
	SendInvestmentAgreement(investorID uuid.UUID, agreementURL string) error
}

type emailService struct {
	smtpConfig SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) EmailService {
	return &emailService{
		smtpConfig: config,
	}
}

func (s *emailService) SendInvestmentAgreement(investorID uuid.UUID, agreementURL string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.smtpConfig.Username)
	//TODO: Replace with actual email
	m.SetHeader("To", fmt.Sprintf("investor-%s@example.com", investorID))
	m.SetHeader("Subject", "Loan Investment Agreement")
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Investment Agreement</h1>
		<p>Please find your investment agreement at the following link:</p>
		<p><a href="%s">View Agreement</a></p>
	`, agreementURL))

	d := gomail.NewDialer(s.smtpConfig.Host, s.smtpConfig.Port,
		s.smtpConfig.Username, s.smtpConfig.Password)

	return d.DialAndSend(m)
}
