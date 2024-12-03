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
		http.Error(res, "serialization error", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)

	res.Write(buffer.Bytes())

}

func HendlerPostTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		fmt.Printf("ошибка при десериализации из json: %s\n", err.Error())
		http.Error(res, "deserialization error", http.StatusBadRequest)
		return
	}
	if _, ok := tasks[task.ID]; ok {
		fmt.Println("элемент уже есть в массиве")
		http.Error(res, "element allready exists", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	res.WriteHeader(http.StatusCreated)
}

func HendlerGetIDtask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(req, "id")
	task, ok := tasks[id]

	if !ok {
		fmt.Println("Элемента нет в маппе")
		http.Error(res, "task not found", http.StatusBadRequest)
		return
	}

	err := json.NewEncoder(res).Encode(task)

	if err != nil {
		fmt.Printf("ошибка при сериализации в json: %s\n", err.Error())
		http.Error(res, "serialization error", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func HendlerDeletTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(req, "id")

	_, ok := tasks[id]
	if !ok {
		fmt.Println("Элемента нет в маппе")
		http.Error(res, "task not found", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	res.WriteHeader(http.StatusOK)
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
