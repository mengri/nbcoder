package agent

type AgentType string

const (
	AgentTypeProduct      AgentType = "PRODUCT"
	AgentTypeArchitecture AgentType = "ARCHITECTURE"
	AgentTypeManagement   AgentType = "MANAGEMENT"
	AgentTypeTechStack    AgentType = "TECH_STACK"
)

func (t AgentType) IsValid() bool {
	switch t {
	case AgentTypeProduct, AgentTypeArchitecture, AgentTypeManagement, AgentTypeTechStack:
		return true
	}
	return false
}
