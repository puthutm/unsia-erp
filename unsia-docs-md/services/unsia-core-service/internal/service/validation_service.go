package service

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"
)

type ValidationService struct {
	validate *validator.Validate
}

func NewValidationService() *ValidationService {
	return &ValidationService{
		validate: validator.New(),
	}
}

// ValidateStruct validates a struct
func (s *ValidationService) ValidateStruct(i interface{}) error {
	return s.validate.Struct(i)
}

// ValidationError represents validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResponse represents validation response
type ValidationResponse struct {
	Valid bool            `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// ValidateAndRespond validates and returns response
func (s *ValidationService) ValidateAndRespond(i interface{}) ValidationResponse {
	err := s.validate.Struct(i)
	if err == nil {
		return ValidationResponse{Valid: true}
	}

	errors := make([]ValidationError, 0)
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: formatValidationMessage(err),
		})
	}

	return ValidationResponse{
		Valid: false,
		Errors: errors,
	}
}

func formatValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Field must be at least " + err.Param() + " characters"
	case "max":
		return "Field must be at most " + err.Param() + " characters"
	case "eqfield":
		return "Field must match " + err.Param()
	case "gte":
		return "Field must be greater than or equal to " + err.Param()
	case "lte":
		return "Field must be less than or equal to " + err.Param()
	default:
		return "Invalid value"
	}
}

// EmailValidator validates email format
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(strings.ToLower(email))
}

// PhoneValidator validates phone number (Indonesian format)
var phoneRegex = regexp.MustCompile(`^(\+62|62|0)[0-9]{9,12}$`)

func IsValidPhone(phone string) bool {
	if phone == "" {
		return false
	}
	return phoneRegex.MatchString(phone)
}

// NIMValidator validates NIM (Nomor Induk Mahasiswa)
func IsValidNIM(nim string) bool {
	if nim == "" {
		return false
	}
	// NIM should be 10-12 digits
	matched, _ := regexp.MatchString(`^[0-9]{10,12}$`, nim)
	return matched
}

// NIPValidator validates NIP (Nomor Induk Pegawai)
func IsValidNIP(nip string) bool {
	if nip == "" {
		return false
	}
	// NIP should be 18 digits (for PNS) or can be shorter for contract
	matched, _ := regexp.MatchString(`^[0-9]{9,18}$`, nip)
	return matched
}

// SendEmail sends email (simplified - in production use proper mail server)
type EmailService struct {
	smtpHost string
	smtpPort int
	from     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost: "smtp.gmail.com",
		smtpPort: 587,
		from:     "noreply@unsia.ac.id",
	}
}

// SendEmail sends an email
func (e *EmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// In production, configure proper SMTP credentials
	// d := gomail.Dialer{Host: e.smtpHost, Port: e.smtpPort}
	// return d.DialAndSend(m, ...)

	return nil
}

// SendWelcomeEmail sends welcome email to new user
func (e *EmailService) SendWelcomeEmail(email, name string) error {
	subject := "Selamat Datang di UNSIA"
	body := `
		<html>
		<body>
			<h2>Halo ` + name + `,</h2>
			<p>Selamat! Akun UNSIA Anda berhasil dibuat.</p>
			<p>Silakan login untuk mulai menggunakan sistem.</p>
		</body>
		</html>
	`
	return e.SendEmail(email, subject, body)
}
