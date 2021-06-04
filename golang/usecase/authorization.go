package usecase

import "app/domain/function"

type AuthorizationUsecase struct {
	DBOperator *function.DatabaseOperater
}
