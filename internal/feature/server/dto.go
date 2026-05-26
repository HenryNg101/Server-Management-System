package server

import "github.com/HenryNg101/server-management-system/internal/model"

type UpdateServerRequest struct {
	Name        *string `json:"name"`
	Status      *bool   `json:"status"`
	IPv4Address *string `json:"ipv4"`
	Port        *uint   `json:"port"`
	Protocol    *string `json:"protocol"`
}

type CreateServerRequest struct {
	Name        string `json:"name" binding:"required"`
	Status      bool   `json:"status" binding:"required"`
	IPv4Address string `json:"ipv4_address" binding:"required,ip"`
	Port        uint   `json:"port" binding:"required"`
	Protocol    string `json:"protocol"`
}

type ImportServersResponse struct {
	SuccessCount int             `json:"success_count"`
	FailedCount  int             `json:"failed_count"`
	Successes    []model.Server  `json:"successes"`
	Failures     []ImportFailure `json:"failures"`
}

type ImportFailure struct {
	Row    int               `json:"row"`
	Error  string            `json:"error"`
	Record map[string]string `json:"record"`
}
