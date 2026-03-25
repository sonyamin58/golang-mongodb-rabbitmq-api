package machinery

import (
	"fmt"
	"log"

	"github.com/gravityblast/machinery/v2/tasks"
	"github.com/ibas/golib-api/internal/model"
	"github.com/ibas/golib-api/internal/repository"
	"gorm.io/gorm"
)

// TaskContext holds database connection for task execution
type TaskContext struct {
	DB *gorm.DB
}

// ProcessTopup handles topup (deposit) transaction processing
func ProcessTopup(txID uint64) error {
	ctx := GetTaskContext()
	if ctx == nil || ctx.DB == nil {
		return fmt.Errorf("database connection not available")
	}

	log.Printf("[TASK] Processing topup for transaction ID: %d", txID)

	// Load transaction
	txRepo := repository.NewTransactionRepository(ctx.DB)
	transaction, err := txRepo.FindByID(uint(txID))
	if err != nil {
		log.Printf("[TASK] Failed to load transaction: %v", err)
		return err
	}
	if transaction == nil {
		err := fmt.Errorf("transaction not found: %d", txID)
		log.Printf("[TASK] %v", err)
		return err
	}

	// Load account
	accountRepo := repository.NewAccountRepository(ctx.DB)
	account, err := accountRepo.FindByID(transaction.AccountID)
	if err != nil {
		log.Printf("[TASK] Failed to load account: %v", err)
		return err
	}
	if account == nil {
		err := fmt.Errorf("account not found for transaction")
		log.Printf("[TASK] %v", err)
		return err
	}

	// Update account balance
	account.Balance += transaction.Amount
	if err := accountRepo.Update(account); err != nil {
		log.Printf("[TASK] Failed to update account balance: %v", err)
		return err
	}

	// Update transaction status to completed
	transaction.Status = model.TransactionStatusCompleted
	if err := txRepo.Update(transaction); err != nil {
		log.Printf("[TASK] Failed to update transaction status: %v", err)
		return err
	}

	log.Printf("[TASK] Topup completed successfully for transaction ID: %d, new balance: %.2f", txID, account.Balance)
	return nil
}

// ProcessWithdraw handles withdrawal transaction processing
func ProcessWithdraw(txID uint64) error {
	ctx := GetTaskContext()
	if ctx == nil || ctx.DB == nil {
		return fmt.Errorf("database connection not available")
	}

	log.Printf("[TASK] Processing withdrawal for transaction ID: %d", txID)

	// Load transaction
	txRepo := repository.NewTransactionRepository(ctx.DB)
	transaction, err := txRepo.FindByID(uint(txID))
	if err != nil {
		log.Printf("[TASK] Failed to load transaction: %v", err)
		return err
	}
	if transaction == nil {
		err := fmt.Errorf("transaction not found: %d", txID)
		log.Printf("[TASK] %v", err)
		return err
	}

	// Load account
	accountRepo := repository.NewAccountRepository(ctx.DB)
	account, err := accountRepo.FindByID(transaction.AccountID)
	if err != nil {
		log.Printf("[TASK] Failed to load account: %v", err)
		return err
	}
	if account == nil {
		err := fmt.Errorf("account not found for transaction")
		log.Printf("[TASK] %v", err)
		return err
	}

	// Check sufficient balance
	if account.Balance < transaction.Amount {
		err := fmt.Errorf("insufficient balance: available %.2f, required %.2f", account.Balance, transaction.Amount)
		log.Printf("[TASK] %v", err)
		// Update transaction status to failed
		transaction.Status = model.TransactionStatusFailed
		txRepo.Update(transaction)
		return err
	}

	// Update account balance
	account.Balance -= transaction.Amount
	if err := accountRepo.Update(account); err != nil {
		log.Printf("[TASK] Failed to update account balance: %v", err)
		return err
	}

	// Update transaction status to completed
	transaction.Status = model.TransactionStatusCompleted
	if err := txRepo.Update(transaction); err != nil {
		log.Printf("[TASK] Failed to update transaction status: %v", err)
		return err
	}

	log.Printf("[TASK] Withdrawal completed successfully for transaction ID: %d, new balance: %.2f", txID, account.Balance)
	return nil
}

// ProcessTransfer handles transfer transaction processing
func ProcessTransfer(txID uint64) error {
	ctx := GetTaskContext()
	if ctx == nil || ctx.DB == nil {
		return fmt.Errorf("database connection not available")
	}

	log.Printf("[TASK] Processing transfer for transaction ID: %d", txID)

	// Load transaction
	txRepo := repository.NewTransactionRepository(ctx.DB)
	transaction, err := txRepo.FindByID(uint(txID))
	if err != nil {
		log.Printf("[TASK] Failed to load transaction: %v", err)
		return err
	}
	if transaction == nil {
		err := fmt.Errorf("transaction not found: %d", txID)
		log.Printf("[TASK] %v", err)
		return err
	}

	// Load source account
	accountRepo := repository.NewAccountRepository(ctx.DB)
	fromAccount, err := accountRepo.FindByID(transaction.AccountID)
	if err != nil {
		log.Printf("[TASK] Failed to load source account: %v", err)
		return err
	}
	if fromAccount == nil {
		err := fmt.Errorf("source account not found")
		log.Printf("[TASK] %v", err)
		return err
	}

	// Check sufficient balance
	if fromAccount.Balance < transaction.Amount {
		err := fmt.Errorf("insufficient balance in source account")
		log.Printf("[TASK] %v", err)
		transaction.Status = model.TransactionStatusFailed
		txRepo.Update(transaction)
		return err
	}

	// Deduct from source account
	fromAccount.Balance -= transaction.Amount
	if err := accountRepo.Update(fromAccount); err != nil {
		log.Printf("[TASK] Failed to update source account: %v", err)
		return err
	}

	// If there's a related account, add to destination
	if transaction.RelatedAccountID != nil {
		toAccount, err := accountRepo.FindByID(*transaction.RelatedAccountID)
		if err != nil {
			log.Printf("[TASK] Failed to load destination account: %v", err)
			// Rollback: restore source account
			fromAccount.Balance += transaction.Amount
			accountRepo.Update(fromAccount)
			return err
		}
		if toAccount != nil {
			toAccount.Balance += transaction.Amount
			if err := accountRepo.Update(toAccount); err != nil {
				log.Printf("[TASK] Failed to update destination account: %v", err)
				// Rollback: restore source account
				fromAccount.Balance += transaction.Amount
				accountRepo.Update(fromAccount)
				return err
			}
		}
	}

	// Update transaction status to completed
	transaction.Status = model.TransactionStatusCompleted
	if err := txRepo.Update(transaction); err != nil {
		log.Printf("[TASK] Failed to update transaction status: %v", err)
		return err
	}

	log.Printf("[TASK] Transfer completed successfully for transaction ID: %d", txID)
	return nil
}

// SendNotification sends notification to user
func SendNotification(userID uint64, message string) error {
	log.Printf("[TASK] Sending notification to user ID: %d, message: %s", userID, message)

	// In production, this would:
	// 1. Store notification in database
	// 2. Send push notification
	// 3. Send email
	// 4. Send SMS

	// For now, we just log and return success
	log.Printf("[TASK] Notification sent to user ID: %d", userID)
	return nil
}

// Task signatures for reference
var (
	ProcessTopupSignature = &tasks.Signature{
		Name: "process_topup",
		Args: []tasks.Arg{
			{Name: "tx_id", Type: "uint64", Value: 0},
		},
	}
	ProcessWithdrawSignature = &tasks.Signature{
		Name: "process_withdraw",
		Args: []tasks.Arg{
			{Name: "tx_id", Type: "uint64", Value: 0},
		},
	}
	ProcessTransferSignature = &tasks.Signature{
		Name: "process_transfer",
		Args: []tasks.Arg{
			{Name: "tx_id", Type: "uint64", Value: 0},
		},
	}
	SendNotificationSignature = &tasks.Signature{
		Name: "send_notification",
		Args: []tasks.Arg{
			{Name: "user_id", Type: "uint64", Value: 0},
			{Name: "message", Type: "string", Value: ""},
		},
	}
)
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc

// Global task context for worker initialization
var globalTaskContext *TaskContext

// SetTaskContext sets the global task context for worker
func SetTaskContext(ctx *TaskContext) {
	globalTaskContext = ctx
}

// GetTaskContext gets the global task context
func GetTaskContext() *TaskContext {
	return globalTaskContext
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
