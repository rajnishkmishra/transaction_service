package vm

type Transaction struct {
	TransactionID uint64  `uri:"transaction_id" json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Type          string  `uri:"type" json:"type"`
	ParentID      uint64  `json:"parent_id"`
}

type TransactionRequest struct {
	Transaction
}

type DummyRequest struct{}

type CreateTransactionResponse struct {
	TransactionID uint64 `json:"transaction_id"`
}

type TransactionIdRequest struct {
	TransactionID uint64 `uri:"transaction_id"`
}

type TransactionTypeRequest struct {
	Type string `uri:"type"`
}

type GetSumResponse struct {
	Sum float64 `json:"sum"`
}
