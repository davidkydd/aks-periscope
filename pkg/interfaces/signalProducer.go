package interfaces

// SignalProducer defines an object producing a signal
type SignalProducer interface {
	GetName() string

	Evaluate() error

	GetParameters() map[string]string

	GetData() map[string]DataValue

	GetSignal() Signal
}
