package git
// branch_policy.go
// 分支策略与命名规范配置与校验


import (
	"regexp"
)

type BranchPolicy struct {
	ProjectID      string
	AllowedPattern string // e.g. ^(main|dev|feature/\w+|fix/\w+|release/\w+)$
}

func (p *BranchPolicy) Validate(branch string) bool {
	re := regexp.MustCompile(p.AllowedPattern)
	return re.MatchString(branch)
}
