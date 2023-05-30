package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Camada de domínio
type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	GetUser(userID int) (*User, error)
	GetAllUsers() ([]*User, error)
	SaveUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(userID int) error
}

// Camada de infraestrutura
type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(dbFile string) (*SQLiteUserRepository, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT
		)
	`)
	if err != nil {
		return nil, err
	}
	return &SQLiteUserRepository{db: db}, nil
}

func (r *SQLiteUserRepository) GetUser(userID int) (*User, error) {
	row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) GetAllUsers() ([]*User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *SQLiteUserRepository) SaveUser(user *User) error {
	result, err := r.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return err
	}
	user.ID, _ = result.LastInsertId()
	return nil
}

func (r *SQLiteUserRepository) UpdateUser(user *User) error {
	_, err := r.db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

func (r *SQLiteUserRepository) DeleteUser(userID int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", userID)
	return err
}

// Camada de aplicação
type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUser(userID int) (*User, error) {
	return s.userRepository.GetUser(userID)
}

func (s *UserService) GetAllUsers() ([]*User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *UserService) CreateUser(name, email string) error {
	user := &User{
		Name:  name,
		Email: email,
	}
	return s.userRepository.SaveUser(user)
}

func (s *UserService) UpdateUser(userID int, name, email string) error {
	user, err := s.userRepository.GetUser(userID)
	if err != nil {
		return err
	}
	user.Name = name
	user.Email = email
	return s.userRepository.UpdateUser(user)
}

func (s *UserService) DeleteUser(userID int) error {
	return s.userRepository.DeleteUser(userID)
}

// Exemplo de uso
func main() {
	dbFile := "users.db"
	repository, err := NewSQLiteUserRepository(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	service := NewUserService(repository)

	// Criação de usuário
	err = service.CreateUser("John Doe", "john@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Obtenção de usuário por ID
	user, err := service.GetUser(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.ID, user.Name, user.Email)

	// Obtenção de todos os usuários
	users, err := service.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		fmt.Println(user.ID, user.Name, user.Email)
	}

	// Atualização de usuário
	err = service.UpdateUser(1, "John Smith", "john.smith@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Deleção de usuário
	err = service.DeleteUser(1)
	if err != nil {
		log.Fatal(err)
	}
}
