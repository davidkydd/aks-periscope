module github.com/Azure/aks-periscope

// 1.16 required for go:embed (used for testing resources)
go 1.16

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2 v2.1.0 // indirect
	github.com/Azure/azure-storage-blob-go v0.14.0
	github.com/Azure/go-autorest/autorest/adal v0.9.14 // indirect
	github.com/containerd/containerd v1.4.13 // indirect
	github.com/docker/docker v20.10.14+incompatible
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	helm.sh/helm/v3 v3.6.3
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/cli-runtime v0.21.3
	k8s.io/client-go v0.21.3
	k8s.io/kubectl v0.21.0
	k8s.io/metrics v0.21.0
	rsc.io/letsencrypt v0.0.3 // indirect
)
