// model_chain.go
// 模型链配置与路由
package airuntime

type ModelRouteStrategy string

const (
	RouteByType   ModelRouteStrategy = "BY_TYPE"
	RouteByWeight ModelRouteStrategy = "BY_WEIGHT"
	RouteCustom   ModelRouteStrategy = "CUSTOM"
)

type ModelChain struct {
	ID     string       `json:"id"`
	Models []string     `json:"models"`
	Routes []ModelRoute `json:"routes"`
}

type ModelRoute struct {
	From     string             `json:"from"`
	To       string             `json:"to"`
	Strategy ModelRouteStrategy `json:"strategy"`
	Weight   int                `json:"weight,omitempty"`
}

// 路由选择示例
func (mc *ModelChain) SelectRoute(from string) *ModelRoute {
	for _, r := range mc.Routes {
		if r.From == from {
			return &r
		}
	}
	return nil
}
