package usecase

import (
	"app/domain/function"
	"log"
)

const (
	failedInt int    = -100
	failedStr string = ""
)

type AuthorizationUsecase struct {
	DBOperator *function.DatabaseOperater
}

func (au *AuthorizationUsecase) IsUser(workspaceID, slackUserID string) (int, string, error) {
	user, err := au.DBOperator.GetUser(workspaceID, slackUserID)
	if err != nil {
		log.Printf("Failed to get user: %+v", err)
		return failedInt, failedStr, err
	}
	return user.ID, user.Name, nil
}

func (au *AuthorizationUsecase) IsAdmin(workspaceID, slackUserID string) (int, string, error) {
	user, err := au.DBOperator.GetUser(workspaceID, slackUserID)
	if err != nil {
		log.Printf("Failed to get user: %+v", err)
		return failedInt, failedStr, err
	}
	if !user.IsAdministrator {
		log.Printf("Failed to get admin user: %+v", err)
		return failedInt, failedStr, err
	}
	return user.ID, user.Name, nil
}
