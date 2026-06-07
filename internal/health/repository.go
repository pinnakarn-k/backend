package health

type Repository interface {
	Check() error
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Check() error {
	return nil
}
