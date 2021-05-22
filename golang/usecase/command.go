package usecase

import "errors"

type CommandUsecase struct {
}

func (c *CommandUsecase) Start() error {
	return errors.New("error")
}

func (c *CommandUsecase) Cancel() error {
	return errors.New("error")
}
