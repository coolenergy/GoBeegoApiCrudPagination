package controllers

import (
	"encoding/json"
	"github.com/melardev/GoBeegoApiCrudPagination/dtos"
	"github.com/melardev/GoBeegoApiCrudPagination/models"
	"github.com/melardev/GoBeegoApiCrudPagination/services"
	"net/http"
	"strconv"
)

type TodosController struct {
	BaseController
}

func (this *TodosController) GetAllTodos() {
	page, pageSize := this.getPagingParams()
	todos, totalTodoCount, err := services.FetchTodos(page, pageSize)
	if err != nil {
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
	}
	this.Data["json"] = dtos.CreateTodoPagedResponse(this.Ctx.Request.URL.Path, todos, page, pageSize, totalTodoCount)
	this.ServeJSON()
}

func (this *TodosController) GetAllPendingTodos() {
	page, pageSize := this.getPagingParams()
	todos, totalTodoCount, err := services.FetchPendingTodos(page, pageSize)
	if err != nil {
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
	}
	this.Data["json"] = dtos.CreateTodoPagedResponse(this.Ctx.Request.URL.Path, todos, page, pageSize, totalTodoCount)
	this.ServeJSON()
}
func (this *TodosController) GetAllCompletedTodos() {
	page, pageSize := this.getPagingParams()
	todos, totalTodoCount, err := services.FetchCompletedTodos(page, pageSize)
	if err != nil {
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
	}

	this.Data["json"] = dtos.CreateTodoPagedResponse(this.Ctx.Request.URL.Path, todos, page, pageSize, totalTodoCount)
	this.ServeJSON()
}

func (this *TodosController) GetTodoById() {
	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	todo, err := services.FetchById(id)
	if err != nil {
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
		return
	}

	this.Data["json"] = dtos.GetSuccessTodoDto(todo)
	this.ServeJSON()
}

func (this *TodosController) CreateTodo() {
	todo := &models.Todo{}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, todo); err != nil {
		// this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
		return
	}

	todo, err := services.CreateTodo(todo.Title, todo.Description, todo.Completed)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
		return
	}

	this.Data["json"] = dtos.CreateTodoCreatedDto(todo)
	this.ServeJSON()
}

func (this *TodosController) UpdateTodo() {
	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.Data["json"] = dtos.CreateErrorDtoWithMessage("You must set an ID")
		this.ServeJSON()
		return
	}

	var todoInput models.Todo
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &todoInput); err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
		return
	}

	todo, err := services.UpdateTodo(id, todoInput.Title, todoInput.Description, todoInput.Completed)
	if err != nil {
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
		return
	}

	this.Data["json"] = dtos.CreateTodoUpdatedDto(todo)
	this.ServeJSON()
}

func (this *TodosController) DeleteTodo() {
	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.Data["json"] = dtos.CreateErrorDtoWithMessage("You must set an ID")
		this.ServeJSON()
		return
	}

	todo, err := services.FetchById(id)

	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Data["json"] = dtos.CreateErrorDtoWithMessage("todo not found")
		this.ServeJSON()

		return
	}

	err = services.DeleteTodo(todo)

	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Data["json"] = dtos.CreateErrorDtoWithMessage("Could not delete Todo")
		this.ServeJSON()
		return
	}

	this.Data["json"] = dtos.CreateSuccessWithMessageDto("todo deleted successfully")
	this.ServeJSON()
}

func (this *TodosController) DeleteAllTodos() {
	err := services.DeleteAllTodos()
	if err != nil {
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		this.ServeJSON()
	}
	this.Data["json"] = dtos.CreateSuccessWithMessageDto("All Todos deleted succesfully")
	this.ServeJSON()
}
