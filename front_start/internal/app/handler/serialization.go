package handler

import (
    "front_start/internal/app/ds"
)

// Centralized JSON DTOs used by handler layer

// Requests
type DTO_registerRequest struct {
    Login    string `json:"login"`
    Password string `json:"password"`
}

type DTO_updateUserRequest struct {
    Password *string `json:"password"`
}

type DTO_serviceCreateRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    FullInfo    string `json:"full_info"`
    TheAxis     string `json:"the_axis"`
    Status      *bool  `json:"status"`
}

type DTO_taskUpdateRequest struct {
    TaskDescription string `json:"task_description"`
}

type DTO_resolveRequest struct {
    Action string `json:"action"` // "complete" | "reject"
}

type DTO_updateDegreesRequest struct {
    Degrees *float32 `json:"degrees"`
}

// Responses
type DTO_GatesListResponse struct {
    Items []ds.Gate `json:"items"`
}

type DTO_TasksListResponse struct {
    Items []ds.QuantumTask `json:"items"`
}

type DTO_SimpleIDResponse struct {
    ID int `json:"id"`
}

type DTO_TaskServiceLinkResponse struct {
    TaskID    uint `json:"task_id"`
    ServiceID int  `json:"service_id"`
}

type DTO_CurrTaskInfoResponse struct {
    TaskID        uint `json:"task_id"`
    ServicesCount int  `json:"services_count"`
}

type DTO_UploadImageResponse struct {
    ID    int    `json:"id"`
    Image string `json:"image"`
}

type DTO_UpdateDegreesResponse struct {
    TaskID    int      `json:"task_id"`
    ServiceID int      `json:"service_id"`
    Degrees   *float32 `json:"degrees"`
}

type DTO_RegisterResponse struct {
    Login string `json:"login"`
}

type DTO_LoginResponse struct {
    Login string `json:"login"`
}

type DTO_LogoutResponse struct {
    Logout bool `json:"logout"`
}


