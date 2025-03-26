package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/config"
	models "github.com/YuriGarciaRibeiro/zipzop-chat/internal/model"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	issuer    string
	expiry    time.Duration
}

func NewAuthService(repo *repository.UserRepository, cfg *config.AppConfig) *AuthService {
	return &AuthService{
		userRepo:  repo,
		jwtSecret: cfg.Auth.JWTSecret,
		issuer:    cfg.Auth.TokenIssuer,
		expiry:    cfg.Auth.JWTExpiry,
	}
}

func (s *AuthService) RegisterUser(user *models.User) (*models.PublicUser, string, error) {
	// 1. Verificar se usuário já existe
	if exists, err := s.userRepo.EmailExists(user.Email); err != nil || exists {
		return nil, "", errors.New("email já está em uso")
	}

	// 2. Criar hash da senha (se ainda não foi criado)
	if user.PasswordHash == "" {
		hashedPwd, err := bcrypt.GenerateFromPassword(
			[]byte(user.PasswordHash+user.Salt),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return nil, "", err
		}
		user.PasswordHash = string(hashedPwd)
	}

	// 3. Persistir no banco
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, "", err
	}

	// 4. Gerar token
	token, err := s.GenerateJWT(user.ID.String())
	if err != nil {
		return nil, "", err
	}

	// 5. Retornar usuário público + token
	return user.ToPublic(), token, nil
}

func (s *AuthService) LoginUser(loginRequest *models.LoginUserRequest) (*models.PublicUser, string, error) {
	user, err := s.userRepo.GetUser(loginRequest.Email)
	if err != nil {
		fmt.Println(err)
		return nil, "", errors.New(err.Error())
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(loginRequest.Password+user.Salt),
	)
	if err != nil {
		return nil, "", errors.New("credenciais inválidas")
	}

	token, err := s.GenerateJWT(user.ID.String())
	if err != nil {
		return nil, "", err
	}

	return user.ToPublic(), token, nil
}

// Métodos auxiliares mantidos com melhorias
func (s *AuthService) GenerateJWT(userID string) (string, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"iss":   s.issuer,
		"exp":   time.Now().Add(s.expiry).Unix(),
		"iat":   time.Now().Unix(),
		"email": user.Email,
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de signing inesperado: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", err
	}


	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return claims["email"].(string), nil
	}

	return "", fmt.Errorf("token inválido")
}
