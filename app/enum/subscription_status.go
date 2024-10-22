package enum

type SubscriptionStatus string

const (
	SUBSCRIPTION_STATUS_ACTIVE   SubscriptionStatus = "active"
	SUBSCRIPTION_STATUS_INACTIVE SubscriptionStatus = "inactive"
)

func GetSubscriptionStatus(t SubscriptionStatus) SubscriptionStatus {
	switch t {
	case SUBSCRIPTION_STATUS_ACTIVE:
		return SUBSCRIPTION_STATUS_ACTIVE
	default:
		return SUBSCRIPTION_STATUS_INACTIVE
	}
}
