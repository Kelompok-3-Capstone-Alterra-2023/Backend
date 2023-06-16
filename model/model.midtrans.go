package model

type MidtransRequest struct {
	OrderNumber string
	Amount      int64
	Doctor      struct {
		ID       uint
		FullName string
		Price    int64
	}
	QTY        int32
	Method     string
	ServiceFee int64
	User       struct {
		FName string
		Email string
		Phone string
	}
}
