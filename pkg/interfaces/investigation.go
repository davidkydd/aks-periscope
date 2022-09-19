package interfaces

// Investigation defines interface for an investigation
type Investigation interface {
	GetName() string

	GetParameters() map[string]string

	GetDataProducer() DataProducer

	GetSignal() Signal
}
