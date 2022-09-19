package applens

import (
	"github.com/Azure/aks-periscope/pkg/interfaces"
	"github.com/Azure/aks-periscope/pkg/utils"
)

// DNSCollector defines a DNS Collector struct
type ApplensSignalProducer struct {
	Parameters   map[string]string
	Signal       interfaces.Signal
	OsIdentifier utils.OSIdentifier
	FilePaths    *utils.KnownFilePaths
	FileSystem   interfaces.FileSystemAccessor
}

// NewApplensSignalProducer is a constructor
func NewApplensSignalProducer(parameters map[string]string, osIdentifier utils.OSIdentifier, filePaths *utils.KnownFilePaths, fileSystem interfaces.FileSystemAccessor) *ApplensSignalProducer {
	return &ApplensSignalProducer{
		Parameters:   parameters,
		OsIdentifier: osIdentifier,
		FilePaths:    filePaths,
		FileSystem:   fileSystem,
	}
}

func (producer *ApplensSignalProducer) GetName() string {
	return "applens"
}

func (producer *ApplensSignalProducer) GetData() map[string]interfaces.DataValue {
	return map[string]interfaces.DataValue{
		"rawResponse": utils.NewStringDataValue(apiResponse),
	}
}

func (producer *ApplensSignalProducer) GetParameters() map[string]string {
	return producer.Parameters
}

func (producer *ApplensSignalProducer) GetSignal() interfaces.Signal {
	return producer.Signal
}

func (producer *ApplensSignalProducer) Evaluate() error {
	// TODO call applens API for a detector

	producer.Signal = NewApplensSignal
	return nil
}
