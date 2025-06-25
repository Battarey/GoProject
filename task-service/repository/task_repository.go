package repository

import (
	"task-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetTaskByID(id string) (*model.Task, error) {
	var task model.Task
	taskID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	if err := r.db.Where("id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) UpdateTask(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) DeleteTask(id string) error {
	taskID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&model.Task{}, "id = ?", taskID).Error
}

func (r *TaskRepository) ListTasks(status, assigneeID string, offset, limit int) ([]model.Task, error) {
	var tasks []model.Task
	query := r.db.Model(&model.Task{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if assigneeID != "" {
		if id, err := uuid.Parse(assigneeID); err == nil {
			query = query.Where("assignee_id = ?", id)
		}
	}
	if err := query.Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) ChangeStatus(id, status string) error {
	taskID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Task{}).Where("id = ?", taskID).Update("status", status).Error
}
