package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type LuminexServiceHandler struct {
	log *log.Helper
}

func NewLuminexServiceHandler(logger log.Logger) *LuminexServiceHandler {
	return &LuminexServiceHandler{
		log: log.NewHelper(logger),
	}
}
