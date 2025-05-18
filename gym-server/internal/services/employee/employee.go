package employee

import (
	"context"
	"health/internal/domain/models"
	"health/lib/logger/sl"
	"log/slog"
)

type Employee struct {
	log             *slog.Logger
	employeeManager EmployeeManager
}

type EmployeeManager interface {
	Employee(ctx context.Context, employee_id int64) (*models.Employee, error)
	UpdateEmployee(ctx context.Context, updateFields map[string]any) (*models.Employee, error)
	DeleteEmployee(ctx context.Context, employee_id int64) error

	ListDepartments(ctx context.Context) ([]*models.Department, error)
}

func New(log *slog.Logger, employeeManager EmployeeManager) *Employee {
	return &Employee{
		log:             log,
		employeeManager: employeeManager,
	}
}

func (e *Employee) GetEmployee(ctx context.Context, employee_id int64) (*models.Employee, error) {
	const op = "employee.GetEmployee"
	log := e.log.With("op", op)

	employee, err := e.employeeManager.Employee(ctx, employee_id)
	if err != nil {
		log.Error("failed to get employee", sl.Err(err))
		return nil, err
	}

	return employee, nil
}

func (e *Employee) UpdateEmployee(ctx context.Context, updateFields map[string]any) (*models.Employee, error) {
	const op = "employee.UpdateEmployee"
	log := e.log.With("op", op)

	updatedEmployee, err := e.employeeManager.UpdateEmployee(ctx, updateFields)
	if err != nil {
		log.Error("failed to update employee", sl.Err(err))
		return nil, err
	}

	return updatedEmployee, nil
}

func (e *Employee) DeleteEmployee(ctx context.Context, employee_id int64) error {
	const op = "employee.DeleteEmployee"
	log := e.log.With("op", op)

	err := e.employeeManager.DeleteEmployee(ctx, employee_id)
	if err != nil {
		log.Error("failed to delete employee", sl.Err(err))
		return err
	}

	return nil
}

func (e *Employee) ListDepartments(ctx context.Context) ([]*models.Department, error) {
	const op = "employee.ListDepartments"
    log := e.log.With("op", op)

    departments, err := e.employeeManager.ListDepartments(ctx)
    if err != nil {
        log.Error("failed to list departments", sl.Err(err))
        return nil, err
    }

    return departments, nil
}
