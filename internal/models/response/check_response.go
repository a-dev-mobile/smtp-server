package response

// @Description holds information about the latest version of the app.
type VersionInfo struct {
	// VersionCode represents the version build of the application.
	VersionCode int `json:"versionCode,omitempty"`

	// Name is the version of the latest application version.
	VersionName string `json:"versionName,omitempty"`

	// Url is the URL where the latest version can be downloaded from.
	Url string `json:"url,omitempty"`
}

// @Description holds the overall update information for an application.
type CheckResponse struct {
	// UpdateType indicates the type of update required ('hard', 'soft', or 'none').
	// When UpdateType is 'none', the LatestVersion field will not be present in the response.
	UpdateType string `json:"updateType,omitempty"`

	// LatestVersion contains the details of the latest version if an update is available.
	// This field is omitted when UpdateType is 'none'.
	LatestVersion *VersionInfo `json:"latestVersion,omitempty"`
}
