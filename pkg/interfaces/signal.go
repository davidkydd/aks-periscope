package interfaces

// Signal defines interface for a signal
type Signal interface {
	CurrentSignal() (interface{}, interface{}, error)

	FollowupCollectors() map[string]map[string]string //map of collector names to map of collector parameter key-values

	FollowupDiagnosers() map[string]map[string]string //map of diagnoser names to map of diagnoser parameter key-values
}
