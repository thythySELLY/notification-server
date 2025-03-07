package models

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

func IsValidStatus(status string) bool {
	switch status {
	case StatusActive, StatusInactive:
		return true
	}
	return false
}
