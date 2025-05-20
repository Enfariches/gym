package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"health/internal/domain/models"
	ctxkey "health/lib/ctxkey"

	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) CreateStatistics(ctx context.Context, statProgress string, percentView, media_id int64) error {
	const op = "postgres.CreateStatistics"

	employee_id := ctx.Value(ctxkey.UserKey).(int64)
	departament_id := ctx.Value(ctxkey.DepartmentKey).(int64)

	if err := s.checkMediaId(media_id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, _ := goqu.
		Insert("statistics").
		Cols("progress", "percentage_view", "employee_id", "department_id", "media_id").
		Vals(goqu.Vals{statProgress, percentView, employee_id, departament_id, media_id}).
		ToSQL()

	_, err := s.db.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return nil
}

func (s *Storage) GetEmployeeStatistics(ctx context.Context, employee_id, media_id int64) (*models.Statistics, error) {
	const op = "postgres.GetEmployeeStatistics"

	query, args, _ := goqu.
		From("statistics").
		Select("id", "progress", "percentage_view", "created_at").
		Where(goqu.And(goqu.C("employee_id").Eq(employee_id), goqu.C("media_id").Eq(media_id))).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	employeeName, employeeSurname, err := s.getInfoEmployee(employee_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mediaTitle, err := s.getTitleMedia(media_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stat, err := scanStatistics(row, func(s *models.Statistics) {
		s.EmployeeName = employeeName
		s.EmployeeSurname = employeeSurname
		s.MediaTitle = mediaTitle
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return nil, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return stat, nil
}

func (s *Storage) ListMediaStatistics(ctx context.Context, media_id int64) ([]*models.Statistics, error) {
	const op = "postgres.ListMediaStatistics"

	query, args, _ := goqu.
		From("statistics").
		Select("id", "progress", "percentage_view", "employee_id", "created_at").
		Where(goqu.C("media_id").Eq(media_id)).
		ToSQL()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}
	defer rows.Close()

	mediaTitle, err := s.getTitleMedia(media_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var stats []*models.Statistics

	for rows.Next() {

		statistics := &models.Statistics{}
		var employee_id int64

		err := rows.Scan(&statistics.ID, &statistics.Progress, &statistics.PercentageView, &employee_id, &statistics.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		employeeName, employeeSurname, err := s.getInfoEmployee(employee_id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		statistics.MediaTitle = mediaTitle
		statistics.EmployeeName = employeeName
		statistics.EmployeeSurname = employeeSurname

		stats = append(stats, statistics)
	}

	return stats, nil
}

func (s *Storage) ListEmployeeStatistics(ctx context.Context, employee_id int64) ([]*models.Statistics, error) {
	const op = "postgres.ListEmployeeStatistics"

	query, args, _ := goqu.
		From("statistics").
		Select("id", "progress", "percentage_view", "media_id", "created_at").
		Where(goqu.C("employee_id").Eq(employee_id)).
		ToSQL()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}
	defer rows.Close()

	employeeName, employeeSurname, err := s.getInfoEmployee(employee_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var stats []*models.Statistics

	for rows.Next() {

		statistics := &models.Statistics{}
		var media_id int64

		err := rows.Scan(&statistics.ID, &statistics.Progress, &statistics.PercentageView, &media_id, &statistics.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		mediaTitle, err := s.getTitleMedia(media_id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		statistics.MediaTitle = mediaTitle
		statistics.EmployeeName = employeeName
		statistics.EmployeeSurname = employeeSurname

		stats = append(stats, statistics)
	}

	return stats, nil
}

func (s *Storage) ListDepartmentStatistics(ctx context.Context, department_id int64) ([]*models.Statistics, error) {
	const op = "postgres.ListDepartmentStatistics"

	query, args, _ := goqu.
		From("statistics").
		Select("id", "progress", "percentage_view", "employee_id", "media_id", "created_at").
		Where(goqu.C("department_id").Eq(department_id)).
		ToSQL()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}
	defer rows.Close()

	var stats []*models.Statistics

	for rows.Next() {

		statistics := &models.Statistics{}
		var media_id, employee_id int64

		err := rows.Scan(&statistics.ID, &statistics.Progress, &statistics.PercentageView, &employee_id, &media_id, &statistics.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		mediaTitle, err := s.getTitleMedia(media_id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		employeeName, employeeSurname, err := s.getInfoEmployee(employee_id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		statistics.MediaTitle = mediaTitle
		statistics.EmployeeName = employeeName
		statistics.EmployeeSurname = employeeSurname

		stats = append(stats, statistics)
	}

	return stats, nil
}

func (s *Storage) checkMediaId(media_id int64) error {
	const op = "postgres.checkMediaId"

	query, args, _ := goqu.
		From("mediafiles").
		Select("id").
		Where(goqu.C("id").Eq(media_id)).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: media with id %d not found", op, media_id)
		}
		return fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return nil
}

func (s *Storage) getInfoEmployee(employee_id int64) (string, string, error) {
	const op = "postgres.getInfoEmployee"

	query, args, _ := goqu.
		From("employees").
		Select("name", "surname").
		Where(goqu.C("id").Eq(employee_id)).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	var firstName, lastName string
	err := row.Scan(&firstName, &lastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("%s: %w", op, err)
		}
		return "", "", fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return firstName, lastName, nil
}

func (s *Storage) getTitleMedia(media_id int64) (string, error) {
	const op = "postgres.getTitleMedia"

	query, args, _ := goqu.
		From("mediafiles").
		Select("title").
		Where(goqu.C("id").Eq(media_id)).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	var title string
	err := row.Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return "", fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return title, nil
}
