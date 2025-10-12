package repository

import (
	"fmt"
	"math"
	"front_start/internal/app/ds"
)

// GetResByGate вычисляет новые коэффициенты на основе гейта и угла поворота
func (r *Repository) GetResByGate(gate ds.Gate, degree float32, koeff_0, koeff_1 *float32) error {
	// Конвертируем float32 в float64 для математических вычислений
	deg := float64(degree)

	k0 := float64(*koeff_0)
	k1 := float64(*koeff_1)

	var newK0, newK1 float64
	// Обработка гейта в зависимости от оси
	switch gate.TheAxis {	
	case "X":
		// X-вращение: R_x(θ) = [[cos(θ/2), -i*sin(θ/2)], [-i*sin(θ/2), cos(θ/2)]]
		// Для вещественных коэффициентов k0, k1:
		cosHalf := -math.Cos(deg / 2)
		sinHalf := -math.Sin(deg / 2)
		
		// Применяем матрицу вращения к вектору [k0, k1]
		// k0' = k0*cos(θ/2) - i*k1*sin(θ/2) -> действительная часть: k0*cos(θ/2)
		// k1' = -i*k0*sin(θ/2) + k1*cos(θ/2) -> действительная часть: k1*cos(θ/2)
		// Но для правильного результата нужно учесть, что мнимая часть "переходит" в вещественную
		if (k0 > k1) {
			newK0 = k0*cosHalf - k1*sinHalf  // k0*cos(θ/2) - k1*sin(θ/2)
			newK1 = 1 - newK0  				 
		} else {
			newK1 = k0*sinHalf + k1*cosHalf  // k0*sin(θ/2) + k1*cos(θ/2)
			newK0 = 1 - newK1  				
		}
		
	
	case "Y":
		// Y-вращение: R_y(θ) = [[cos(θ/2), -sin(θ/2)], [sin(θ/2), cos(θ/2)]]
		cosHalf := -math.Cos(deg / 2)
		sinHalf := -math.Sin(deg / 2)
		
		// Применяем матрицу вращения к вектору [k0, k1]
		if (k0 > k1) {
			newK0 = k0*cosHalf - k1*sinHalf  // k0*cos(θ/2) - k1*sin(θ/2)
			newK1 = 1 - newK0  
		} else {
			newK1 = k0*sinHalf + k1*cosHalf	// k0*sin(θ/2) + k1*cos(θ/2)
			newK0 = 1 - newK1
		}
	
	case "Z":
		// Z-вращение: R_z(θ) = [[e^(-iθ/2), 0], [0, e^(iθ/2)]]
		// Для вещественных коэффициентов k0, k1:
		cosHalf := -math.Cos(deg / 2)
		//sinHalf := math.Sin(deg / 2)
		
		// Применяем матрицу вращения к вектору [k0, k1]
		// k0' = k0 * e^(-iθ/2) = k0 * (cos(θ/2) - i*sin(θ/2))
		// k1' = k1 * e^(iθ/2) = k1 * (cos(θ/2) + i*sin(θ/2))
		// Для вещественных коэффициентов берем только действительные части
		if (k0 > k1) {
			newK0 = k0 * cosHalf  // k0*cos(θ/2)
			newK1 = 1 - newK0 
		} else {
			newK1 = k1 * cosHalf  // k1*cos(θ/2)
			newK0 = 1 - newK1  
		}
		
		
	default:
		// Для гейтов без оси (например, H, I) возвращаем исходные коэффициенты
		newK0 = k0
		newK1 = k1
	}

	// Обновляем коэффициенты
	*koeff_0 = float32(newK0)
	*koeff_1 = float32(newK1)
	
	return nil
}

// GetQTaskRes вычисляет результат квантовой задачи, применяя все гейты по очереди
func (r *Repository) GetQTaskRes(taskID uint) error {
	var task ds.QuantumTask
	
	// Загружаем задачу со всеми связанными гейтами и их углами
	err := r.db.
		Preload("GatesDegrees.Gate").
		First(&task, taskID).
		Error
	if err != nil {
		return fmt.Errorf("failed to load task: %v", err)
	}

	// Проверяем, что у задачи есть описание (необходимо для вычисления результата)
	if task.TaskDescription == "" {
		return fmt.Errorf("task description is required for result calculation")
	}

	// Инициализируем коэффициенты начальными значениями
	resKoeff0 := task.Res_koeff_0
	resKoeff1 := task.Res_koeff_1

	// Обходим все гейты в задаче и применяем их по очереди
	for _, gateDegree := range task.GatesDegrees {
		// Проверяем, что у гейта указан угол поворота
		if gateDegree.Degrees == nil {
			return fmt.Errorf("degrees not specified for gate %s (ID: %d)", 
				gateDegree.Gate.Title, gateDegree.Gate.ID_gate)
		}

		// Применяем гейт к текущим коэффициентам
		err := r.GetResByGate(gateDegree.Gate, *gateDegree.Degrees, &resKoeff0, &resKoeff1)
		if err != nil {
			return fmt.Errorf("failed to apply gate %s: %v", gateDegree.Gate.Title, err)
		}
	}

	// Обновляем коэффициенты в базе данных
	err = r.db.Model(&task).Updates(map[string]interface{}{
		"res_koeff_0": resKoeff0,
		"res_koeff_1": resKoeff1,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update task coefficients: %v", err)
	}

	return nil
}