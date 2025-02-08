package models

type TransactionStatus uint

const (
	IN_PROGRESS TransactionStatus = iota
	COMPLETED
	FAILED
)

type Transaction struct {
	ID       uint64            `json:"id" gorm:"primary_key;not null;auto_increment"`
	Amount   float64           `json:"amount"`
	Type     string            `json:"type"`
	Status   TransactionStatus `json:"status"`
	ParentID *uint64           `json:"parent_id,omitempty" gorm:"parent_id"`
	Parent   *Transaction      `json:"parent,omitempty" gorm:"ForeignKey:ParentID;AssociationForeignKey:ID"`
}
