package interfaces

// Diagnoser defines interface for a diagnoser, conceptually identical to a collector but can be an entry point for starting a "diagnosis"
type Diagnoser interface {
	GetName() string

	Diagnose() error

	GetData() map[string]DataValue
}
