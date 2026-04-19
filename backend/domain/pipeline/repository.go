package pipeline

type PipelineRepo interface {
	Save(pipeline *Pipeline) error
	FindByID(id string) (*Pipeline, error)
	FindByCardID(cardID string) (*Pipeline, error)
	Update(pipeline *Pipeline) error
}
