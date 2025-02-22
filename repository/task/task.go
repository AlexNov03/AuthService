package taskrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/AlexNov03/AuthService/errors/internalerr"
	"github.com/AlexNov03/AuthService/models"
)

type TaskRepo struct {
	DB *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{DB: db}
}

func (tr *TaskRepo) GetUserTasks(ctx context.Context, userID string) ([]*models.TaskInfo, error) {
	rows, err := tr.DB.QueryContext(ctx, `SELECT task_id, header, task_description, start_time, end_time
	FROM "task" WHERE user_id=$1`, userID)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, errors.New("TaskRepo.GetUserTasks: client canceled request")
		}
		return nil, internalerr.NewInternalError(http.StatusInternalServerError, fmt.Sprintf("TaskRepo.GetUserTasks query: %v", err))
	}

	defer rows.Close()

	userTasks := make([]*models.TaskInfo, 0)

	for rows.Next() {
		userTask := new(models.TaskInfo)

		err := rows.Scan(&userTask.ID, &userTask.Header, &userTask.Description, &userTask.StartTime, &userTask.EndTime)
		if err != nil {
			return nil, internalerr.NewInternalError(http.StatusInternalServerError, fmt.Sprintf("TaskRepo.GetUserTasks rows.Next: %v", err))
		}
		userTasks = append(userTasks, userTask)
	}
	return userTasks, nil
}

func (tr *TaskRepo) AddUserTask(ctx context.Context, userID string, task *models.Task) error {
	err := tr.DB.QueryRowContext(ctx, `SELECT 1 FROM "task" WHERE 
	NOT (end_time < $1 OR start_time > $2) LIMIT 1`, task.StartTime, task.EndTime).Scan(new(int))

	if err == nil {
		return internalerr.NewInternalError(http.StatusConflict, "task overlaps in time other tasks")
	}

	if errors.Is(err, context.Canceled) {
		return errors.New("TaskRepo.AddUserTask: client canceled request")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return internalerr.NewInternalError(http.StatusInternalServerError, fmt.Sprintf("TaskRepo.AddUserTask queryRow: %v", err))
	}

	_, err = tr.DB.ExecContext(ctx, `INSERT INTO "task" (user_id, header, task_description, start_time, end_time)
	VALUES ($1, $2, $3, $4, $5)`, userID, task.Header, task.Description, task.StartTime, task.EndTime)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("TaskRepo.AddUserTask: client cancelled request")
		}
		return internalerr.NewInternalError(http.StatusInternalServerError, fmt.Sprintf("TaskRepo.AddUserTask exec: %v", err))
	}
	return nil
}
