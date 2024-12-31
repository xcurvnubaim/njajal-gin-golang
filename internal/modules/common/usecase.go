package common

type IUseCase interface {
}

type useCase struct {
	repository IRepository
}

func NewuseCase(repository IRepository) *useCase {
	return &useCase{repository}
}

// func (u *useCase) GetFile(path string)  