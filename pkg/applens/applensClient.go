package applens

// ApplensSignal defines an Applens Signal struct
type ApplensClient struct {
}

// NewApplensClient is a constructor
func NewApplensClient() *ApplensSignal {
	return &ApplensSignal{}
}

/* sample pattern to call azurecontainerservice API from azure-sdk-go, adapt to call getDetectors url

// GetUpgradeProfile - Gets the upgrade profile of a managed cluster.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-07-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// resourceName - The name of the managed cluster resource.
// options - ManagedClustersClientGetUpgradeProfileOptions contains the optional parameters for the ManagedClustersClient.GetUpgradeProfile
// method.
func (client *ManagedClustersClient) GetUpgradeProfile(ctx context.Context, resourceGroupName string, resourceName string, options *ManagedClustersClientGetUpgradeProfileOptions) (ManagedClustersClientGetUpgradeProfileResponse, error) {
	req, err := client.getUpgradeProfileCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return ManagedClustersClientGetUpgradeProfileResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ManagedClustersClientGetUpgradeProfileResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ManagedClustersClientGetUpgradeProfileResponse{}, runtime.NewResponseError(resp)
	}
	return client.getUpgradeProfileHandleResponse(resp)
}


// getUpgradeProfileCreateRequest creates the GetUpgradeProfile request.
func (client *ManagedClustersClient) getUpgradeProfileCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *ManagedClustersClientGetUpgradeProfileOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/upgradeProfiles/default"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getUpgradeProfileHandleResponse handles the GetUpgradeProfile response.
func (client *ManagedClustersClient) getUpgradeProfileHandleResponse(resp *http.Response) (ManagedClustersClientGetUpgradeProfileResponse, error) {
	result := ManagedClustersClientGetUpgradeProfileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedClusterUpgradeProfile); err != nil {
		return ManagedClustersClientGetUpgradeProfileResponse{}, err
	}
	return result, nil
}
*/
