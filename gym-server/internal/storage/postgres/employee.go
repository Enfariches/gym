package postgres

import (
	"context"
	"fmt"
	"health/internal/domain/models"
	"health/internal/storage"
	ctxkey "health/lib/ctxkey"

	"github.com/doug-martin/goqu/v9"
)

var filedsTable = []interface{}{"id", "name", "second_name", "surname", "age", "sex", "phone", "email", "post"}

func (s *Storage) Employee(ctx context.Context, employee_id int64) (*models.Employee, error) {
	const op = "postgres.Employee"

	department_id := ctx.Value(ctxkey.DepartmentKey).(int64)

	query, args, _ := goqu.
		From("employees").
		Select(filedsTable...).
		Where(goqu.Ex{"id": employee_id}).
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	departmentName, err := s.GetDepNameByDepID(department_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department by employee id: %w", err)
	}

	employee, err := scanEmployee(row, func(s *models.Employee) { s.Department = departmentName })
	if err != nil {
		if HandleDBError(err) == storage.ErrUserNotFound {
			return &models.Employee{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return &models.Employee{}, fmt.Errorf("%s: %w", op, err)
	}

	return employee, nil
}

func (s *Storage) UpdateEmployee(ctx context.Context, updateFields map[string]any) (*models.Employee, error) {
	const op = "postgres.UpdateEmployee"

	employee_id := ctx.Value(ctxkey.UserKey).(int64)
	department_id := ctx.Value(ctxkey.DepartmentKey).(int64)

	if deptName, ok := updateFields["department"]; ok {
		err := s.updateDepNameByDepID(department_id, deptName.(string))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		delete(updateFields, "department")
	}

	if len(updateFields) == 0 {
		employee, err := s.Employee(ctx, employee_id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return employee, nil
	}
	query, args, _ := goqu.
		Update("employees").
		Set(updateFields).
		Where(goqu.C("id").Eq(employee_id)).
		Returning(filedsTable...).
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	departmentName, err := s.GetDepNameByDepID(department_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department by employee id: %w", err)
	}

	updatedEmployee, err := scanEmployee(row, func(s *models.Employee) { s.Department = departmentName })
	if err != nil {
		return &models.Employee{}, fmt.Errorf("%s: %w", op, err)
	}


	return updatedEmployee, nil
}

func (s *Storage) DeleteEmployee(ctx context.Context, employee_id int64) error {
	const op = "postgres.DeleteEmployee"

	query, args, _ := goqu.Delete("employees").Where(goqu.C("id").Eq(employee_id)).ToSQL()

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) ListDepartments(ctx context.Context) ([]*models.Department, error) {
	const op = "postgres.ListDepartments"

	query, _, _ := goqu.
		From("departments").
		Select("id", "name").
		ToSQL()

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	resultDepartments := make([]*models.Department, 0)

	for rows.Next() {
		d, err := scanDepartment(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		resultDepartments = append(resultDepartments, d)
	}

	return resultDepartments, nil
}
