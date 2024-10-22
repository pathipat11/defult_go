package enum

type SprintStatus string

const (
	SPRINT_STATUS_IN_PROGRESS SprintStatus = "in_progress"
	SPRINT_STATUS_COMPLETED   SprintStatus = "completed"
	SPRINT_STATUS_BACKLOG     SprintStatus = "backlog"
)

func GetSprintStatus(t SprintStatus) SprintStatus {
	switch t {
	case SPRINT_STATUS_COMPLETED:
		return SPRINT_STATUS_COMPLETED
	case SPRINT_STATUS_BACKLOG:
		return SPRINT_STATUS_BACKLOG
	default:
		return SPRINT_STATUS_IN_PROGRESS
	}
}
