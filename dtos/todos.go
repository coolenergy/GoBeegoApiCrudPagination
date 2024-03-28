package dtos

import (
	"github.com/melardev/GoBeegoApiCrudPagination/models"
)

type CreateTodo struct {
	Title       string `form:"title" json:"title" xml:"title" binding:"required"`
	Description string `form:"description" json:"description" xml:"description"`
	Completed   bool   `form:"completed" json:"completed" xml:"completed"`
}

func GetTodoDto(todo *models.Todo, includeDescription bool) map[string]interface{} {
	dto := map[string]interface{}{
		"id":    todo.Id,
		"title": todo.Title,
	}

	if includeDescription {
		dto["description"] = todo.Description
	}

	dto["completed"] = todo.Completed
	dto["created_at"] = todo.CreatedAt
	dto["updated_at"] = todo.UpdatedAt

	return dto
}

func GetTodoListDto(todos []models.Todo) []interface{} {

	result := make([]interface{}, len(todos))

	for i := 0; i < len(todos); i++ {
		result[i] = GetTodoDto(&todos[i], false)
	}

	return result
}

func CreateTodoPagedResponse(path string, todos []models.Todo, page, pageSize int, count int64) map[string]interface{} {
	var resources = make([]interface{}, len(todos))
	for index, todo := range todos {
		resources[index] = GetTodoDto(&todo, false)
	}
	return CreatePagedResponse(path, resources, "todos", page, pageSize, count)
}

func CreateTodoCreatedDto(todo *models.Todo) interface{} {
	return CreateSuccessWithDtoAndMessageDto(GetTodoDto(todo, true), "Todo created successfully")
}

func CreateTodoUpdatedDto(todo *models.Todo) interface{} {
	return CreateSuccessWithDtoAndMessageDto(GetTodoDto(todo, false), "Todo updated successfully")
}
