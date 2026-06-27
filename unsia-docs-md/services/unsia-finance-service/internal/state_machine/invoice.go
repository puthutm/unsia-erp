package state_machine

import (
	"fmt"
)

// InvoiceStatus represents valid invoice statuses
type InvoiceStatus string

const (
	InvoiceStatusDraft           InvoiceStatus = "DRAFT"
	InvoiceStatusIssued          InvoiceStatus = "ISSUED"
	InvoiceStatusPartiallyPaid    InvoiceStatus = "PARTIALLY_PAID"
	InvoiceStatusPaid            InvoiceStatus = "PAID"
	InvoiceStatusCancelled       InvoiceStatus = "CANCELLED"
	InvoiceStatusExpired         InvoiceStatus = "EXPIRED"
)

// InvoiceStateMachine defines valid state transitions for Invoice
type InvoiceStateMachine struct{}

// NewInvoiceStateMachine creates a new invoice state machine
func NewInvoiceStateMachine() *InvoiceStateMachine {
	return &InvoiceStateMachine{}
}

// Valid transitions map
var invoiceTransitions = map[InvoiceStatus][]InvoiceStatus{
	InvoiceStatusDraft: {
		InvoiceStatusIssued,
		InvoiceStatusCancelled,
	},
	InvoiceStatusIssued: {
		InvoiceStatusPartiallyPaid,
		InvoiceStatusPaid,
		InvoiceStatusCancelled,
		InvoiceStatusExpired,
	},
	InvoiceStatusPartiallyPaid: {
		InvoiceStatusPaid,
		InvoiceStatusCancelled,
		InvoiceStatusExpired,
	},
	InvoiceStatusPaid: {}, // No further transitions
	InvoiceStatusCancelled: {}, // No further transitions
	InvoiceStatusExpired: {}, // No further transitions
}

// CanTransition checks if a status transition is valid
func (sm *InvoiceStateMachine) CanTransition(from, to InvoiceStatus) bool {
	allowedTransitions, exists := invoiceTransitions[from]
	if !exists {
		return false
	}

	for _, valid := range allowedTransitions {
		if valid == to {
			return true
		}
	}

	return false
}

// ValidateTransition returns an error if the transition is invalid
func (sm *InvoiceStateMachine) ValidateTransition(from, to InvoiceStatus) error {
	if from == to {
		return nil // No change = valid
	}

	if !sm.CanTransition(from, to) {
		return fmt.Errorf("invalid status transition from %s to %s", from, to)
	}

	return nil
}

// GetValidTransitions returns all valid transitions from a given status
func (sm *InvoiceStateMachine) GetValidTransitions(from InvoiceStatus) []InvoiceStatus {
	return invoiceTransitions[from]
}

// IsTerminal checks if the status is a terminal state (no further transitions allowed)
func (sm *InvoiceStateMachine) IsTerminal(status InvoiceStatus) bool {
	return len(invoiceTransitions[status]) == 0
}

// UpdateStatusWithTransition updates invoice status with validation
func (sm *InvoiceStateMachine) UpdateStatusWithTransition(currentStatus *string, newStatus string) error {
	current := InvoiceStatus(*currentStatus)
	to := InvoiceStatus(newStatus)

	if err := sm.ValidateTransition(current, to); err != nil {
		return err
	}

	*currentStatus = newStatus
	return nil
}

// PaymentStatus represents valid payment statuses
type PaymentStatus string

const (
	PaymentStatusReceived  PaymentStatus = "RECEIVED"
	PaymentStatusVerified PaymentStatus = "VERIFIED"
	PaymentStatusPosted    PaymentStatus = "POSTED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusReversed PaymentStatus = "REVERSED"
	// Legacy statuses from existing model
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
)

// PaymentStateMachine defines valid state transitions for Payment
type PaymentStateMachine struct{}

// NewPaymentStateMachine creates a new payment state machine
func NewPaymentStateMachine() *PaymentStateMachine {
	return &PaymentStateMachine{}
}

// Payment valid transitions
var paymentTransitions = map[PaymentStatus][]PaymentStatus{
	PaymentStatusReceived: {
		PaymentStatusVerified,
		PaymentStatusFailed,
	},
	PaymentStatusVerified: {
		PaymentStatusPosted,
		PaymentStatusReversed,
	},
	PaymentStatusPosted: {
		PaymentStatusReversed,
	},
	PaymentStatusFailed: {}, // No further transitions from failed
	PaymentStatusReversed: {}, // Terminal state
}

// CanTransition checks if a payment status transition is valid
func (sm *PaymentStateMachine) CanTransition(from, to PaymentStatus) bool {
	allowedTransitions, exists := paymentTransitions[from]
	if !exists {
		return false
	}

	for _, valid := range allowedTransitions {
		if valid == to {
			return true
		}
	}

	return false
}

// ValidateTransition returns an error if the transition is invalid
func (sm *PaymentStateMachine) ValidateTransition(from, to PaymentStatus) error {
	if from == to {
		return nil
	}

	if !sm.CanTransition(from, to) {
		return fmt.Errorf("invalid payment status transition from %s to %s", from, to)
	}

	return nil
}

// ClearanceStatus represents valid clearance statuses
type ClearanceStatus string

const (
	ClearanceStatusBlocked    ClearanceStatus = "BLOCKED"
	ClearanceStatusConditional ClearanceStatus = "CONDITIONAL"
	ClearanceStatusCleared    ClearanceStatus = "CLEARED"
	ClearanceStatusRevoked    ClearanceStatus = "REVOKED"
)

// ClearanceStateMachine defines valid state transitions for Clearance
type ClearanceStateMachine struct{}

// NewClearanceStateMachine creates a new clearance state machine
func NewClearanceStateMachine() *ClearanceStateMachine {
	return &ClearanceStateMachine{}
}

// Clearance valid transitions
var clearanceTransitions = map[ClearanceStatus][]ClearanceStatus{
	ClearanceStatusBlocked: {
		ClearanceStatusConditional,
		ClearanceStatusCleared,
	},
	ClearanceStatusConditional: {
		ClearanceStatusCleared,
		ClearanceStatusBlocked,
	},
	ClearanceStatusCleared: {
		ClearanceStatusRevoked,
	},
	ClearanceStatusRevoked: {}, // Terminal state
}

// CanTransition checks if a clearance status transition is valid
func (sm *ClearanceStateMachine) CanTransition(from, to ClearanceStatus) bool {
	allowedTransitions, exists := clearanceTransitions[from]
	if !exists {
		return false
	}

	for _, valid := range allowedTransitions {
		if valid == to {
			return true
		}
	}

	return false
}

// ValidateTransition returns an error if the transition is invalid
func (sm *ClearanceStateMachine) ValidateTransition(from, to ClearanceStatus) error {
	if from == to {
		return nil
	}

	if !sm.CanTransition(from, to) {
		return fmt.Errorf("invalid clearance status transition from %s to %s", from, to)
	}

	return nil
}

// ValidateTransitionWithReason validates that REVOKED status requires a reason
func (sm *ClearanceStateMachine) ValidateTransitionWithReason(from, to ClearanceStatus, reason *string) error {
	if err := sm.ValidateTransition(from, to); err != nil {
		return err
	}

	// REVOKED status requires a reason
	if to == ClearanceStatusRevoked {
		if reason == nil || *reason == "" {
			return fmt.Errorf("reason is required when revoking clearance")
		}
	}

	return nil
}
