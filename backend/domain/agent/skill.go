package agent

type Skill struct {
	Name        string
	Description string
	Invoke      func(args ...interface{}) (interface{}, error)
}

type SkillError struct {
	Msg string
}

func (e *SkillError) Error() string {
	return e.Msg
}
