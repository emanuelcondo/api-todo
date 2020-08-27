package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/emanuelcondo/api-todo/models"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	"github.com/emanuelcondo/api-todo/repositories"
	"github.com/emanuelcondo/api-todo/utils"
)

// Search TODO
func SearchTodos(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	params := repositories.SearchTodoParams{
		Text:  "",
		Page:  1,
		Limit: 25,
	}

	if val, ok := qs["search"]; ok && len(val) > 0 {
		params.Text = val[0]
	}
	if val, ok := qs["page"]; ok && len(val) > 0 {
		tmp, err := strconv.ParseUint(val[0], 10, 32)
		if err == nil && tmp > 0 {
			params.Page = uint(tmp)
		}
	}
	if val, ok := qs["limit"]; ok && len(val) > 0 {
		tmp, err := strconv.ParseUint(val[0], 10, 32)
		if err == nil && tmp > 0 && tmp <= 100 {
			params.Limit = uint(tmp)
		}
	}

	result, err := todoRepository.Search(params)

	if err != nil {
		fmt.Println(err)
		err := map[string]string{
			"message": "An error occurred when looking for todos.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
	} else {
		utils.Reply(w, r, http.StatusCreated, result)
	}
}

// Create TODO
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&todo)

	if err != nil {
		fmt.Println(err)
		err := map[string]string{
			"message": "An error occurred when parsing todo from body.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
		return
	}

	params := &models.Todo{
		Title:       todo.Title,
		Description: todo.Description,
	}

	result, err := todoRepository.Create(params)

	if err != nil {
		fmt.Println(err)
		err := map[string]string{
			"message": "An error occurred when creating todo.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
	} else {
		result := map[string]interface{}{
			"todo": result,
		}
		utils.Reply(w, r, http.StatusCreated, result)
	}
}

// Update TODO
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	tmp := chi.URLParam(r, "id")

	todoID, err := strconv.ParseUint(tmp, 10, 32)
	if err != nil {
		fmt.Println(err)
		err := map[string]string{
			"message": "Invalid ID.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
		return
	}

	var todo models.Todo

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&todo)

	if err != nil {
		fmt.Println(err)
		err := map[string]string{
			"message": "An error occurred when parsing todo from body.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
		return
	}

	params := &models.Todo{
		ID:          uint(todoID),
		Title:       todo.Title,
		Description: todo.Description,
	}

	result, err := todoRepository.Update(params)

	if err != nil {
		fmt.Println(err)
		err := map[string]string{
			"message": "An error occurred when updating todo.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
	} else {
		result := map[string]interface{}{
			"todo": result,
		}
		utils.Reply(w, r, http.StatusOK, result)
	}
}
