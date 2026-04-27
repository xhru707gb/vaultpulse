package secretexpiry

// Exported constants and types for use in external test packages.

const (
	StateOK      = stateOK
	StateWarning = stateWarning
	StateExpired = stateExpired
)

// Status is re-exported so alert_test.go can construct values directly.
type Status = status
