package domain

import (
	"context"
)

type Role struct {
	ID   int64     `json:"id"`
	Name RoleConst `json:"name"`
}

type RoleConst string

const (
	RoleAdmin  RoleConst = "admin"
	RoleViewer RoleConst = "viewer"
)

type RoleService interface {
	HasPermission(ctx context.Context, perm PermissionConst) error
	SyncPermissions(ctx context.Context) error
}

type RoleRepository interface {
	SyncPermissions(ctx context.Context, data map[RoleConst]map[PermissionConst]bool) error
}
