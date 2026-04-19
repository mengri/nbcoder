package git

import "time"

type Commit struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

type CommitHistory struct {
	Commits []*Commit `json:"commits"`
}

func (ch *CommitHistory) Add(commit *Commit) {
	ch.Commits = append(ch.Commits, commit)
}

func (ch *CommitHistory) Latest() *Commit {
	if len(ch.Commits) == 0 {
		return nil
	}
	return ch.Commits[len(ch.Commits)-1]
}
