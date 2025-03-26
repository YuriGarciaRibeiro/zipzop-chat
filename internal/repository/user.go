package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/config"
	models "github.com/YuriGarciaRibeiro/zipzop-chat/internal/model"
	_ "github.com/lib/pq" // Import do driver PostgreSQL
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(cfg *config.DatabaseConfig) (*UserRepository, error) {
	// Monta a URL de conexão
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DatabaseName,
		cfg.SSLMode,
	)

	// Abre conexão com o banco
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configuração do pool de conexões
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxConns / 2)
	db.SetConnMaxLifetime(time.Minute * 30)

	// Testa a conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Close() error {
	return r.db.Close()
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, phone, password_hash, salt, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.Salt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, phone, password_hash, salt, created_at, updated_at
		FROM users WHERE id = $1`
	
	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Salt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // Usuário não encontrado
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}


func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, phone, password_hash, salt, created_at, updated_at
		FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Salt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // Usuário não encontrado
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// Método adicional recomendado
func (r *UserRepository) GetUser(email string) (*models.User, error) {
	fmt.Println("Buscando usuário por email:", email)

	query := `
        SELECT 
            id, 
            email, 
            phone, 
            password_hash, 
            salt, 
            created_at, 
            updated_at 
        FROM users 
        WHERE email = $1`

	var user models.User

	// Ordem DEVE corresponder exatamente à query SQL
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Salt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Usuário não encontrado
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %v", err)
	}

	return &user, nil
}
