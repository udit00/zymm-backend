package bussinessAuth

import (
	"errors"
	"strings"
	"zymm/internal/models"
)

func ValidateRegistrationRequest(req models.RegistrationApiRequestModel) error {
	missing := []string{}

	if req.DisplayName == "" {
		missing = append(missing, "displayName")
	}
	if req.Mobile == "" {
		missing = append(missing, "mobile")
	}
	if req.Password == "" {
		missing = append(missing, "password")
	}
	if req.Gender == "" {
		missing = append(missing, "gender")
	}
	if req.AppVersion == "" {
		missing = append(missing, "appVersion")
	}
	if req.UserAgent == "" {
		missing = append(missing, "userAgent")
	}
	if req.LocationLat == "" {
		missing = append(missing, "locationLat")
	}
	if req.LocationLong == "" {
		missing = append(missing, "locationLong")
	}
	if req.IpAddress == "" {
		missing = append(missing, "ipAddress")
	}

	if len(missing) > 0 {
		return errors.New("Missing Required Fields: " + strings.Join(missing, ", "))
	} else {
		return nil
	}
}

func ValidateOwnerRegistrationRequest(req models.RegistrationOwnerApiRequestModel) error {
	missing := []string{}

	if req.DisplayName == "" {
		missing = append(missing, "displayName")
	}
	if req.Mobile == "" {
		missing = append(missing, "mobile")
	}
	if req.Password == "" {
		missing = append(missing, "password")
	}
	if req.Gender == "" {
		missing = append(missing, "gender")
	}
	if req.AppVersion == "" {
		missing = append(missing, "appVersion")
	}
	if req.UserAgent == "" {
		missing = append(missing, "userAgent")
	}
	if req.IpAddress == "" {
		missing = append(missing, "ipAddress")
	}
	if req.GymName == "" {
		missing = append(missing, "gymName")
	}
	if req.State == "" {
		missing = append(missing, "state")
	}
	if req.City == "" {
		missing = append(missing, "city")
	}
	if req.GymAddress == "" {
		missing = append(missing, "gymAddress")
	}
	if req.ContactNo == "" {
		missing = append(missing, "gymOfficialContactNo")
	}
	if req.OfficialEmail == "" {
		missing = append(missing, "gymOfficialEmail")
	}
	if req.LocationLat == "" {
		missing = append(missing, "gymOfficialLocationLat")
	}
	if req.LocationLong == "" {
		missing = append(missing, "gymOfficialLocationLong")
	}

	if len(missing) > 0 {
		return errors.New("Missing Required Fields: " + strings.Join(missing, ", "))
	} else {
		return nil
	}
}
