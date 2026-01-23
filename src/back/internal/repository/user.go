package repository

import (
	"database/sql"
	"errors"

	"github.com/CreateLab/laritmo/internal/models"
	sq "github.com/Masterminds/squirrel"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// userColumns returns the list of columns for user queries
func userColumns() []string {
	return []string{
		"id", "email", "username", "password_hash", "role",
		"is_active", "last_login_at", "created_at", "updated_at",
	}
}

// scanUser scans a row into a User struct
func scanUser(row interface{ Scan(...any) error }) (*models.User, error) {
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query, args, _ := sq.Select(userColumns()...).
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()

	return scanUser(r.db.QueryRow(query, args...))
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query, args, _ := sq.Select(userColumns()...).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	return scanUser(r.db.QueryRow(query, args...))
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query, args, _ := sq.Select(userColumns()...).
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()

	return scanUser(r.db.QueryRow(query, args...))
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	query, args, _ := sq.Select(userColumns()...).
		From("users").
		OrderBy("id ASC").
		ToSql()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) GetAllAdmins() ([]*models.User, error) {
	query, args, _ := sq.Select(userColumns()...).
		From("users").
		Where(sq.Eq{"role": []string{models.RoleOwner, models.RoleAdmin}}).
		OrderBy("id ASC").
		ToSql()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) Create(user *models.User) error {
	query, args, _ := sq.Insert("users").
		Columns("email", "username", "password_hash", "role", "is_active").
		Values(user.Email, user.Username, user.PasswordHash, user.Role, user.IsActive).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	query, args, _ := sq.Update("users").
		Set("email", user.Email).
		Set("username", user.Username).
		Set("role", user.Role).
		Set("is_active", user.IsActive).
		Where(sq.Eq{"id": user.ID}).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	query, args, _ := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdatePassword(id int, passwordHash string) error {
	query, args, _ := sq.Update("users").
		Set("password_hash", passwordHash).
		Where(sq.Eq{"id": id}).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateLastLogin(id int) error {
	query, args, _ := sq.Update("users").
		Set("last_login_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.Eq{"id": id}).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
