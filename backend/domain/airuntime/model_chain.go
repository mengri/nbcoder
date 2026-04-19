package airuntime

type ModelRouteStrategy string

const (
	RouteByType   ModelRouteStrategy = "BY_TYPE"
	RouteByWeight ModelRouteStrategy = "BY_WEIGHT"
	RouteCustom   ModelRouteStrategy = "CUSTOM"
)

type Chain struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Models []string     `json:"models"`
	Routes []ModelRoute `json:"routes"`
}

type ModelRoute struct {
	From     string             `json:"from"`
	To       string             `json:"to"`
	Strategy ModelRouteStrategy `json:"strategy"`
	Weight   int                `json:"weight,omitempty"`
}

func (c *Chain) SelectRoute(from string) *ModelRoute {
	for _, r := range c.Routes {
		if r.From == from {
			return &r
		}
	}
	return nil
}

func (c *Chain) NextModel(currentModelID string) string {
	route := c.SelectRoute(currentModelID)
	if route != nil {
		return route.To
	}
	for i, m := range c.Models {
		if m == currentModelID && i+1 < len(c.Models) {
			return c.Models[i+1]
		}
	}
	return ""
}
