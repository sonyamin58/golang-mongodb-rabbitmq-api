package machinery

import (
	"fmt"

	"github.com/gravityblast/machinery/v2"
	"github.com/gravityblast/machinery/v2/config"
)

// NewServer creates a new Machinery server with the given broker and result backend URLs
func NewServer(brokerURL, resultBackendURL string) (*machinery.Server, error) {
	cfg := &config.Config{
		Broker:                 brokerURL,
		DefaultQueue:           "machinery_tasks",
		ResultBackend:          resultBackendURL,
		ResultsExpireIn:        3600,
		WorkerConnectionRetries: 5,
		TaskTimeLimit:          1800,        // 30 minutes
		TaskHardTimeLimit:      2100,        // 35 minutes
		MaxWorkerRetries:       3,
	}

	server, err := machinery.NewServer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create machinery server: %w", err)
	}

	return server, nil
}

// RegisterTasks registers all task signatures with the server
func RegisterTasks(server *machinery.Server) error {
	// Register process_topup task
	if err := server.RegisterTask("process_topup", ProcessTopup); err != nil {
		return fmt.Errorf("failed to register process_topup task: %w", err)
	}

	// Register process_withdraw task
	if err := server.RegisterTask("process_withdraw", ProcessWithdraw); err != nil {
		return fmt.Errorf("failed to register process_withdraw task: %w", err)
	}

	// Register process_transfer task
	if err := server.RegisterTask("process_transfer", ProcessTransfer); err != nil {
		return fmt.Errorf("failed to register process_transfer task: %w", err)
	}

	// Register send_notification task
	if err := server.RegisterTask("send_notification", SendNotification); err != nil {
		return fmt.Errorf("failed to register send_notification task: %w", err)
	}

	return nil
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
