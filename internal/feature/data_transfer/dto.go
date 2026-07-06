package data_transfer

type ImportServersResponse struct {
	SuccessCount int             `json:"success_count"`
	FailedCount  int             `json:"failed_count"`
	Failures     []ImportFailure `json:"failures"`
}

type ImportFailure struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

type ImportJobMessage struct {
	JobID     string `json:"job_id"`
	ObjectKey string `json:"key"` // MinIO object key
}
