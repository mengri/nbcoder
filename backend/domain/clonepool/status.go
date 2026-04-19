package clonepool

type CloneInstanceStatus string

const (
	InstanceIdle   CloneInstanceStatus = "IDLE"
	InstanceBusy   CloneInstanceStatus = "BUSY"
	InstanceDirty  CloneInstanceStatus = "DIRTY"
)

func (s CloneInstanceStatus) IsValid() bool {
	switch s {
	case InstanceIdle, InstanceBusy, InstanceDirty:
		return true
	}
	return false
}
