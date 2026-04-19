package airuntime

import "time"

type ProviderRepo interface {
	Save(provider *Provider) error
	FindByID(id string) (*Provider, error)
	FindAll() ([]*Provider, error)
	Update(provider *Provider) error
}

type ChainRepo interface {
	Save(chain *Chain) error
	FindByID(id string) (*Chain, error)
	FindAll() ([]*Chain, error)
}

type CallLogRepo interface {
	Save(log *CallLog) error
	FindByAgentID(agentID string) ([]*CallLog, error)
	FindByTimeRange(start, end time.Time) ([]*CallLog, error)
}
