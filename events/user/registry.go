package userevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the user functional area
func GetV1Registry() map[string]sharedevents.EventInfo {
	return map[string]sharedevents.EventInfo{
		UserCreationRequestedV1: {
			Payload:     &UserCreationRequestedPayloadV1{},
			Summary:     "User Creation Requested",
			Description: "Request to create a new user.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}},
		},
		UserCreatedV1: {
			Payload:     &UserCreatedPayloadV1{},
			Summary:     "User Created",
			Description: "User successfully created.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},
		UserCreationFailedV1: {
			Payload:     &UserCreationFailedPayloadV1{},
			Summary:     "User Creation Failed",
			Description: "User creation failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},

		UserSignupRequestedV1: {
			Payload:     &UserSignupRequestedPayloadV1{},
			Summary:     "User Signup Requested",
			Description: "User initiated signup.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}},
		},
		UserSignupSucceededV1: {
			Payload:     &UserSignupSucceededPayloadV1{},
			Summary:     "User Signup Succeeded",
			Description: "User signup succeeded.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},
		UserSignupFailedV1: {
			Payload:     &UserSignupFailedPayloadV1{},
			Summary:     "User Signup Failed",
			Description: "User signup failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},

		// Retrieval flow
		GetUserRequestedV1: {
			Payload:     &GetUserRequestedPayloadV1{},
			Summary:     "Get User Requested",
			Description: "Request to retrieve a user's data.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}},
		},
		GetUserResponseV1: {
			Payload:     &GetUserResponsePayloadV1{},
			Summary:     "Get User Response",
			Description: "Response containing user data.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},
		GetUserFailedV1: {
			Payload:     &GetUserFailedPayloadV1{},
			Summary:     "Get User Failed",
			Description: "User retrieval failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},

		// Role flow
		UserRoleUpdateRequestedV1: {
			Payload:     &UserRoleUpdateRequestedPayloadV1{},
			Summary:     "User Role Update Requested",
			Description: "Request to update a user's role.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}},
		},
		UserRoleUpdatedV1: {
			Payload:     &UserRoleUpdatedPayloadV1{},
			Summary:     "User Role Updated",
			Description: "A user's role was updated.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},
		UserRoleUpdateFailedV1: {
			Payload:     &UserRoleUpdateFailedPayloadV1{},
			Summary:     "User Role Update Failed",
			Description: "User role update failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},

		GetUserRoleRequestedV1: {
			Payload:     &GetUserRoleRequestedPayloadV1{},
			Summary:     "Get User Role Requested",
			Description: "Request to retrieve a user's role.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}},
		},
		GetUserRoleResponseV1: {
			Payload:     &GetUserRoleResponsePayloadV1{},
			Summary:     "Get User Role Response",
			Description: "Response containing the user's role.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},
		GetUserRoleFailedV1: {
			Payload:     &GetUserRoleFailedPayloadV1{},
			Summary:     "Get User Role Failed",
			Description: "User role retrieval failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},

		UserPermissionsCheckRequestedV1: {
			Payload:     &UserPermissionsCheckRequestedPayloadV1{},
			Summary:     "User Permissions Check Requested",
			Description: "Request to check user permissions.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "round"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}},
		},
		UserPermissionsCheckResponseV1: {
			Payload:     &UserPermissionsCheckResponsePayloadV1{},
			Summary:     "User Permissions Check Response",
			Description: "Response containing permissions check result.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "round"}},
		},
		UserPermissionsCheckFailedV1: {
			Payload:     &UserPermissionsCheckFailedPayloadV1{},
			Summary:     "User Permissions Check Failed",
			Description: "Permissions check failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "round"}},
		},

		// UDisc flow
		UpdateUDiscIdentityRequestedV1: {
			Payload:     &UpdateUDiscIdentityRequestedPayloadV1{},
			Summary:     "Update UDisc Identity Requested",
			Description: "Request to update a user's UDisc identity.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}},
		},
		UDiscIdentityUpdatedV1: {
			Payload:     &UDiscIdentityUpdatedPayloadV1{},
			Summary:     "UDisc Identity Updated",
			Description: "UDisc identity was updated.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},
		UDiscIdentityUpdateFailedV1: {
			Payload:     &UDiscIdentityUpdateFailedPayloadV1{},
			Summary:     "UDisc Identity Update Failed",
			Description: "UDisc identity update failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},

		UDiscMatchConfirmationRequiredV1: {
			Payload:     &UDiscMatchConfirmationRequiredPayloadV1{},
			Summary:     "UDisc Match Confirmation Required",
			Description: "Player matches require admin confirmation.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "user"}},
		},
		UDiscMatchConfirmedV1: {
			Payload:     &UDiscMatchConfirmedPayloadV1{},
			Summary:     "UDisc Match Confirmed",
			Description: "Admin confirmed player matches.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "round"}},
		},
	}
}
