package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Gate struct {
	ID          int
	Title       string
	Description string
	FullInfo    string
	Image       string
	IsEditable  bool
	TheAxis     string
}

func (r *Repository) GetGates() ([]Gate, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	gates := []Gate{ // массив элементов из наших структур
		{
			ID:          1,
			Title:       "Identity Gate",
			Description: "Не изменяет состояния кубита.",
			FullInfo:    "\tНичего не делает с состоянием кубита. Оставляет его без изменений.",
			Image:       "http://127.0.0.1:9000/ibm-pictures/img/I-gate.png",
			IsEditable:  false,
			TheAxis:     "",
		},
		{
			ID:          2,
			Title:       "Pauli-X Gate (NOT gate)",
			Description: "Инвертирует состояние кубита.",
			FullInfo:    "\tАналог классического NOT-гейта. Переворачивает состояние кубита.",
			Image:       "http://127.0.0.1:9000/ibm-pictures/img/X-gate.png",
			IsEditable:  false,
			TheAxis:     "",
		},
		{
			ID:          3,
			Title:       "X-axis Rotation Gate",
			Description: "Вращает кубит вокруг оси X на угол тэта.",
			FullInfo:    "\tЭта операция вращает состояние кубита на сфере Блоха вокруг оси X.\n\tЗначение угла поворота можно задать при компановке выражения (в деталях калькуляции).",
			Image:       "http://127.0.0.1:9000/ibm-pictures/img/X-rot-gate.png",
			IsEditable:  true,
			TheAxis:     "X",
		},
		{
			ID:          4,
			Title:       "Y-axis Rotation Gate",
			Description: "Вращает кубит вокруг оси Y на угол тэта.",
			FullInfo:    "\tЭта операция вращает состояние кубита на сфере Блоха вокруг оси Y.\n\tЗначение угла поворота можно задать при компановке выражения (в деталях калькуляции).",
			Image:       "http://127.0.0.1:9000/ibm-pictures/img/Y-rot-gate.png",
			IsEditable:  true,
			TheAxis:     "Y",
		},
		{
			ID:          5,
			Title:       "Z-axis Rotation Gate",
			Description: "Вращает кубит вокруг оси Z на угол тэта.",
			FullInfo:    "\tЭта операция вращает состояние кубита на сфере Блоха вокруг оси Z.\n\tЗначение угла поворота можно задать при компановке выражения (в деталях калькуляции).",
			Image:       "http://127.0.0.1:9000/ibm-pictures/img/Z-rot-gate.png",
			IsEditable:  true,
			TheAxis:     "Z",
		},
		{
			ID:          6,
			Title:       "H (Hadamard) Gate",
			Description: "Создает равномерную суперпозицию из базисного состояния.",
			FullInfo:    "\tОперация поворачивает кубит на 90 градусов вокруг оси Y, затем на 180 градусов вокруг оси X.\n\tЭто один из самых важных гейтов.",
			Image:       "http://127.0.0.1:9000/ibm-pictures/img/H-gate.png",
			IsEditable:  false,
			TheAxis:     "",
		},
	}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(gates) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return gates, nil
}

func (r *Repository) GetGate(id int) (Gate, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	gates, err := r.GetGates()
	if err != nil {
		return Gate{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, gate := range gates {
		if gate.ID == id {
			return gate, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Gate{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetGatesByTitle(title string) ([]Gate, error) {
	gates, err := r.GetGates()
	if err != nil {
		return []Gate{}, err
	}

	var result []Gate
	for _, gate := range gates {
		if strings.Contains(strings.ToLower(gate.Title), strings.ToLower(title)) {
			result = append(result, gate)
		}
	}

	return result, nil
}

type Task struct {
	ID_task int
	Gates   []Gate
}

func (r *Repository) GetTasks() ([]Task, error) {
	gates, err := r.GetGates()
	if err != nil {
		return nil, err
	}

	tasks := []Task{
		{
			ID_task: 1,
			Gates: []Gate{
				gates[0],
				gates[1],
				gates[3],
			},
		},
	}

	return tasks, nil
}

func (r *Repository) GetTask(id int) (Task, error) {
	tasks, err := r.GetTasks()
	if err != nil {
		return Task{}, err
	}

	for _, task := range tasks {
		if task.ID_task == id {
			return task, nil
		}
	}

	return Task{}, fmt.Errorf("фактор не найден")
}
