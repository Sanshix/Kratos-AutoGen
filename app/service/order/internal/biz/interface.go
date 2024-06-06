package biz

import (
	"context"
	"github.com/Sanshix/Kratos-AutoGen/app/service/order/internal/model"
)

type OrderRepo interface {
	// CreateOrder 写入Order
	CreateOrder(ctx context.Context, m *model.Order) (int64, error)
	// UpdateOrder 更新Order
	UpdateOrder(ctx context.Context, m *model.Order) (rows int64, error error)
	// DeleteOrder 删除Order
	DeleteOrder(ctx context.Context, id int64) (rows int64, error error)
	// GetOrder 获取Order记录
	GetOrder(ctx context.Context, m *model.Order) (result interface{}, err error)
	// ListOrder 分页获取Order记录
	ListOrder(ctx context.Context, m *model.Order, pageSize, page int) (result interface{}, counts int64, err error)
}
