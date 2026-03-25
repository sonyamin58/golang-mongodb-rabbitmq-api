package machinery

import (
	"fmt"

	"github.com/gravityblast/machinery/v2"
	"github.com/gravityblast/machinery/v2/tasks"
)

// AsyncClient wraps Machinery async client for sending tasks
type AsyncClient struct {
	server *machinery.Server
}

// NewAsyncClient creates a new async client with the given server
func NewAsyncClient(server *machinery.Server) *AsyncClient {
	return &AsyncClient{server: server}
}

// SendTopupTask sends a topup processing task asynchronously
func (c *AsyncClient) SendTopupTask(txID uint64) (*tasks.Signature, error) {
	sig := &tasks.Signature{
		Name: "process_topup",
		Args: []tasks.Arg{
			{Name: "tx_id", Type: "uint64", Value: txID},
		},
		RetryStrategy: &tasks.RetryStrategy{
			MaxRetries: 3,
			MinRetry:   5,
			MaxRetry:   30,
		},
	}

	asyncResult, err := c.server.SendTask(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to send topup task: %w", err)
	}

	// Wait for result (non-blocking)
	go func() {
		if err := asyncResult.Get(); err != nil {
			fmt.Printf("Topup task failed: %v\n", err)
		}
	}()

	return sig, nil
}

// SendWithdrawTask sends a withdrawal processing task asynchronously
func (c *AsyncClient) SendWithdrawTask(txID uint64) (*tasks.Signature, error) {
	sig := &tasks.Signature{
		Name: "process_withdraw",
		Args: []tasks.Arg{
			{Name: "tx_id", Type: "uint64", Value: txID},
		},
		RetryStrategy: &tasks.RetryStrategy{
			MaxRetries: 3,
			MinRetry:   5,
			MaxRetry:   30,
		},
	}

	asyncResult, err := c.server.SendTask(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to send withdraw task: %w", err)
	}

	go func() {
		if err := asyncResult.Get(); err != nil {
			fmt.Printf("Withdraw task failed: %v\n", err)
		}
	}()

	return sig, nil
}

// SendTransferTask sends a transfer processing task asynchronously
func (c *AsyncClient) SendTransferTask(txID uint64) (*tasks.Signature, error) {
	sig := &tasks.Signature{
		Name: "process_transfer",
		Args: []tasks.Arg{
			{Name: "tx_id", Type: "uint64", Value: txID},
		},
		RetryStrategy: &tasks.RetryStrategy{
			MaxRetries: 3,
			MinRetry:   5,
			MaxRetry:   30,
		},
	}

	asyncResult, err := c.server.SendTask(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to send transfer task: %w", err)
	}

	go func() {
		if err := asyncResult.Get(); err != nil {
			fmt.Printf("Transfer task failed: %v\n", err)
		}
	}()

	return sig, nil
}

// SendNotificationTask sends a notification task asynchronously
func (c *AsyncClient) SendNotificationTask(userID uint64, message string) (*tasks.Signature, error) {
	sig := &tasks.Signature{
		Name: "send_notification",
		Args: []tasks.Arg{
			{Name: "user_id", Type: "uint64", Value: userID},
			{Name: "message", Type: "string", Value: message},
		},
		RetryStrategy: &tasks.RetryStrategy{
			MaxRetries: 3,
			MinRetry:   2,
			MaxRetry:   10,
		},
	}

	asyncResult, err := c.server.SendTask(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification task: %w", err)
	}

	go func() {
		if err := asyncResult.Get(); err != nil {
			fmt.Printf("Notification task failed: %v\n", err)
		}
	}()

	return sig, nil
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
