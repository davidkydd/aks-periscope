package applens

import (
	"github.com/Azure/aks-periscope/pkg/interfaces"
	"github.com/Azure/aks-periscope/pkg/utils"
)

// ApplensSignal defines an Applens Signal struct
type ApplensSignal struct {
	SignalValue        string
	SignalError        error
	FollowupCollectors map[string]map[string]string //collectors to run as a "followup" to this signal value, represented by map of collector names to map of collector parameter key-values
	FollowupDiagnosers map[string]map[string]string //diagnoers to run as a "followup" to this signal value, represented by map of diagnoser names to map of diagnoser parameter key-values
	FileSystem         interfaces.FileSystemAccessor
}

// NewApplensSignal is a constructor
func NewApplensSignal(followupCollectors map[string]map[string]string, followupDiagnosers map[string]map[string]string, filePaths *utils.KnownFilePaths, fileSystem interfaces.FileSystemAccessor) *ApplensSignal {
	return &ApplensSignal{
		FollowupCollectors: followupCollectors,
		FollowupDiagnosers: followupDiagnosers,
	}
}

func (signal *ApplensSignal) GetName() string {
	return "applensSignal"
}

func (signal *ApplensSignal) GetData() map[string]interfaces.DataValue {
	return map[string]interfaces.DataValue{
		"apiResponse": utils.NewStringDataValue(),
	}
}

func (signal *ApplensSignal) GetFollowupCollectors() map[string]map[string]string {
	return signal.FollowupCollectors
}

func (signal *ApplensSignal) GetFollowupDiagnosers() map[string]map[string]string {
	return signal.FollowupDiagnosers
}

func (signal *ApplensSignal) CurrentSignal() (string, error) {
	return signal.SignalValue, nil
}
