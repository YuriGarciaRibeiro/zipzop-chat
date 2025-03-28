package router

import (
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/auth"
)

func NewAuthRoutes(authHandler *auth.AuthHandler) []Route {
	return []Route{
		{
			URI:      "/auth/login",
			Method:   "POST",
			Function: authHandler.Login,
			IsAuth:   false,
		},
		{
			URI:      "/auth/register",
			Method:   "POST",
			Function: authHandler.Register,
			IsAuth:   false,
		},
	}
}

