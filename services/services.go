package services

import (
	"github.com/CPU-commits/Intranet_BAnnoucements/models"
	"github.com/CPU-commits/Intranet_BAnnoucements/stack"
)

// Models
var (
	annoucementModel = models.NewAnnoucementModel()
)

var nats = stack.NewNats()

// Error Response
type ErrorRes struct {
	Err        error
	StatusCode int
}
