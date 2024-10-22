package enum

type TransactionStatus string

const (
	TRANSACTION_STATUS_SUCCESSFUL TransactionStatus = "successful"
	TRANSACTION_STATUS_FAILED     TransactionStatus = "failed"
)

func GetTransactionStatus(t TransactionStatus) TransactionStatus {
	switch t {
	case TRANSACTION_STATUS_SUCCESSFUL:
		return TRANSACTION_STATUS_SUCCESSFUL
	default:
		return TRANSACTION_STATUS_FAILED
	}
}
