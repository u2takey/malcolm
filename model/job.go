package model

import (
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
)

type Work batchv1.Job

// represent data model for job running in k8s
type Job struct {
	PipeLine []Work `json:"pipeline"`
}
