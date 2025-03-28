package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userRequest model.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := mapRegisterRequestToUser(userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publicUser, token, err := h.service.RegisterUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(publicUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"user":  string(jsonUser),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest model.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.service.LoginUser(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  map[string]interface{}{"id": user.ID.String(), "email": user.Email},
	})
}

func mapRegisterRequestToUser(req model.RegisterUserRequest) (*model.User, error) {
	salt := GenerateRandomSalt() // Implemente esta função
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password+salt),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
		Salt:         salt,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}, nil
}
