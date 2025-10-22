package handler

import (
    "time"
)

// Centralized JSON DTOs used by handler layer

//====== Requests ======
// Users
type DTO_Req_UserReg struct {
    Login    string `json:"login"`
    Password string `json:"password"`
}

type DTO_Req_UserUpd struct {
    Password *string `json:"password"`
}

// Gates
type DTO_Req_GateCreate struct { // исправить
    Title       string `json:"title"`
    Description string `json:"description"`
    FullInfo    string `json:"full_info"`
    TheAxis     string `json:"the_axis"`
    Status      *bool  `json:"status"`
    I0j0 		*int 	`gorm:"column:i0_j0"`
	I0j1 		*int 	`gorm:"column:i0_j1"`
	I1j0 		*int 	`gorm:"column:i1_j0"`
	I1j1 		*int 	`gorm:"column:i1_j1"`
	Matrix_koeff *float32 `gorm:"column:matrix_koeff"`
}

// Tasks
type DTO_Req_TaskUpd struct {
    TaskDescription string `json:"task_description"`
}

type DTO_Req_TaskResolve struct {
    Action string `json:"action"` // "complete" | "reject"
}

// DegreesToGates
type DTO_Req_DegreesUpd struct {
    Degrees *float32 `json:"degrees"`
}

//====== Responses ======
// Gates
type DTO_Resp_Gate struct {
    ID_gate     uint    `gorm:"column:id_gate;primaryKey;autoIncrement"`
	Title       string  `gorm:"column:title;size:255;not null;default:gate-no-name;unique"`
	Description string  `gorm:"column:description;not null"`
	Status      bool    `gorm:"column:status; default:true"`
	Image       *string `gorm:"column:image"`
	I0j0 		*int 	`gorm:"column:i0_j0"`
	I0j1 		*int 	`gorm:"column:i0_j1"`
	I1j0 		*int 	`gorm:"column:i1_j0"`
	I1j1 		*int 	`gorm:"column:i1_j1"`
	Matrix_koeff *float32 `gorm:"column:matrix_koeff"`
	// subject area
	FullInfo string `gorm:"column:full_info;not null"`
	TheAxis  string `gorm:"column:the_axis"`
}

type DTO_Resp_UploadImg struct {
    ID    int    `json:"id"`
    Image string `json:"image"`
}

type DTO_Resp_CurrTaskInfo struct {
    TaskID        uint `json:"task_id"`
    ServicesCount int  `json:"services_count"`
}
// Tasks
type DTO_Resp_Tasks struct {
    ID_task         uint `json:"id_task"`
    TaskStatus      string `json:"task_status"`
    CreationDate    time.Time `json:"creation_date"`
    ID_user         uint `json:"id_user"`
    ConclusionDate  time.Time `json:"conclusion_date"`
    TaskDescription string  `gorm:"column:task_description"`
	Res_koeff_0     float32 `gorm:"column:res_koeff_0;default:1"`
	Res_koeff_1     float32 `gorm:"column:res_koeff_1;default:0"`
    GatesDegrees    []DTO_Resp_GatesDegrees `json:"gates_degrees"`
}

type DTO_Resp_GatesDegrees struct {
    ID_gate         uint `json:"id_gate"`
    ID_task         uint `json:"id_task"`
    Degrees         *float32 `json:"degrees"`
}

type DTO_Resp_SimpleID struct {
    ID int `json:"id"`
}

// DegreesToGates
type DTO_Resp_UpdateDegrees struct {
    TaskID    int      `json:"task_id"`
    ServiceID int      `json:"service_id"`
    Degrees   *float32 `json:"degrees"`
}

type DTO_Resp_TaskServiceLink struct {
    TaskID    uint `json:"task_id"`
    ServiceID int  `json:"service_id"`
}

// Users
type DTO_User struct {
    ID_user     uint `json:"id_user"`
    Login       string `json:"login"`
    IsAdmin     bool `json:"is_admin"`
}

type DTO_Resp_User struct {
    Login string `json:"login"`
}

type DTO_Resp_TokenLogin struct {
    Token string `json:"token"`
    User  DTO_User `json:"user"`
}

/* trash
type DTO_Resp_UserLogout struct {
    Logout bool `json:"logout"`
}
*/