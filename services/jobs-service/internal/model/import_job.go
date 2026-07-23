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
	ID               string    `gorm:"column:id;type:uuid;primaryKey"`
	UserID           uint      `gorm:"uniqueIndex:idx_user_server_name"`
	FilePath         string    `gorm:"column:file_path;type:text;not null"`
	Status           JobStatus `gorm:"column:status;type:text;not null;check:status IN ('pending','processing','done','failed')"`
	ProcessedRows    int       `gorm:"column:processed_rows;type:int;not null;default:0;check:processed_rows >= 0"`
	SuccessRowsCount int       `gorm:"column:success_rows_count;type:int;not null;default:0;check:success_rows_count >= 0"`
	FailedRowsCount  int       `gorm:"column:failed_rows_count;type:int;not null;default:0;check:failed_rows_count >= 0"`
	ResultPath       *string   `gorm:"column:result_path;type:text"`
	Error            *string   `gorm:"column:error;type:text"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp;not null;default:now()"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp;not null;default:now()"`

	User User `gorm:"foreignKey:UserID"`
}
