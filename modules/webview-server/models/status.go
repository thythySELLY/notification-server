package models

const (
	StatusActive  = "active"
	StatusInactive= "inactive"
)

// Function to validate the status
func IsValidStatus(status string) bool {
	switch (status) {
	case StatusActive, StatusInactive:
		return true
	}
	return false
}
