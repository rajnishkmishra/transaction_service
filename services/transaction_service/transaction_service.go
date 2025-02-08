package transaction_service

import (
	"errors"
	"net/http"

	"bitbucket.org/transaction_service/models"
	"bitbucket.org/transaction_service/models/vm"
	"bitbucket.org/transaction_service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		db: db,
	}
}

func (t *TransactionService) CreateTransaction(ctx *utils.Context, request vm.DummyRequest) (response vm.CreateTransactionResponse, werr utils.WrapperError) {
	dbTransaction := models.Transaction{}
	err := t.db.Save(&dbTransaction).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response.TransactionID = dbTransaction.ID
	return
}

func (t *TransactionService) MakePayment(ctx *utils.Context, request vm.TransactionRequest) (response string, werr utils.WrapperError) {
	if request.TransactionID == 0 {
		err := errors.New("invalid transaction id")
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}
	if request.Amount <= 0 {
		err := errors.New("invalid amount")
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}
	if request.Type == "" {
		err := errors.New("invalid type")
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}
	dbTransaction := models.Transaction{}
	err := t.db.Where("id = ?", request.TransactionID).First(&dbTransaction).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	dbTransaction.Amount = request.Amount
	dbTransaction.Type = request.Type
	dbTransaction.Status = models.COMPLETED
	if request.ParentID > 0 {
		dbTransaction.ParentID = &request.ParentID
	}
	err = t.db.Where("id = ?", request.TransactionID).Save(&dbTransaction).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response = "Payment done successfully!"
	return
}

func (t *TransactionService) GetTransactionByID(ctx *utils.Context, request vm.TransactionIdRequest) (response vm.Transaction, werr utils.WrapperError) {
	if request.TransactionID == 0 {
		err := errors.New("invalid transaction id")
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}
	dbTransaction := models.Transaction{}
	err := t.db.Where("id = ? and status = 1", request.TransactionID).First(&dbTransaction).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response.TransactionID = dbTransaction.ID
	response.Amount = dbTransaction.Amount
	response.Type = dbTransaction.Type
	if dbTransaction.ParentID != nil {
		response.ParentID = *dbTransaction.ParentID
	}
	return
}

func (t *TransactionService) GetTransactionByType(ctx *utils.Context, request vm.TransactionTypeRequest) (response []uint64, werr utils.WrapperError) {
	if request.Type == "" {
		err := errors.New("invalid type")
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}
	dbTransactions := []models.Transaction{}
	err := t.db.Where("type = ? and status = 1", request.Type).Find(&dbTransactions).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response = make([]uint64, 0)
	for _, dbTransaction := range dbTransactions {
		response = append(response, dbTransaction.ID)
	}
	return
}

func (t *TransactionService) GetSum(ctx *utils.Context, request vm.TransactionIdRequest) (response vm.GetSumResponse, werr utils.WrapperError) {
	if request.TransactionID == 0 {
		err := errors.New("invalid transaction id")
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}
	transactionAmount := float64(0)
	err := t.db.Table("transactions").Select("amount").Where("id = ?", request.TransactionID).Find(&transactionAmount).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}

	totalAmount, err := t.sumOfAmountOfChildren(ctx, request.TransactionID)
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response.Sum = totalAmount + transactionAmount
	return
}

func (t *TransactionService) sumOfAmountOfChildren(ctx *utils.Context, transactionID uint64) (float64, error) {
	totalAmount := float64(0)
	dbTransactions := []models.Transaction{}
	err := t.db.Table("transactions").Where("parent_id = ?", transactionID).Find(&dbTransactions).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		return totalAmount, err
	}

	for _, dbTransaction := range dbTransactions {
		totalAmount += dbTransaction.Amount
		childSum, err := t.sumOfAmountOfChildren(ctx, dbTransaction.ID)
		if err != nil {
			logrus.WithContext(ctx.Ctx).Error(err)
			return totalAmount, err
		}
		totalAmount += childSum
	}
	return totalAmount, err
}
