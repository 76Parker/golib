package health

import "context"

var (
	IsReady  bool = true
	NotReady bool = false
)

type SystemType string

var (
	KafkaSystemType    SystemType = "kafka"
	PostgresSystemType SystemType = "postgres"
	MinioSystemType    SystemType = "minio"
	TemporalSystemType SystemType = "temporal"
)

type Status struct {
	System   SystemType `json:"system"`
	IsReady  bool       `json:"is_ready"`
	ErrorMsg string     `json:"error_msg"`
}

type Checker interface {
	Check(ctx context.Context) Status
}
