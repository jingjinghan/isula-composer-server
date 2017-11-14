package health

import (
	"fmt"
)

// Health is used for services to register and report their status
type Health interface {
	// GetStatus return the service status
	GetStatus() (string, string)
}

// RegisterHealth is used for services to regist themself
func RegisterHealth(name string, h Health) error {
	fmt.Println("rh")
	return nil
}
