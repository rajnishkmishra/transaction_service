package apis

import (
	"bitbucket.org/transaction_service/services/transaction_service"
	"bitbucket.org/transaction_service/utils"
	"github.com/gin-gonic/gin"
)

const (
	CreateTransaction    = "/transaction"
	MakePayment          = "/transaction/:transaction_id"
	GetTransactionByID   = "/transaction/:transaction_id"
	GetTransactionByType = "/types/:type"
	GetSum               = "/sum/:transaction_id"
)

func NewTransactionController(engine *gin.Engine, transactionService *transaction_service.TransactionService) {
	router := engine.Group("/transactionservice")
	router.POST(CreateTransaction, utils.Controller(utils.NewOptions(transactionService.CreateTransaction).ForPost()))
	router.PUT(MakePayment, utils.Controller(utils.NewOptions(transactionService.MakePayment).ForPost()))
	router.GET(GetTransactionByID, utils.Controller(utils.NewOptions(transactionService.GetTransactionByID)))
	router.GET(GetTransactionByType, utils.Controller(utils.NewOptions(transactionService.GetTransactionByType)))
	router.GET(GetSum, utils.Controller(utils.NewOptions(transactionService.GetSum)))
}
