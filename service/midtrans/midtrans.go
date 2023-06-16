package midtrans

import (
	"capstone/model"
	"fmt"
	"os"
	"strconv"

	"github.com/midtrans/midtrans-go"
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
