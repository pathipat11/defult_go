package enum

type AdDisplayStatus string

const (
	AD_DISPLAY_STATUS_ACTIVE   AdDisplayStatus = "active"
	AD_DISPLAY_STATUS_INACTIVE AdDisplayStatus = "inactive"
)

func GetAdDisplayStatus(t AdDisplayStatus) AdDisplayStatus {
	switch t {
	case AD_DISPLAY_STATUS_ACTIVE:
		return AD_DISPLAY_STATUS_ACTIVE
	default:
		return AD_DISPLAY_STATUS_INACTIVE
	}
}
