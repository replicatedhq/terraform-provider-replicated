package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type CustomersAPIClient struct {
	Client  HttpRequestDoer
	ApiKey  string
	BaseURL string
}

type EntitlementValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GetCustomerResponse struct {
	Customer struct {
		ID                               string    `json:"id"`
		TeamID                           string    `json:"teamId"`
		Name                             string    `json:"name"`
		Email                            string    `json:"email"`
		CreatedAt                        time.Time `json:"createdAt"`
		UpdatedAt                        time.Time `json:"updatedAt"`
		ExpiresAt                        string    `json:"expiresAt"`
		IsArchived                       bool      `json:"isArchived"`
		Type                             string    `json:"type"`
		InstallationID                   string    `json:"installationId"`
		InstallationVersion              string    `json:"installationVersion"`
		AppType                          string    `json:"appType"`
		Airgap                           bool      `json:"airgap"`
		IsGitopsSupported                bool      `json:"isGitopsSupported"`
		IsIdentityServiceSupported       bool      `json:"isIdentityServiceSupported"`
		IsGeoaxisSupported               bool      `json:"isGeoaxisSupported"`
		IsSnapshotSupported              bool      `json:"isSnapshotSupported"`
		IsSupportBundleUploadEnabled     bool      `json:"isSupportBundleUploadEnabled"`
		IsHelmVMDownloadEnabled          bool      `json:"isHelmVmDownloadEnabled"`
		IsEmbeddedClusterDownloadEnabled bool      `json:"isEmbeddedClusterDownloadEnabled"`
		IsKotsInstallEnabled             bool      `json:"isKotsInstallEnabled"`
		IsInstallerSupportEnabled        bool      `json:"isInstallerSupportEnabled"`
		InstalledReleaseSequence         int       `json:"installedReleaseSequence"`
		InstalledReleaseLabel            string    `json:"installedReleaseLabel"`
		LastActive                       time.Time `json:"lastActive"`
		Actions                          struct {
			ShipApplyDocker  string `json:"shipApplyDocker"`
			ShipApply        string `json:"shipApply"`
			ShipInitHomebrew string `json:"shipInitHomebrew"`
			ShipInitCloud    string `json:"shipInitCloud"`
		} `json:"actions"`
		ShipInstallStatus struct {
			Status    string    `json:"status"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"shipInstallStatus"`
		Channels []struct {
			ID                                string    `json:"id"`
			AppID                             string    `json:"appId"`
			AppSlug                           string    `json:"appSlug"`
			AppName                           string    `json:"appName"`
			ChannelSlug                       string    `json:"channelSlug"`
			Name                              string    `json:"name"`
			Description                       string    `json:"description"`
			ChannelIcon                       string    `json:"channelIcon"`
			IsArchived                        bool      `json:"isArchived"`
			IsDefault                         bool      `json:"isDefault"`
			Created                           time.Time `json:"created"`
			Updated                           time.Time `json:"updated"`
			NumReleases                       int       `json:"numReleases"`
			IsHelmOnly                        bool      `json:"isHelmOnly"`
			BuildAirgapAutomatically          bool      `json:"buildAirgapAutomatically"`
			AirgapDockerRegistryFormatEnabled bool      `json:"airgapDockerRegistryFormatEnabled"`
			SemverRequired                    bool      `json:"semverRequired"`
			SemverWarning                     string    `json:"semverWarning"`
			TargetKotsVersion                 string    `json:"targetKotsVersion"`
			EnterprisePartnerChannelID        string    `json:"enterprisePartnerChannelID"`
			ChannelSequence                   int       `json:"channelSequence"`
			CurrentVersion                    string    `json:"currentVersion"`
			ExtraLintRules                    any       `json:"extraLintRules"`
			GitHubRef                         any       `json:"gitHubRef"`
			ReleaseNotes                      string    `json:"releaseNotes"`
			ReleaseSequence                   int       `json:"releaseSequence"`
			ReplicatedRegistryDomain          string    `json:"replicatedRegistryDomain"`
			CustomHostnameOverrides           struct {
				Registry struct {
					Hostname string `json:"hostname"`
				} `json:"registry"`
				Proxy struct {
					Hostname string `json:"hostname"`
				} `json:"proxy"`
				DownloadPortal struct {
					Hostname string `json:"hostname"`
				} `json:"downloadPortal"`
				ReplicatedApp struct {
					Hostname string `json:"hostname"`
				} `json:"replicatedApp"`
			} `json:"customHostnameOverrides"`
			AdoptionRate []any `json:"adoptionRate"`
			Customers    struct {
				TotalCustomers    int `json:"totalCustomers"`
				ActiveCustomers   int `json:"activeCustomers"`
				InactiveCustomers int `json:"inactiveCustomers"`
			} `json:"customers"`
			Releases      []any `json:"releases"`
			ChartReleases any   `json:"chartReleases"`
		} `json:"channels"`
		Instances []struct {
			LicenseID      string    `json:"licenseId"`
			InstanceID     string    `json:"instanceId"`
			ClusterID      string    `json:"clusterId"`
			CreatedAt      time.Time `json:"createdAt"`
			LastActive     time.Time `json:"lastActive"`
			AppStatus      string    `json:"appStatus"`
			Active         bool      `json:"active"`
			VersionHistory []struct {
				InstanceID                string    `json:"instanceId"`
				ClusterID                 string    `json:"clusterId"`
				VersionLabel              string    `json:"versionLabel"`
				DownstreamChannelID       string    `json:"downstreamChannelId"`
				DownstreamReleaseSequence int       `json:"downstreamReleaseSequence"`
				IntervalStart             time.Time `json:"intervalStart"`
				IntervalLast              time.Time `json:"intervalLast"`
				ReplHelmCount             int       `json:"replHelmCount"`
				NativeHelmCount           int       `json:"nativeHelmCount"`
			} `json:"versionHistory"`
			KotsVersion          string `json:"kotsVersion"`
			ReplicatedSdkVersion string `json:"replicatedSdkVersion"`
			Cloud                string `json:"cloud"`
			IsAirgap             bool   `json:"isAirgap"`
			IsKurl               bool   `json:"isKurl"`
			KurlNodeCountTotal   int    `json:"kurlNodeCountTotal"`
			KurlNodeCountReady   int    `json:"kurlNodeCountReady"`
			K8SVersion           string `json:"k8sVersion"`
			Client               string `json:"client"`
			IsDummyInstance      bool   `json:"isDummyInstance"`
			Tags                 []any  `json:"tags"`
		} `json:"instances"`
		IsInstancesLimited bool               `json:"isInstancesLimited"`
		Entitlements       []EntitlementValue `json:"entitlements"`
		DownloadPortalURL  string             `json:"downloadPortalUrl"`
		CreatedBy          struct {
			ID          string    `json:"id"`
			Type        string    `json:"type"`
			Description string    `json:"description"`
			Link        string    `json:"link"`
			Timestamp   time.Time `json:"timestamp"`
		} `json:"createdBy"`
		UpdatedBy struct {
			ID          string    `json:"id"`
			Type        string    `json:"type"`
			Description string    `json:"description"`
			Link        string    `json:"link"`
			Timestamp   time.Time `json:"timestamp"`
		} `json:"updatedBy"`
	} `json:"customer"`
}

type UpdateCustomerOpts struct {
	AppID                            string             `json:"app_id"`
	ChannelId                        string             `json:"channel_id"`
	Email                            string             `json:"email"`
	EntitlementValues                []EntitlementValue `json:"entitlementValues"`
	ExpiresAt                        string             `json:"expires_at"`
	IsAirgapEnabled                  bool               `json:"is_airgap_enabled"`
	IsEmbeddedClusterDownloadEnabled bool               `json:"is_embedded_cluster_download_enabled"`
	IsGeoaxisSupported               bool               `json:"is_geoaxis_supported"`
	IsGitopsSupported                bool               `json:"is_gitops_supported"`
	IsHelmVMDownloadEnabled          bool               `json:"is_helm_vm_download_enabled"`
	IsIdentityServiceSupported       bool               `json:"is_identity_service_supported"`
	IsKotsInstallEnabled             bool               `json:"is_kots_install_enabled"`
	IsSnapshotSupported              bool               `json:"is_snapshot_supported"`
	IsSupportBundleUploadEnabled     bool               `json:"is_support_bundle_upload_enabled"`
	Name                             string             `json:"name"`
	Type                             string             `json:"type"`
}

type CreateCustomerOpts struct {
	AppID                            string             `json:"app_id"`
	ChannelId                        string             `json:"channel_id"`
	Email                            string             `json:"email"`
	EntitlementValues                []EntitlementValue `json:"entitlementValues"`
	ExpiresAt                        string             `json:"expires_at"`
	IsAirgapEnabled                  bool               `json:"is_airgap_enabled"`
	IsEmbeddedClusterDownloadEnabled bool               `json:"is_embedded_cluster_download_enabled"`
	IsGeoaxisSupported               bool               `json:"is_geoaxis_supported"`
	IsGitopsSupported                bool               `json:"is_gitops_supported"`
	IsHelmVMDownloadEnabled          bool               `json:"is_helm_vm_download_enabled"`
	IsIdentityServiceSupported       bool               `json:"is_identity_service_supported"`
	IsInstallerSupportEnabled        bool               `json:"is_installer_support_enabled"`
	IsKotsInstallEnabled             bool               `json:"is_kots_install_enabled"`
	IsSnapshotSupported              bool               `json:"is_snapshot_supported"`
	IsSupportBundleUploadEnabled     bool               `json:"is_support_bundle_upload_enabled"`
	Name                             string             `json:"name"`
	Type                             string             `json:"type"`
}

func NewCustomersAPIClient(apiKey string) *CustomersAPIClient {
	return &CustomersAPIClient{
		Client:  http.DefaultClient,
		ApiKey:  apiKey,
		BaseURL: "https://api.replicated.com/vendor/v3",
	}
}

func (c *CustomersAPIClient) CreateCustomer(opts CreateCustomerOpts) (*GetCustomerResponse, error) {
	url := fmt.Sprintf("%s/customer", c.BaseURL)

	// Serialize customer to JSON
	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.ApiKey)

	// Send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Check response status
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create customer failed with status code: %d, %s", res.StatusCode, body)
	}

	// Unmarshal response body
	var result GetCustomerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *CustomersAPIClient) GetCustomer(appId string, customerId string) (*GetCustomerResponse, error) {
	url := fmt.Sprintf("%s/app/%s/customer/%s", c.BaseURL, appId, customerId)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", c.ApiKey)

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Check if the response status code is OK (200)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get customer. status code: %d, %s", res.StatusCode, url)
	}

	var result GetCustomerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	// return nil, fmt.Errorf("failed to unmarshal response body: %s", result)
	return &result, nil
}

func (c *CustomersAPIClient) UpdateCustomer(customerId string, opts UpdateCustomerOpts) error {
	url := fmt.Sprintf("%s/customer/%s", c.BaseURL, customerId)

	// Serialize customer to JSON
	payload, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.ApiKey)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("update customer failed with status code: %d", err)
		}
		return fmt.Errorf("update customer failed with status code: %d, %s", resp.StatusCode, responseData)
	}

	return nil
}
