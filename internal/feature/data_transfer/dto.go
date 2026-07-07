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

type CreateImportJobResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

type ImportJobMessage struct {
	JobID     string `json:"job_id"`
	ObjectKey string `json:"key"` // MinIO object key
}

type GetJobResponse struct {
	ID               string  `json:"id"`
	Status           string  `json:"status"`
	ProcessedRows    int     `json:"processed_rows"`
	SuccessRowsCount int     `json:"success_rows_count"`
	FailedRowsCount  int     `json:"failed_rows_count"`
	Error            *string `json:"error,omitempty"`

	InputFileURL    *string `json:"input_file_url,omitempty"`
	FailuresFileURL *string `json:"failures_file_url,omitempty"`
}
