package collector

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/aks-periscope/pkg/interfaces"
	"github.com/Azure/aks-periscope/pkg/utils"
)

// KubeObjectsCollector defines a KubeObjects Collector struct
type KubeObjectsCollector struct {
	BaseCollector
}

var _ interfaces.Collector = &KubeObjectsCollector{}

// NewKubeObjectsCollector is a constructor
func NewKubeObjectsCollector(exporters []interfaces.Exporter) *KubeObjectsCollector {
	return &KubeObjectsCollector{
		BaseCollector: BaseCollector{
			collectorType: KubeObjects,
			exporters:     exporters,
		},
	}
}

type KubeObjectToCollect struct {
	Id string `yaml:"id"`
	Kubeobject string `yaml:"kubeobject"`
	Selector string `yaml:"selector"`
	Output string `yaml:"output"`
	Namespace string `yaml:"namespace"`
}

type KubeObjectsToCollect struct {
	Kubeobjects []KubeObjectToCollect `yaml:"kubeobjects_to_collect"`
}

// Collect implements the interface method
func (collector *KubeObjectsCollector) Collect() error {

	rootPath, err := utils.CreateCollectorDir(collector.GetName())
	//if err != nil {
		//return err
	//}[5]int{10, 20, 30, 40, 50}

	properties := os.Getenv("kubeobjects_to_collect")
	log.Printf("Kubeobject properties : %s\n", properties)

	var objects = &KubeObjectsToCollect{
		Kubeobjects: []KubeObjectToCollect{
			KubeObjectToCollect{Id: "test", Kubeobject: "test2", Selector: "test3", Output: "test4", Namespace: "test5"},
			KubeObjectToCollect{Id: "test6", Kubeobject: "test7", Selector: "test8", Output: "test9", Namespace: "test10"},
		},
	}

	marshaled, err := yaml.Marshal(&objects)
	if err != nil {
		return err
	}

	log.Printf("Marshaled properties : %s\n", marshaled)

	var kubeObjectsToCollect KubeObjectsToCollect
	err = yaml.Unmarshal([]byte(properties), &kubeObjectsToCollect)
	if err != nil {
		panic(err)
	}

	fmt.Printf("KubeObjectsToCollect Value: %#v\n", kubeObjectsToCollect)

	for _, ko := range kubeObjectsToCollect.Kubeobjects {
		output, err := utils.RunCommandOnContainer("kubectl", "-n", ko.Namespace, "get", ko.Kubeobject, "--output="+ko.Output, "--selector="+ko.Selector)
		if err != nil {
			return err
		}
		objects := []string{}
		objects = strings.Split(output, " ")

		for _, object := range objects {

			//in the case that both the original and new implementations are executed, and the same object is
			//to be "collected" by both implementations, the call to utils.writetofile will be made twice
			//for the same filename, this should however just result in the file being truncated (overwritten)
			//by whichever runs second - due to utils.WriteToFile calling os.Create
			kubernetesObjectFile := filepath.Join(rootPath, ko.Namespace+"_"+ko.Kubeobject+"_"+object)

			output, err := utils.RunCommandOnContainer("kubectl", "-n", ko.Namespace, "describe", ko.Kubeobject, object)
			if err != nil {
				return err
			}

			err = utils.WriteToFile(kubernetesObjectFile, output)
			if err != nil {
				return err
			}

			collector.AddToCollectorFiles(kubernetesObjectFile)
		}
	}

	//original implementation
	kubernetesObjects := strings.Fields(os.Getenv("DIAGNOSTIC_KUBEOBJECTS_LIST"))
	log.Printf("KubernetesObjects: %s\n", kubernetesObjects)

	for _, kubernetesObject := range kubernetesObjects {
		kubernetesObjectParts := strings.Split(kubernetesObject, "/")
		nameSpace := kubernetesObjectParts[0]
		objectType := kubernetesObjectParts[1]
		objects := []string{}
		if len(kubernetesObjectParts) == 3 {
			objects = append(objects, kubernetesObjectParts[2])
		}

		if len(objects) == 0 {
			output, err := utils.RunCommandOnContainer("kubectl", "-n", nameSpace, "get", objectType, "--output=jsonpath={.items..metadata.name}")
			if err != nil {
				return err
			}

			objects = strings.Split(output, " ")
		}

		for _, object := range objects {
			kubernetesObjectFile := filepath.Join(rootPath, nameSpace+"_"+objectType+"_"+object)

			output, err := utils.RunCommandOnContainer("kubectl", "-n", nameSpace, "describe", objectType, object)
			if err != nil {
				return err
			}

			err = utils.WriteToFile(kubernetesObjectFile, output)
			if err != nil {
				return err
			}

			collector.AddToCollectorFiles(kubernetesObjectFile)
		}
	}

	return nil
}
