package enum

type TaskStatus string

const (
	TASK_STATUS_TODO        TaskStatus = "todo"
	TASK_STATUS_IN_PROGRESS TaskStatus = "in_progress"
	TASK_STATUS_DONE        TaskStatus = "done"
)

func GetTaskStatus(t TaskStatus) TaskStatus {
	switch t {
	case TASK_STATUS_IN_PROGRESS:
		return TASK_STATUS_IN_PROGRESS
	case TASK_STATUS_DONE:
		return TASK_STATUS_DONE
	default:
		return TASK_STATUS_TODO
	}
}
