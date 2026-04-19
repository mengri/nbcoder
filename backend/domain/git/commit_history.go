package git

import (
	"fmt"
	"time"
)

type Commit struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

func NewCommit(hash, message, author string) *Commit {
	return &Commit{
		Hash:      hash,
		Message:   message,
		Author:    author,
		Timestamp: time.Now().UTC(),
	}
}

type CommitHistory struct {
	Commits []*Commit `json:"commits"`
}

func NewCommitHistory() *CommitHistory {
	return &CommitHistory{
		Commits: []*Commit{},
	}
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

func (ch *CommitHistory) FindByHash(hash string) *Commit {
	for _, c := range ch.Commits {
		if c.Hash == hash {
			return c
		}
	}
	return nil
}

func (ch *CommitHistory) FindByAuthor(author string) []*Commit {
	var result []*Commit
	for _, c := range ch.Commits {
		if c.Author == author {
			result = append(result, c)
		}
	}
	return result
}

func (ch *CommitHistory) FindByTimeRange(start, end time.Time) []*Commit {
	var result []*Commit
	for _, c := range ch.Commits {
		if !c.Timestamp.Before(start) && !c.Timestamp.After(end) {
			result = append(result, c)
		}
	}
	return result
}

type IncrementalCommitBuilder struct {
	history    *CommitHistory
	pendingMsgs []string
	author     string
}

func NewIncrementalCommitBuilder(history *CommitHistory, author string) *IncrementalCommitBuilder {
	return &IncrementalCommitBuilder{
		history: history,
		author:  author,
	}
}

func (b *IncrementalCommitBuilder) AddChange(message string) {
	b.pendingMsgs = append(b.pendingMsgs, message)
}

func (b *IncrementalCommitBuilder) Commit() (*Commit, error) {
	if len(b.pendingMsgs) == 0 {
		return nil, fmt.Errorf("no changes to commit")
	}
	hash := fmt.Sprintf("%x", len(b.history.Commits)+1)
	msg := b.pendingMsgs[0]
	if len(b.pendingMsgs) > 1 {
		msg = fmt.Sprintf("%s (+%d more)", b.pendingMsgs[0], len(b.pendingMsgs)-1)
	}
	commit := NewCommit(hash, msg, b.author)
	b.history.Add(commit)
	b.pendingMsgs = nil
	return commit, nil
}

func (b *IncrementalCommitBuilder) HasPendingChanges() bool {
	return len(b.pendingMsgs) > 0
}

func (b *IncrementalCommitBuilder) SquashCommit(message string) (*Commit, error) {
	if len(b.pendingMsgs) == 0 {
		return nil, fmt.Errorf("no changes to squash")
	}
	hash := fmt.Sprintf("%x", len(b.history.Commits)+1)
	commit := NewCommit(hash, message, b.author)
	b.history.Add(commit)
	b.pendingMsgs = nil
	return commit, nil
}
