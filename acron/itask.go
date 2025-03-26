package acron

// ITask is an interface for job data that can return its type.
type ITask interface {
	GetType() TaskType
	Validate() error
	Run(ICronControlCenter) error
}
