package applens

// ApplensSignal defines an Applens Signal struct
type ApplensSignalError struct {
	Message     string
	SignalError error
}

func (e *ApplensSignalError) Error() string { return e.Message + ": " + e.SignalError.Error() }

func (e *ApplensSignalError) Unwrap() error { return e.SignalError }

// NewApplensSignal is a constructor
func NewApplensSignalError(message string, signalError error) *ApplensSignalError {
	return &ApplensSignalError{
		Message:     message,
		SignalError: signalError,
	}
}
