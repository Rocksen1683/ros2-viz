package ros2viz

import (
	"log"
	"os/exec"
)

// DataProvider defines the interface for getting ROS graph data.
type DataProvider interface {
	GetROSGraphData() ([]byte, error)
}

// ROSInspector implements the DataProvider interface.
type ROSInspector struct {
	scriptPath string
}

// NewROSInspector creates a new inspector.
func NewROSInspector(scriptPath string) *ROSInspector {
	return &ROSInspector{scriptPath: scriptPath}
}

// GetROSGraphData executes the Python introspection script.
func (i *ROSInspector) GetROSGraphData() ([]byte, error) {
	log.Println("Executing Python introspection script at:", i.scriptPath)
	cmd := exec.Command("python3", i.scriptPath)

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("Python script stderr: %s", string(exitError.Stderr))
		}
		return nil, err
	}
	return output, nil
}