package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/Azure/aks-periscope/pkg/collector"
	"github.com/Azure/aks-periscope/pkg/diagnoser"
	"github.com/Azure/aks-periscope/pkg/exporter"
	"github.com/Azure/aks-periscope/pkg/interfaces"
	"github.com/Azure/aks-periscope/pkg/utils"
	restclient "k8s.io/client-go/rest"
)

func main() {
	osIdentifier, err := utils.StringToOSIdentifier(runtime.GOOS)
	if err != nil {
		log.Fatalf("cannot determine OS: %v", err)
	}

	knownFilePaths, err := utils.GetKnownFilePaths(osIdentifier)
	if err != nil {
		log.Fatalf("failed to get file paths: %v", err)
	}

	fileSystem := utils.NewFileSystem()

	// Create a watcher for the Run ID file that checks its content every 10 seconds
	fileWatcher := utils.NewFileContentWatcher(fileSystem, 10*time.Second)

	// Create a channel for unrecoverable errors
	errChan := make(chan error)

	// Add a watcher for the run ID file content
	runIdChan := make(chan string)
	fileWatcher.AddHandler(knownFilePaths.GetConfigPath(utils.RunIdKey), runIdChan, errChan)

	go func() {
		for {
			runId := <-runIdChan
			log.Printf("Starting Periscope run %s", runId)
			err := run(osIdentifier, knownFilePaths, fileSystem)
			if err != nil {
				errChan <- err
			}

			log.Printf("Completed Periscope run %s", runId)
		}
	}()

	fileWatcher.Start()

	// Run until unrecoverable error
	err = <-errChan
	log.Fatalf("Error running Periscope: %v", err)
}

func run(osIdentifier utils.OSIdentifier, knownFilePaths *utils.KnownFilePaths, fileSystem interfaces.FileSystemAccessor) error {
	runtimeInfo, err := utils.GetRuntimeInfo(fileSystem, knownFilePaths)
	if err != nil {
		log.Fatalf("Failed to get runtime information: %v", err)
	}

	config, err := restclient.InClusterConfig()
	if err != nil {
		return fmt.Errorf("cannot load kubeconfig: %w", err)
	}

	// Copies self-signed cert information to container if application is running on Azure Stack Cloud.
	// We need the cert in order to communicate with the storage account.
	if utils.IsAzureStackCloud(knownFilePaths) {
		if err := utils.CopyFile(knownFilePaths.AzureStackCertHost, knownFilePaths.AzureStackCertContainer); err != nil {
			return fmt.Errorf("cannot copy cert for Azure Stack Cloud environment: %w", err)
		}
	}

	//collectorGrp := new(sync.WaitGroup)
	collectorMap, diagnoserMap, exporterMap := initializeComponents(osIdentifier, knownFilePaths, fileSystem, config, runtimeInfo)

	dataProducers := []interfaces.DataProducer{}

	/* don't run collectors directly anymore
	for _, c := range collectors {
		if err := c.CheckSupported(); err != nil {
			// Log the reason why this collector is not supported, and skip to the next
			log.Printf("Skipping unsupported collector %s: %v", c.GetName(), err)
			continue
		}

		dataProducers = append(dataProducers, c)
		collectorGrp.Add(1)
		go func(c interfaces.Collector) {
			defer collectorGrp.Done()

			log.Printf("Collector: %s, collect data", c.GetName())
			err := c.Collect()
			if err != nil {
				log.Printf("Collector: %s, collect data failed: %v", c.GetName(), err)
				return
			}

			log.Printf("Collector: %s, export data", c.GetName())
			if err = exp.Export(c); err != nil {
				log.Printf("Collector: %s, export data failed: %v", c.GetName(), err)
			}
		}(c)
	}

	collectorGrp.Wait()*/

	diagnosers := []interfaces.Diagnoser{
		diagnoser.NewNetworkConfigDiagnoser(runtimeInfo, dnsCollector, kubeletCmdCollector),
		diagnoser.NewNetworkOutboundDiagnoser(runtimeInfo, networkOutboundCollector),
	}

	diagnoserGrp := new(sync.WaitGroup)

	for _, d := range diagnosers {
		dataProducers = append(dataProducers, d)
		diagnoserGrp.Add(1)
		go func(d interfaces.Diagnoser) {
			defer diagnoserGrp.Done()

			log.Printf("Diagnoser: %s, diagnose data", d.GetName())
			diagnosers, err := d.Diagnose()
			if err != nil {
				log.Printf("Diagnoser: %s, diagnose data failed: %v", d.GetName(), err)
				return
			}

			log.Printf("Diagnoser: %s, export data", d.GetName())
			if err = exp.Export(d); err != nil {
				log.Printf("Diagnoser: %s, export data failed: %v", d.GetName(), err)
			}
		}(d)
	}

	diagnoserGrp.Wait()

	zip, err := exporter.Zip(dataProducers)
	if err != nil {
		log.Printf("Could not zip data: %v", err)
	} else {
		if err := exp.ExportReader(runtimeInfo.HostNodeName+".zip", bytes.NewReader(zip.Bytes())); err != nil {
			log.Printf("Could not export zip archive: %v", err)
		}
	}

	return nil
}

func runDiagnoser(diagnoser interfaces.Diagnoser) ([]interfaces.Collector, []interfaces.Diagnoser) {

	//DataStructure changes:
	//We need to introduce a new struct representing a "signal" coming back from a collector / diagnoser.
	//that includes the name of the collector / diagnoser + the params to run it with
	//then we maintain a single map of "executable" IDs to the collector / diagnoser object they should run so we can look them up.

	//Logical changes:
	//call runDiagnoser / runCollector in a loop, collecting the followupCollectors / followupDiagnosers that are returned in an aggregated collection
	//before executing these in the next pass through the loop, until we run out of returned collectors / detectors.
	//Each pass through the loop will correspond to an additional level of "depth" in the diagnostic tree, e.g. breadth first execution.

	//Need to include a max depth to make sure it doesn't run forever

	//followupCollectors, followupDiagnosers, err := diagnoser.Diagnose() //we need to add an API like the below
}

// initializeComponents initializes and returns collectors, diagnosers and exporters
func initializeComponents(
	osIdentifier utils.OSIdentifier,
	knownFilePaths *utils.KnownFilePaths,
	fileSystem interfaces.FileSystemAccessor,
	config *restclient.Config,
	runtimeInfo *utils.RuntimeInfo) (map[string]interfaces.Collector, map[string]interfaces.Diagnoser, map[string]interfaces.Exporter) {

	// exporters
	azureBlobExporter := exporter.NewAzureBlobExporter(runtimeInfo, knownFilePaths, runtimeInfo.RunId)
	exporters := map[string]interfaces.Exporter{
		azureBlobExporter.GetName(): azureBlobExporter,
	}

	// collectors
	dnsCollector := collector.NewDNSCollector(osIdentifier, knownFilePaths, fileSystem)
	helmCollector := collector.NewHelmCollector(config, runtimeInfo)
	ipTablesCollector := collector.NewIPTablesCollector(osIdentifier, runtimeInfo)
	kubeletCmdCollector := collector.NewKubeletCmdCollector(osIdentifier, runtimeInfo)
	kubeObjectsCollector := collector.NewKubeObjectsCollector(config, runtimeInfo)
	networkOutboundCollector := collector.NewNetworkOutboundCollector()
	nodeLogsCollector := collector.NewNodeLogsCollector(runtimeInfo, fileSystem)
	osmCollector := collector.NewOsmCollector(config, runtimeInfo)
	pdbCollector := collector.NewPDBCollector(config, runtimeInfo)
	podsContainerLogsCollector := collector.NewPodsContainerLogsCollector(config, runtimeInfo)
	smiCollector := collector.NewSmiCollector(config, runtimeInfo)
	systemLogsCollector := collector.NewSystemLogsCollector(osIdentifier, runtimeInfo)
	systemPerfCollector := collector.NewSystemPerfCollector(config, runtimeInfo)
	windowsLogsCollector := collector.NewWindowsLogsCollector(osIdentifier, runtimeInfo, knownFilePaths, fileSystem, 10*time.Second, 20*time.Minute)

	collectors := map[string]interfaces.Collector{
		dnsCollector.GetName():               dnsCollector,
		helmCollector.GetName():              helmCollector,
		ipTablesCollector.GetName():          ipTablesCollector,
		kubeletCmdCollector.GetName():        kubeletCmdCollector,
		kubeObjectsCollector.GetName():       kubeObjectsCollector,
		networkOutboundCollector.GetName():   networkOutboundCollector,
		nodeLogsCollector.GetName():          nodeLogsCollector,
		osmCollector.GetName():               osmCollector,
		pdbCollector.GetName():               pdbCollector,
		podsContainerLogsCollector.GetName(): podsContainerLogsCollector,
		smiCollector.GetName():               smiCollector,
		systemLogsCollector.GetName():        systemLogsCollector,
		systemPerfCollector.GetName():        systemPerfCollector,
		windowsLogsCollector.GetName():       windowsLogsCollector,
	}

	//diagnosers
	networkConfigDiagnoser := diagnoser.NewNetworkConfigDiagnoser(runtimeInfo, dnsCollector, kubeletCmdCollector)
	networkOutboundDiagnoser := diagnoser.NewNetworkOutboundDiagnoser(runtimeInfo, networkOutboundCollector)

	diagnosers := map[string]interfaces.Diagnoser{
		networkConfigDiagnoser.GetName():   networkConfigDiagnoser,
		networkOutboundDiagnoser.GetName(): networkOutboundDiagnoser,
	}

	return collectors, diagnosers, exporters
}

// selectedExporters select the exporters to run
// func selectExporters(allExporters map[string]interfaces.Exporter) []interfaces.Exporter {
// 	exporters := []interfaces.Exporter{}

// 	//read list of collectors that are enabled
// 	enabledExporterNames := strings.Fields(os.Getenv("ENABLED_EXPORTERS"))

// 	for _, exporter := range enabledExporterNames {
// 		exporters = append(exporters, allExporters[exporter])
// 	}

// 	return exporters
// }
