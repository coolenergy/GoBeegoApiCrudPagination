package services

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/melardev/GoBeegoApiCrudPagination/models"
)

func FetchTodos(page, pageSize int) ([]models.Todo, int64, error) {
	tableName := new(models.Todo).TableName()
	o := orm.NewOrm()
	// var todos []models.Todo
	// o.QueryTable(new(models.Todo)).All(&todos)

	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, 0, err
	}

	var temp []models.Todo
	totalTodosCount, err := o.Raw(qb.Select(fmt.Sprintf("%s.id", tableName)).From(tableName).String()).
		QueryRows(&temp)

	if err != nil {
		return nil, 0, err
	}

	qb, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, 0, err
	}

	var todos []models.Todo

	qb.Select(
		fmt.Sprintf("%s.id", tableName),
		fmt.Sprintf("%s.title", tableName),
		fmt.Sprintf("%s.completed", tableName),
		fmt.Sprintf("%s.created_at", tableName),
		fmt.Sprintf("%s.updated_at", tableName)).
		From(tableName).
		OrderBy("created_at").Desc().
		Limit(pageSize).Offset((page - 1) * pageSize)

	rawQuery := qb.String()
	count, err := o.Raw(rawQuery).QueryRows(&todos)
	println("fetched", count, "rows")
	return todos, totalTodosCount, err
}

func FetchPendingTodos(page, pageSize int) ([]models.Todo, int64, error) {
	return FetchTodosByCompleted(false, page, pageSize)
}

func FetchCompletedTodos(page, pageSize int) ([]models.Todo, int64, error) {
	return FetchTodosByCompleted(true, page, pageSize)
}

func FetchTodosByCompleted(completed bool, page, pageSize int) ([]models.Todo, int64, error) {
	tableName := new(models.Todo).TableName()
	o := orm.NewOrm()
	// Retrieve the Todos Count
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, 0, err
	}

	qb.Select(
		fmt.Sprintf("%s.id", tableName)).
		From(tableName).
		Where("completed = ?")

	totalTodosCount, err := o.Raw(qb.String(), completed).QueryRows(&[]models.Todo{})
	if err != nil {
		return nil, 0, err
	}

	// Retrieve the Todos
	var todos []models.Todo
	qb, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, 0, err
	}
	qb.Select(
		fmt.Sprintf("%s.id", tableName),
		fmt.Sprintf("%s.title", tableName),
		fmt.Sprintf("%s.completed", tableName),
		fmt.Sprintf("%s.created_at", tableName),
		fmt.Sprintf("%s.updated_at", tableName)).
		From(tableName).
		Where("completed = ?").
		OrderBy("created_at").Desc().
		Limit(pageSize).Offset((page - 1) * pageSize)

	rawQuery := qb.String()
	count, err := o.Raw(rawQuery, completed).QueryRows(&todos)
	if err != nil {
		return nil, 0, err
	}

	println("fetched", count, "rows")
	return todos, totalTodosCount, err
}

func FetchById(id int) (todo *models.Todo, err error) {
	o := orm.NewOrm()
	todo = &models.Todo{Id: id}
	err = o.Read(todo)
	return
}

func CreateTodo(title, description string, completed bool) (todo *models.Todo, err error) {
	todo = &models.Todo{
		Title:       title,
		Description: description,
		Completed:   completed,
	}

	o := orm.NewOrm()
	qs := o.QueryTable(todo)
	inserter, err := qs.PrepareInsert()
	if err != nil {
		return
	}

	id, err := inserter.Insert(todo)

	// We assume the Todo has been saved as such in the Db, but it may not be
	// For example if a string truncation took place
	todo.Id = int(id)
	return todo, err
}

func UpdateTodo(id int, title, description string, completed bool) (todo *models.Todo, err error) {
	todo, err = FetchById(id)

	if err != nil {
		return
	}

	todo.Title = title

	// TODO: handle this in a better way, the user should be able to set description to empty string
	// The intention is to check against nil but in go there are no nil strings, so we can not know
	// if the user intended to update the description to empty string or just update the other fields other than description.
	if description != "" {
		todo.Description = description
	}

	todo.Completed = completed
	o := orm.NewOrm()

	count, err := o.Update(todo)
	println("updated", count, "rows")
	return
}

func DeleteTodo(todo *models.Todo) error {
	o := orm.NewOrm()
	count, err := o.Delete(todo)
	println("deleted", count, "rows")
	return err
}

func DeleteAllTodos() (err error) {
	o := orm.NewOrm()
	todo := models.Todo{}
	count, err := o.QueryTable(todo.TableName()).Filter("id__isnull", false).Delete()
	println("Deleted", count, "rows")
	return err
}
