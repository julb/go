package monitoring

// A string enumeration to indicate the system health status.
type SystemStatus string

// The list of possible system health status.
// Up: everything is normal
// Partial: the system is working in degraded mode.
// Down: the system is not usable.
// OutOfService: the system is running but not in measure to deliver service.
// Unknown: the system status is unknown.
const (
	Up           SystemStatus = "UP"
	Partial      SystemStatus = "PARTIAL"
	Down         SystemStatus = "DOWN"
	OutOfService SystemStatus = "OUT_OF_SERVICE"
	Unknown      SystemStatus = "UNKNOWN"
)

var (
	systemStatus = Up
)

// Returns the actual system status.
func GetSystemStatus() SystemStatus {
	return systemStatus
}

// Returns the actual system status.
func setSystemStatus(newSystemStatus SystemStatus) {
	systemStatus = newSystemStatus
}
