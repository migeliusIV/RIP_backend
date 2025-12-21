package handler

import (
	"time"
)

//====== Requests ======

// DTO_Req_UserReg запрос регистрации пользователя
// @Description Данные для регистрации нового пользователя
type DTO_Req_UserReg struct {
	Login    string `json:"login" example:"quantum_user" binding:"required"`
	Password string `json:"password" example:"secure_password_123" binding:"required"`
}

// DTO_Req_UserUpd запрос обновления пользователя
// @Description Данные для обновления пароля пользователя
type DTO_Req_UserUpd struct {
	Password *string `json:"password" example:"new_secure_password"`
}

// DTO_Req_GateCreate запрос создания гейта
// @Description Данные для создания нового квантового гейта
type DTO_Req_GateCreate struct {
	Title        string   `json:"title" example:"Hadamard Gate" binding:"required"`
	Description  string   `json:"description" example:"Quantum Fourier transform gate"`
	FullInfo     string   `json:"full_info" example:"Creates superposition states"`
	TheAxis      string   `json:"the_axis" example:"X" enums:"X,Y,Z,non"`
	Status       *bool    `json:"status" example:"true"`
	I0j0         *int     `gorm:"column:i0_j0" json:"i0j0" example:"1"`
	I0j1         *int     `gorm:"column:i0_j1" json:"i0j1" example:"0"`
	I1j0         *int     `gorm:"column:i1_j0" json:"i1j0" example:"0"`
	I1j1         *int     `gorm:"column:i1_j1" json:"i1j1" example:"1"`
	Matrix_koeff *float32 `gorm:"column:matrix_koeff" json:"matrix_koeff" example:"0.707"`
}

// DTO_Req_TaskUpd запрос обновления задачи
// @Description Данные для обновления описания квантовой задачи
type DTO_Req_TaskUpd struct {
	TaskDescription string `json:"task_description" example:"Calculate quantum state probabilities" binding:"required"`
}

// DTO_Req_TaskResolve запрос решения задачи
// @Description Данные для завершения или отклонения квантовой задачи
type DTO_Req_TaskResolve struct {
	Action string `json:"action" example:"complete" enums:"complete,reject" binding:"required"`
}

// DTO_Req_DegreesUpd запрос обновления градусов
// @Description Данные для обновления углов поворота гейта
type DTO_Req_DegreesUpd struct {
	Degrees *float32 `json:"degrees" example:"45.5" binding:"required"`
}

//====== Responses ======

// DTO_Resp_Gate ответ с данными гейта
// @Description Полная информация о квантовом гейте
type DTO_Resp_Gate struct {
	ID_gate      uint     `gorm:"column:id_gate;primaryKey;autoIncrement" json:"ID_gate" example:"1"`
	Title        string   `gorm:"column:title;size:255;not null;default:gate-no-name;unique" json:"Title" example:"Pauli-X"`
	Description  string   `gorm:"column:description;not null" json:"Description" example:"Quantum NOT gate"`
	Status       bool     `gorm:"column:status; default:true" json:"Status" example:"true"`
	Image        *string  `gorm:"column:image" json:"Image" example:"gate_x.png"`
	I0j0         *int     `gorm:"column:i0_j0" json:"I0j0" example:"0"`
	I0j1         *int     `gorm:"column:i0_j1" json:"I0j1" example:"1"`
	I1j0         *int     `gorm:"column:i1_j0" json:"I1j0" example:"1"`
	I1j1         *int     `gorm:"column:i1_j1" json:"I1j1" example:"0"`
	Matrix_koeff *float32 `gorm:"column:matrix_koeff" json:"Matrix_koeff" example:"1.0"`
	FullInfo     string   `gorm:"column:full_info;not null" json:"Full_info" example:"Bit flip gate"`
	TheAxis      string   `gorm:"column:the_axis" json:"TheAxis" example:"X"`
}

// DTO_Resp_UploadImg ответ загрузки изображения
// @Description Результат загрузки изображения для гейта
type DTO_Resp_UploadImg struct {
	ID    int    `json:"id" example:"1"`
	Image string `json:"image" example:"gate_image.png"`
}

// DTO_Resp_CurrTaskInfo информация о текущей задаче
// @Description Статистика по текущей задаче пользователя
type DTO_Resp_CurrTaskInfo struct {
	TaskID        uint `json:"task_id" example:"5"`
	ServicesCount int  `json:"services_count" example:"3"`
}

// DTO_Resp_Tasks ответ с данными задачи
// @Description Полная информация о квантовой задаче включая гейты
type DTO_Resp_Tasks struct {
	ID_task         uint                    `json:"id_task" example:"1"`
	TaskStatus      string                  `json:"task_status" example:"completed" enums:"draft,formed,completed,rejected"`
	CreationDate    time.Time               `json:"creation_date" example:"2023-10-01T15:04:05Z"`
	ID_user         uint                    `json:"id_user" example:"1"`
	ConclusionDate  time.Time               `json:"conclusion_date" example:"2023-10-01T16:04:05Z"`
	FormedDate      time.Time               `json:"formed_date" example:"2023-10-01T16:04:05Z"`
	TaskDescription string                  `gorm:"column:task_description" json:"task_description" example:"Quantum state calculation"`
	Res_koeff_0     float32                 `gorm:"column:res_koeff_0;default:1" json:"res_koeff_0" example:"0.707"`
	Res_koeff_1     float32                 `gorm:"column:res_koeff_1;default:0" json:"res_koeff_1" example:"0.707"`
	GatesDegrees    []DTO_Resp_GatesDegrees `json:"gates_degrees"`
}

// DTO_Resp_GatesDegrees связь гейта и задачи с углами
// @Description Информация о гейте в задаче с указанными углами поворота
type DTO_Resp_GatesDegrees struct {
	Title   string   `json:"titile" example:"hi"`
	TheAxis string   `gorm:"column:the_axis" json:"TheAxis" example:"X"`
	Image   *string  `gorm:"column:image" json:"Image" example:"gate_x.png"`
	ID_gate uint     `json:"id_gate" example:"2"`
	ID_task uint     `json:"id_task" example:"1"`
	Degrees *float32 `json:"degrees" example:"90.0"`
}

// DTO_Resp_SimpleID простой ответ с ID
// @Ответ содержащий только идентификатор
type DTO_Resp_SimpleID struct {
	ID int `json:"id" example:"1"`
}

// DTO_Resp_UpdateDegrees ответ обновления градусов
// @Description Результат обновления углов поворота гейта
type DTO_Resp_UpdateDegrees struct {
	TaskID    int      `json:"task_id" example:"1"`
	ServiceID int      `json:"service_id" example:"2"`
	Degrees   *float32 `json:"degrees" example:"45.5"`
}

// DTO_Resp_TaskServiceLink связь задачи и сервиса
// @Description Информация о связи между задачей и гейтом
type DTO_Resp_TaskServiceLink struct {
	TaskID    uint `json:"task_id" example:"1"`
	ServiceID int  `json:"service_id" example:"3"`
}

// DTO_User данные пользователя
// @Description Базовая информация о пользователе системы
type DTO_User struct {
	ID_user uint   `json:"id_user" example:"1"`
	Login   string `json:"login" example:"quantum_researcher"`
	IsAdmin bool   `json:"is_admin" example:"false"`
}

// DTO_Resp_User ответ с пользователем
// @Description Упрощенные данные пользователя
type DTO_Resp_User struct {
	Login string `json:"login" example:"quantum_researcher"`
}

// DTO_Resp_TokenLogin ответ аутентификации
// @Description Ответ содержащий JWT токен и данные пользователя
type DTO_Resp_TokenLogin struct {
	Token string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  DTO_User `json:"user"`
}

// DTO_Resp_GatesDegrees связь гейта и задачи с углами
// @Description Информация о гейте в задаче с указанными углами поворота
type DTO_Resp_TasksGatesInfo struct {
	Title   string   `json:"titile" example:"hi"`
	TheAxis string   `gorm:"column:the_axis" json:"TheAxis" example:"X"`
	ID_task uint     `json:"id_task" example:"1"`
	Degrees *float32 `json:"degrees" example:"90.0"`
}
