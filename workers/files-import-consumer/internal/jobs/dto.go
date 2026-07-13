package jobs

type ImportFailure struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

type ImportJobMessage struct {
	JobID     string `json:"job_id"`
	ObjectKey string `json:"key"` // MinIO object key
}
