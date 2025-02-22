package taskuc

import (
	"context"

	"github.com/AlexNov03/AuthService/models"
)

type TaskRepository interface {
	GetUserTasks(ctx context.Context, userID string) ([]*models.TaskInfo, error)
	AddUserTask(ctx context.Context, userID string, task *models.Task) error
}

type TaskUsecase struct {
	TD TaskRepository
}

func NewTaskUsecase(td TaskRepository) *TaskUsecase {
	return &TaskUsecase{TD: td}
}

func (tu *TaskUsecase) GetUserTasks(ctx context.Context, userID string) ([]*models.TaskInfo, error) {
	taskInfo, err := tu.TD.GetUserTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return taskInfo, nil
}
func (tu *TaskUsecase) AddUserTask(ctx context.Context, userID string, task *models.Task) error {
	err := tu.TD.AddUserTask(ctx, userID, task)
	if err != nil {
		return err
	}
	return nil
}
