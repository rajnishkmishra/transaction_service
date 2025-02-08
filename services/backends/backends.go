package backends

import (
	"bitbucket.org/transaction_service/migration"
	"bitbucket.org/transaction_service/services/apis"
	"bitbucket.org/transaction_service/services/transaction_service"
	"bitbucket.org/transaction_service/utils"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func PathHandler(backends utils.Backends) {
	migrationService := migration.NewMigrationService(backends.MySQLConn)
	go migrationService.InitMigration()
	r = backends.GinEngine

	trasactionService := transaction_service.NewTransactionService(backends.MySQLConn.DB)
	apis.NewTransactionController(r, trasactionService)
}
