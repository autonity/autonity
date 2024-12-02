package params

import "time"

const (
	// approximate time the node might need to start all services after shutting down
	// in case of a crash or restart
	BootingTime = time.Minute * 2
)
