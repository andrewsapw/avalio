package status

type Status struct {
	Success      bool
	ResourceName string
}

func NewStatus(success bool, resourceName string) Status {
	return Status{success, resourceName}
}
