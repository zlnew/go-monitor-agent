package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"horizonx-server/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers(ctx context.Context, opts domain.ListOptions) ([]*domain.User, int64, error) {
	baseQuery := `
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.deleted_at IS NULL
	`
	args := []any{}

	if opts.Search != "" {
		baseQuery += " AND (u.email LIKE ? OR u.name LIKE ?)"
		searchParam := "%" + opts.Search + "%"
		args = append(args, searchParam, searchParam)
	}

	var total int64
	if opts.IsPaginate {
		countQuery := "SELECT COUNT(*) " + baseQuery
		if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("failed to count users: %w", err)
		}
	}

	selectQuery := `
		SELECT 
			u.id, u.name, u.email, u.password, u.role_id, u.created_at, u.updated_at,
			r.id, r.name
	` + baseQuery

	selectQuery += " ORDER BY u.created_at DESC"

	if opts.IsPaginate {
		offset := (opts.Page - 1) * opts.Limit
		selectQuery += " LIMIT ? OFFSET ?"
		args = append(args, opts.Limit, offset)
	} else {
		selectQuery += " LIMIT 1000"
	}

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var role domain.Role

		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt, &user.UpdatedAt,
			&role.ID, &role.Name,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}

		user.Role = &role
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, ID int64) (*domain.User, error) {
	query := `
		SELECT 
			u.id, u.name, u.email, u.password, u.role_id, u.created_at, u.updated_at,
			r.id, r.name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = ? AND u.deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, ID)

	var user domain.User
	var role domain.Role

	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt, &user.UpdatedAt,
		&role.ID, &role.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user.Role = &role
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT 
			u.id, u.name, u.email, u.password, u.role_id, u.created_at, u.updated_at,
			r.id, r.name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = ? AND u.deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	var role domain.Role

	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt, &user.UpdatedAt,
		&role.ID, &role.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user.Role = &role
	return &user, nil
}

func (r *UserRepository) GetRoleByID(ctx context.Context, roleID int64) (*domain.Role, error) {
	query := `SELECT id, name FROM roles WHERE id = ?`

	var role domain.Role
	err := r.db.QueryRowContext(ctx, query, roleID).Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRoleNotFound
		}
		return nil, err
	}

	return &role, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password, role_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.RoleID, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User, userID int64) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, password = ?, role_id = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.RoleID, now, userID)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found or deleted", userID)
	}

	user.UpdatedAt = now

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	query := `UPDATE users SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to execute soft delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found or already deleted", userID)
	}

	return nil
}
