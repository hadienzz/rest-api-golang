package config

import (
	"github.com/midtrans/midtrans-go"
)

func InitMidtrans() {
	cfg := Get()

	midtrans.ServerKey = cfg.MidtransServerKey
	midtrans.Environment = midtrans.Sandbox
}
