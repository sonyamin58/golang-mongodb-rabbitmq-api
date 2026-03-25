package celery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type Publisher struct {
	redisClient *redis.Client
	taskPrefix  string
}

type TaskMessage struct {
	TaskName string      `json:"task"`
	Args     interface{} `json:"args,omitempty"`
	Kwargs   interface{} `json:"kwargs,omitempty"`
	ID       string      `json:"id"`
}

func NewPublisher(redisClient *redis.Client) *Publisher {
	return &Publisher{
		redisClient: redisClient,
		taskPrefix:  "celery",
	}
}

func (p *Publisher) PublishTask(taskName string, args interface{}, kwargs interface{}) error {
	msg := TaskMessage{
		TaskName: taskName,
		Args:     args,
		Kwargs:   kwargs,
		ID:       generateTaskID(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal task message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queue := fmt.Sprintf("%s:default", p.taskPrefix)
	if err := p.redisClient.LPush(ctx, queue, data).Err(); err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	return nil
}

func (p *Publisher) PublishTransactionNotification(transactionID uint) error {
	args := map[string]interface{}{
		"transaction_id": transactionID,
	}
	return p.PublishTask("tasks.process_transaction", args, nil)
}

func (p *Publisher) PublishEmailNotification(userID uint, subject, body string) error {
	args := map[string]interface{}{
		"user_id":  userID,
		"subject": subject,
		"body":    body,
	}
	return p.PublishTask("tasks.send_email", args, nil)
}

func (p *Publisher) PublishAccountStatement(accountID uint, startDate, endDate time.Time) error {
	args := map[string]interface{}{
		"account_id": accountID,
		"start_date": startDate.Format(time.RFC3339),
		"end_date":   endDate.Format(time.RFC3339),
	}
	return p.PublishTask("tasks.generate_statement", args, nil)
}

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
