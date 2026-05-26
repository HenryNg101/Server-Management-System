package server

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
