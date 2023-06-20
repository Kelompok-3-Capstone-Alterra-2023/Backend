package midtrans

import (
	"capstone/model"
	"fmt"
	"os"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/iris"
	"github.com/midtrans/midtrans-go/snap"
)

func CreateSnapToken(request *model.MidtransRequest) (*snap.Response, error) {
	s := snap.Client{}
	s.New(os.Getenv("MT_SERVER_KEY"), midtrans.Sandbox)
	itemName := fmt.Sprintf("Sesi Konsultasi via %s", request.Method)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderNumber,
			GrossAmt: request.Amount,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    strconv.Itoa(int(request.Doctor.ID)),
				Name:  itemName,
				Price: int64(request.Doctor.Price),
				Qty:   1,
			},
			{
				Name:  "Service Fee",
				Price: request.ServiceFee,
				Qty:   1,
			},
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.User.FName,
			Email: request.User.Email,
			Phone: request.User.Phone,
		},
	}

	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		return &snap.Response{}, err
	}

	return snapResp, nil
}

func Payout(data *model.Withdraw) (*iris.CreatePayoutResponse, *midtrans.Error) {
	amount := fmt.Sprintf("%f", data.Total)
	req := iris.CreatePayoutReq{
		Payouts: []iris.CreatePayoutDetailReq{
			{
				BeneficiaryName:    data.AccountName,
				BeneficiaryAccount: data.AccountNumber,
				BeneficiaryBank:    data.Bank,
				BeneficiaryEmail:   data.Doctor.Email,
				Amount:             amount,
				Notes:              data.Notes,
			},
		},
	}

	var i iris.Client
	i.New(os.Getenv("MT_IRIS-API-KEY"), midtrans.Sandbox)

	irisResponse, err := i.CreatePayout(req)
	if err != nil {
		return nil, err
	}

	return irisResponse, nil
}

func ApprovePayout(ReferenceNumber string) error {
	var i iris.Client
	i.New(os.Getenv("MT_IRIS-API-KEY"), midtrans.Sandbox)
	req := iris.ApprovePayoutReq{
		ReferenceNo: []string{ReferenceNumber},
	}
	_, err := i.ApprovePayout(req)
	if err != nil {
		return err
	}

	return nil
}
