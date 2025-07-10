package usertypes

import (
	"fmt"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// User interface
type User interface {
	GetID() int64
	GetUserID() sharedtypes.DiscordID
	GetRole() sharedtypes.UserRoleEnum
}

// UserData struct implementing the User interface
type UserData struct {
	ID     int64                    `json:"id"`
	UserID sharedtypes.DiscordID    `json:"user_id"`
	Role   sharedtypes.UserRoleEnum `json:"role"`
}

func (u UserData) GetID() int64 {
	return u.ID
}

func (u UserData) GetUserID() sharedtypes.DiscordID {
	return u.UserID
}

func (u UserData) GetRole() sharedtypes.UserRoleEnum {
	return u.Role
}

// ParseUserRoleEnum converts a string to a UserRoleEnum.
func ParseUserRoleEnum(role string) (sharedtypes.UserRoleEnum, error) {
	switch role {
	case "User":
		return sharedtypes.UserRoleUser, nil
	case "Editor":
		return sharedtypes.UserRoleEditor, nil
	case "Admin":
		return sharedtypes.UserRoleAdmin, nil
	default:
		return sharedtypes.UserRoleUnknown, fmt.Errorf("invalid user role: %s", role)
	}
}
