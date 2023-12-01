package request

import "github.com/go-playground/validator/v10"

// @Description The request payload for checking application updates.
type CheckRequest struct {
	// PackageName is the unique identifier of the application package.
	// @Required
	PackageName string `json:"packageName" validate:"required" example:"com.wayofdt.kidneysmart"`

	// BuildNumber is the current version build of the application.
	// Must be greater than 0.
	// @Required
	VersionCode int `json:"versionCode" validate:"gt=0" example:"2"`

	// BuildNumber is the current version build of the application.
	// Must be greater than 0.
	// @Required
	VersionName string `json:"versionName" validate:"required" example:"2.0.0"`


	// InstallerPackageName is the package name of the installer.
	// @Required
	InstallerPackageName string `json:"installerPackageName" validate:"required" example:"apk"`
}

// Validate использует валидатор для проверки полей структуры AppUpdateRequest.
func (a *CheckRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
