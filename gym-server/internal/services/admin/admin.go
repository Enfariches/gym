package admin

import (
	"context"
	"fmt"
	"health/internal/domain/models"
	"health/lib/logger/sl"
	"log/slog"
)

type Admin struct {
	log          *slog.Logger
	adminManager AdminManager
}

type AdminManager interface {
	Admin(ctx context.Context, admin_id int64) (*models.Admin, error)
	UpdateAdmin(ctx context.Context, updateFields map[string]any) (*models.Admin, error)
	ListAdminEmployees(ctx context.Context, department_id int64) ([]*models.Employee, error)
}

func New(log *slog.Logger, adminManager AdminManager) *Admin {
	return &Admin{
		log:          log,
		adminManager: adminManager,
	}
}

func (a *Admin) GetAdmin(ctx context.Context, admin_id int64) (*models.Admin, error) {
	const op = "admin.GetAdmin"
	log := a.log.With("op", op)

	admin, err := a.adminManager.Admin(ctx, admin_id)
	if err != nil {
		log.Error("failed to get admin", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return admin, nil
}

func (a *Admin) UpdateAdmin(ctx context.Context, updateFields map[string]any) (*models.Admin, error) {
	const op = "admin.UpdateAdmin"
	log := a.log.With("op", op)

	updatedAdmin, err := a.adminManager.UpdateAdmin(ctx, updateFields)
	if err != nil {
		log.Error("failed to update admin", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedAdmin, nil
}

func (a *Admin) ListAdminEmployees(ctx context.Context, department_id int64) ([]*models.Employee, error) {
	const op = "admin.ListAdminEmployees"
	log := a.log.With("op", op)

	employees, err := a.adminManager.ListAdminEmployees(ctx, department_id)
	if err != nil {
		log.Error("failed to list admin employees", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return employees, nil
}
