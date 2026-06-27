package handler

import "time"

type InvoiceBriefResponse struct {
	ID               string     `json:"id"`
	InvoiceNumber    string     `json:"invoice_number"`
	TargetType       string     `json:"target_type"` // applicant, student
	ApplicantID      *string    `json:"applicant_id,omitempty"`
	StudentID        *string    `json:"student_id,omitempty"`
	AcademicPeriodID *string    `json:"academic_period_id,omitempty"`
	TotalAmount      float64    `json:"total_amount"`
	PaidAmount       float64    `json:"paid_amount"`
	Status           string     `json:"status"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type InvoiceItemResponse struct {
	ID                 string  `json:"id"`
	PaymentComponentID *string `json:"payment_component_id,omitempty"`
	Description        string  `json:"description"`
	Amount             float64 `json:"amount"`
	DiscountAmount     float64 `json:"discount_amount"`
	FinalAmount        float64 `json:"final_amount"`
}

type PaymentResponse struct {
	ID              string     `json:"id"`
	PaymentMethodID *string    `json:"payment_method_id,omitempty"`
	PaymentNumber   *string    `json:"payment_number,omitempty"`
	Amount          float64    `json:"amount"`
	PaymentStatus   string     `json:"payment_status"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
}

type InvoiceDetailResponse struct {
	ID               string                 `json:"id"`
	InvoiceNumber    string                 `json:"invoice_number"`
	TargetType       string                 `json:"target_type"`
	ApplicantID      *string                `json:"applicant_id,omitempty"`
	StudentID        *string                `json:"student_id,omitempty"`
	AcademicPeriodID *string                `json:"academic_period_id,omitempty"`
	TotalAmount      float64                `json:"total_amount"`
	PaidAmount       float64                `json:"paid_amount"`
	Status           string                 `json:"status"`
	DueDate          *time.Time             `json:"due_date,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Items            []InvoiceItemResponse  `json:"items"`
	Payments         []PaymentResponse      `json:"payments"`
}

type ClearanceStatusResponse struct {
	StudentID        string   `json:"student_id"`
	AcademicPeriodID string   `json:"academic_period_id"`
	ServiceCode      string   `json:"service_code"`
	ClearanceStatus  string   `json:"clearance_status"`
	BlockReasons     []string `json:"block_reasons"`
}
