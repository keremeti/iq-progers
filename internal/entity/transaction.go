package entity

import "strings"

type Transaction struct {
	Id      int32
	Date    int64
	Type    TransactionType
	UserId  int32
	Comment string
	Sum     float64
}

type TransactionType int

const (
	Recharge TransactionType = 0
	WriteOff TransactionType = 1
)

func (s TransactionType) ToString() string {
	switch s {
	case Recharge:
		return "RECHARGE"
	case WriteOff:
		return "WRITE_OFF"
	default:
		panic("неизвестный тип TransactionType")
	}
}

var (
	transactionTypeMap = map[string]TransactionType{
		"RECHARGE":  Recharge,
		"WRITE_OFF": WriteOff,
	}
)

func ParseStringToTransactionType(str string) (TransactionType, bool) {
	t, ok := transactionTypeMap[strings.ToUpper(str)]
	return t, ok
}
