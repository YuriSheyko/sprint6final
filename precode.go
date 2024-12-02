package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

func HendlerGetAllTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	buffer := new(bytes.Buffer)
	buffer.Reset()
	err := json.NewEncoder(buffer).Encode(tasks)

	if err != nil {
		fmt.Printf("ошибка при сериализации в json: %s\n", err.Error())
		res.WriteHeader(400)
		return
	}

	res.WriteHeader(201)

	res.Write(buffer.Bytes())

}

func HendlerPostTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		fmt.Printf("ошибка при десериализации из json: %s\n", err.Error())
		res.WriteHeader(400)
		return
	}
	if _, ok := tasks[task.ID]; ok {
		fmt.Println("элемент уже есть в массиве")
		res.WriteHeader(400)
		return
	}
	tasks[task.ID] = task
	res.WriteHeader(201)

}

func HendlerGetIDtask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(req, "id")
	var task *Task = nil

	for _, item := range tasks {
		if item.ID == id {
			task = &item
			break
		}
	}
	if task == nil {
		res.WriteHeader(400)
		return
	}
	buffer := new(bytes.Buffer)
	buffer.Reset()
	err := json.NewEncoder(buffer).Encode(task)

	if err != nil {
		fmt.Printf("ошибка при сериализации в json: %s\n", err.Error())
		res.WriteHeader(400)
		return
	}

	res.WriteHeader(201)

	res.Write(buffer.Bytes())

}

func HendlerDeletTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(req, "id")

	_, ok := tasks[id]

	if !ok {
		fmt.Println("Элемента нет в маппе")
		res.WriteHeader(400)
		return
	}
	delete(tasks, id)
	res.WriteHeader(201)

}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...

	r.Get("/tasks", HendlerGetAllTasks)
	r.Post("/tasks", HendlerPostTask)
	r.Get("/tasks/{id}", HendlerGetIDtask)
	r.Delete("/tasks/{id}", HendlerDeletTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
