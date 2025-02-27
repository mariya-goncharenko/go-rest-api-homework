package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

/*
Обработчик для получения всех задач

Обработчик должен вернуть все задачи, которые хранятся в мапе.

Конечная точка /tasks.

Метод GET.

При успешном запросе сервер должен вернуть статус 200 OK.

При ошибке сервер должен вернуть статус 500 Internal Server Error.
*/
func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var taskList []Task
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	if err := json.NewEncoder(w).Encode(taskList); err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

/*
Обработчик для получения задачи по ID

Обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе.

В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.

Конечная точка /tasks/{id}.

Метод GET.

При успешном выполнении запроса сервер должен вернуть статус 200 OK.

В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
*/
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	task, found := tasks[id]
	if !found {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

/*
Обработчик для отправки задачи на сервер

Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.

Конечная точка /tasks.

Метод POST.

При успешном запросе сервер должен вернуть статус 201 Created.

При ошибке сервер должен вернуть статус 400 Bad Request.
*/
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, что ID уникален
	if _, exists := tasks[newTask.ID]; exists {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

/*
Обработчик удаления задачи по ID

Обработчик должен удалить задачу из мапы по её ID. Здесь так же нужно сначала проверить, есть ли задача с таким ID в мапе, если нет вернуть соответствующий статус.

Конечная точка /tasks/{id}.

Метод DELETE.

При успешном выполнении запроса сервер должен вернуть статус 200 OK.

В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
*/
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	if _, found := tasks[id]; !found {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasksHandler) // Получить все задачи

	r.Get("/tasks/{id}", getTaskHandler) // Получить задачу по ID

	r.Post("/tasks", createTaskHandler) // Создать новую задачу

	r.Delete("/tasks/{id}", deleteTaskHandler) // Удалить задачу по ID

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
