package interfaces

// Exporter defines interface for an exporter
type Exporter interface {
	GetName() string

	Export(DataProducer) error
}
