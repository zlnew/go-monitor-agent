package postgres

import (
	"context"
	"fmt"

	"horizonx/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) domain.RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) SyncPermissions(ctx context.Context, data map[domain.RoleConst]map[domain.PermissionConst]bool) (err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	roleNameToID := make(map[string]int)
	permNameToID := make(map[string]int)

	// --- Upsert Roles ---
	for roleName := range data {
		var id int
		err := tx.QueryRow(ctx,
			`INSERT INTO roles (name)
			 VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
			 RETURNING id`,
			roleName,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed upsert role %s: %w", roleName, err)
		}
		roleNameToID[string(roleName)] = id
	}

	// --- Collect all unique permissions ---
	uniquePerms := make(map[string]struct{})
	for _, perms := range data {
		for permName := range perms {
			uniquePerms[string(permName)] = struct{}{}
		}
	}

	// --- Upsert Permissions ---
	for permName := range uniquePerms {
		var id int
		err := tx.QueryRow(ctx,
			`INSERT INTO permissions (name)
			 VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
			 RETURNING id`,
			permName,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed upsert permission %s: %w", permName, err)
		}
		permNameToID[permName] = id
	}

	// --- Delete pivot table ---
	if _, err := tx.Exec(ctx, `DELETE FROM role_has_permissions`); err != nil {
		return fmt.Errorf("failed to clear role_has_permissions: %w", err)
	}

	// --- Batch insert pivot ---
	batch := &pgx.Batch{}
	const query = `INSERT INTO role_has_permissions (role_id, permission_id) VALUES ($1, $2)`
	for roleName, perms := range data {
		roleID := roleNameToID[string(roleName)]
		for permName, allowed := range perms {
			if !allowed {
				continue
			}
			permID := permNameToID[string(permName)]
			batch.Queue(query, roleID, permID)
		}
	}

	br := tx.SendBatch(ctx, batch)
	if err = br.Close(); err != nil {
		return fmt.Errorf("failed to insert role_has_permissions batch: %w", err)
	}

	// --- Commit transaction ---
	return tx.Commit(ctx)
}
