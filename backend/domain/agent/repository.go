package agent

type TaskRepo interface {
	Save(task *Task) error
	FindByID(id string) (*Task, error)
	FindByStatus(status TaskStatus) ([]*Task, error)
	Update(task *Task) error
}
