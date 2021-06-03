package collector

import (
	"os"
	"reflect"
	"testing"

	"github.com/Azure/aks-periscope/pkg/interfaces"
)

func TestKubeObjectsCollector_Collect(t *testing.T) {
	type fields struct {
		BaseCollector BaseCollector
		config        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"test", fields{config: `kubeobjects_to_collect:
- id: all-resources-table
  kubeobject: all
  selector: app.kubernetes.io/name=openservicemesh.io
  output: wide
  namespace: all`,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &KubeObjectsCollector{
				BaseCollector: tt.fields.BaseCollector,
			}
			os.Setenv("kubeobjects_to_collect", tt.fields.config)

			if err := collector.Collect(); (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewKubeObjectsCollector(t *testing.T) {
	type args struct {
		exporters []interfaces.Exporter
	}
	tests := []struct {
		name string
		args args
		want *KubeObjectsCollector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKubeObjectsCollector(tt.args.exporters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKubeObjectsCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}
