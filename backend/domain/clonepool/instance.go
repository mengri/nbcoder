package clonepool
// instance.go
// 克隆实例状态流转与校验

import "time"

type CloneInstanceStatus string

const (
	InstanceIdle      CloneInstanceStatus = "IDLE"
	InstanceAllocating CloneInstanceStatus = "ALLOCATING"
	InstanceInUse     CloneInstanceStatus = "IN_USE"
	InstanceRecycling CloneInstanceStatus = "RECYCLING"
	InstanceReleased  CloneInstanceStatus = "RELEASED"
)

type CloneInstance struct {
	ID        string              `json:"id"`
	Status    CloneInstanceStatus `json:"status"`
	UpdatedAt time.Time           `json:"updated_at"`
	Logs      []string            `json:"logs"`
}

func (ci *CloneInstance) Allocate() {
	ci.Status = InstanceAllocating
	ci.UpdatedAt = time.Now().UTC()
	ci.Logs = append(ci.Logs, "Instance allocating at "+ci.UpdatedAt.String())
}

func (ci *CloneInstance) Use() {
	ci.Status = InstanceInUse
	ci.UpdatedAt = time.Now().UTC()
	ci.Logs = append(ci.Logs, "Instance in use at "+ci.UpdatedAt.String())
}

func (ci *CloneInstance) Recycle() {
	ci.Status = InstanceRecycling
	ci.UpdatedAt = time.Now().UTC()
	ci.Logs = append(ci.Logs, "Instance recycling at "+ci.UpdatedAt.String())
}

func (ci *CloneInstance) Release() {
	ci.Status = InstanceReleased
	ci.UpdatedAt = time.Now().UTC()
	ci.Logs = append(ci.Logs, "Instance released at "+ci.UpdatedAt.String())
}

func (ci *CloneInstance) IsValidTransition(newStatus CloneInstanceStatus) bool {
	switch ci.Status {
	case InstanceIdle:
		return newStatus == InstanceAllocating
	case InstanceAllocating:
		return newStatus == InstanceInUse
	case InstanceInUse:
		return newStatus == InstanceRecycling
	case InstanceRecycling:
		return newStatus == InstanceReleased || newStatus == InstanceIdle
	case InstanceReleased:
		return newStatus == InstanceIdle
	default:
		return false
	}
}
