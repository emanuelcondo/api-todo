package repositories

import (
	"github.com/emanuelcondo/api-todo/db"
	"github.com/emanuelcondo/api-todo/models"
	"math"
	"strings"
)

type TodoRepository struct{}

type SearchTodoParams struct {
	Text  string
	Page  uint
	Limit uint
}

type TodoPagination struct {
	TotalCount  uint          `json:"total_count"`
	TotalPages  uint          `json:"total_pages"`
	CurrentPage uint          `json:"current_page"`
	ResultLimit uint          `json:"result_limit"`
	Todos       []models.Todo `json:"todos"`
}

// Create Todo
func (todoRepository *TodoRepository) Create(todo *models.Todo) (*models.Todo, error) {
	DB := db.GetDBConnection()
	createdTodo := DB.Create(todo)

	if createdTodo.Error != nil {
		return nil, createdTodo.Error
	} else {
		return todo, nil
	}
}

// Update Todo
func (todoRepository *TodoRepository) Update(todo *models.Todo) (*models.Todo, error) {
	DB := db.GetDBConnection()
	updatedTodo := DB.Model(todo).Where("id = ?", todo.ID).Update(map[string]interface{}{"title": todo.Title, "description": todo.Description})

	if updatedTodo.Error != nil {
		return nil, updatedTodo.Error
	}

	updatedTodo = DB.First(todo, todo.ID)

	if updatedTodo.Error != nil {
		return nil, updatedTodo.Error
	} else {
		return todo, nil
	}
}

// Search Todos
func (todoRepository *TodoRepository) Search(params SearchTodoParams) (interface{}, error) {
	DB := db.GetDBConnection()
	var todos []models.Todo
	var count uint
	offset := params.Limit * (params.Page - 1)
	limit := params.Limit

	queryCount := DB.Find(&todos)
	queryRows := DB.Find(&todos)

	if params.Text != "" && strings.TrimSpace(params.Text) != "" {
		text := strings.ToLower(params.Text)
		text = "%" + strings.TrimSpace(text) + "%"
		queryCount = queryCount.Where("LOWER(title) LIKE ?", text).Or("LOWER(description) LIKE ?", text)
		queryRows = queryRows.Where("LOWER(title) LIKE ?", text).Or("LOWER(description) LIKE ?", text)
	}
	countResult := queryCount.Count(&count)

	if countResult.Error != nil {
		return nil, countResult.Error
	}

	rowsResult := queryRows.Order("created_at ASC").Offset(offset).Limit(limit).Scan(&todos)

	if rowsResult.Error != nil {
		return nil, rowsResult.Error
	} else {
		result := TodoPagination{
			TotalCount:  count,
			TotalPages:  uint(math.Ceil(float64(count) / float64(params.Limit))),
			CurrentPage: params.Page,
			ResultLimit: params.Limit,
			Todos:       todos,
		}
		return result, nil
	}
}
