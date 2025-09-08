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
}

func (r *Repository) GetGates() ([]Gate, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	gates := []Gate{ // массив элементов из наших структур
		{
			ID:          1,
			Title:       "Identity Gate",
			Description: "Ничего не делает с состоянием кубита. Оставляет его без изменений.",
		},
		{
			ID:          2,
			Title:       "Pauli-X Gate (NOT gate)",
			Description: "Инвертирует состояние кубита.",
		},
		{
			ID:          3,
			Title:       "X-axis Rotation Gate",
			Description: "Вращает кубит вокруг оси X на угол тэта.",
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
