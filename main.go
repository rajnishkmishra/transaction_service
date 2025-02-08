package main

import (
	"bitbucket.org/transaction_service/services/backends"
	"bitbucket.org/transaction_service/utils"
)

func main() {
	utils.SetupAndRun(backends.PathHandler)
}
