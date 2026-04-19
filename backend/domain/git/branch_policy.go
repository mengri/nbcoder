package git

import "regexp"

type BranchPolicy struct {
	ProjectID      string
	AllowedPattern string
}

func (p *BranchPolicy) Validate(branch string) bool {
	re := regexp.MustCompile(p.AllowedPattern)
	return re.MatchString(branch)
}

type Branch struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

func NewBranch(name string, isDefault bool) *Branch {
	return &Branch{
		Name:      name,
		IsDefault: isDefault,
	}
}
