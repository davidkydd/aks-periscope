package exporter

import (
	"github.com/Azure/aks-periscope/pkg/interfaces"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// LocalMachineExporter defines an Local Machine Exporter
type LocalMachineExporter struct{
	hostname     string
}

func NewLocalMachineExporter(hostname string) *LocalMachineExporter {
	return &LocalMachineExporter{
		hostname:     hostname,
	}
}

func (exporter *LocalMachineExporter) GetName() string {
	return "localmachine"
}

// Export implements the interface method
func (exporter *LocalMachineExporter) Export(_ interfaces.DataProducer) error {
	//TODO perhaps this should export to files on the localMachine also
	return nil
}

//ExportReader implements the interface method
func (exporter *LocalMachineExporter) ExportReader(name string, reader io.ReadSeeker) error {

	//TODO might need to change the location this file is written
	outFile, err := os.Create(name)
	defer outFile.Close()

	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("read from reader: %v",  err)
	}

	_, err = outFile.Write(bytes)
	if err != nil {
		log.Printf("write to outfile: %v",  err)
		return err
	}

	return nil
}