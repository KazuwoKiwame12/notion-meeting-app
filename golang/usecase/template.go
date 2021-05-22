package usecase

import "errors"

type TemplateUsecase struct {
}

func (t *TemplateUsecase) Create() error {
	return errors.New("error")
}
