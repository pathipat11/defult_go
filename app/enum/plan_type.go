package enum

type PlanType string

const (
	PLAN_TYPE_BASIC   PlanType = "basic"
	PLAN_TYPE_PREMIUM PlanType = "premium"
)

func GetPlanType(t PlanType) PlanType {
	switch t {
	case PLAN_TYPE_BASIC:
		return PLAN_TYPE_BASIC
	default:
		return PLAN_TYPE_PREMIUM
	}
}
