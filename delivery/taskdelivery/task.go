package taskdelivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexNov03/AuthService/errors/externalerr"
	"github.com/AlexNov03/AuthService/models"
	"github.com/go-playground/validator/v10"
)

type TaskUsecase interface {
	GetUserTasks(ctx context.Context, userID string) ([]*models.TaskInfo, error)
	AddUserTask(ctx context.Context, userID string, task *models.Task) error
}

type TaskDelivery struct {
	TUC       TaskUsecase
	Validator *validator.Validate
}

func NewTaskDelivery(tuc TaskUsecase) *TaskDelivery {
	return &TaskDelivery{TUC: tuc, Validator: validator.New(validator.WithRequiredStructEnabled())}
}

func (td *TaskDelivery) AddTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		externalerr.ProcessBadRequestError(w, "user id is not valid")
		return
	}

	task := new(models.Task)
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		externalerr.ProcessBadRequestError(w, err.Error())
		return
	}

	err = td.Validator.Struct(task)
	if err != nil {
		externalerr.ProcessBadRequestError(w, err.Error())
		return
	}

	fmt.Println(task)

	err = td.TUC.AddUserTask(ctx, userID, task)
	if err != nil {
		externalerr.ProcessError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})

}

func (td *TaskDelivery) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		externalerr.ProcessBadRequestError(w, "user id is not valid")
		return
	}

	userTasks, err := td.TUC.GetUserTasks(ctx, userID)
	if err != nil {
		externalerr.ProcessError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userTasks)
}
