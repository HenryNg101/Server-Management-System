package model

import "time"

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusDone       JobStatus = "done"
	JobStatusFailed     JobStatus = "failed"
)

type ImportJob struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey"`
	FilePath  string    `gorm:"column:file_path;type:text;not null"`
	Status    JobStatus `gorm:"column:status;type:text;not null;check:status IN ('pending','processing','done','failed')"`
	Progress  float64   `gorm:"column:progress;type:numeric(5,2);not null;default:0;check:progress >= 0 AND progress <= 100"`
	Error     *string   `gorm:"column:error;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;default:now()"`
}
