package data

import (
	"autoCodeExamples/app/service/order/internal/biz"
	"autoCodeExamples/app/service/order/internal/data"
	"autoCodeExamples/app/service/order/internal/data/query"
	"autoCodeExamples/app/service/order/internal/model"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

var _ biz.OrderRepo = (*OrderRepo)(nil)

type OrderRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewOrderRepo(data *data.Data, logger log.Logger) biz.OrderRepo {
	return &OrderRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/order")),
	}
}

// CreateOrder 写入Order
func (repo *OrderRepo) CreateOrder(ctx context.Context, m *model.Order) (int64, error) {
	err := query.Use(repo.data.Db.GetDebugDB()).Order.Create(m)
	if err != nil {
		repo.log.Error("CreateOrder Create: ", err)
		return 0, err
	}
	return 1, err
}

// UpdateOrder 更新Order
func (repo *OrderRepo) UpdateOrder(ctx context.Context, m *model.Order) (rows int64, error error) {
	u := query.Use(repo.data.Db.GetDebugDB()).Order
	res, err := u.Updates(m)
	if err != nil {
		repo.log.Error("UpdatePricingStrategy ", err)
		return 0, err
	}
	return res.RowsAffected, nil
}

// DeleteOrder 删除Order
func (repo *OrderRepo) DeleteOrder(ctx context.Context, id int64) (rows int64, error error) {
	u := query.Use(repo.data.Db.GetDebugDB()).Order
	res, err := u.Where(u.ID.Eq(id)).Delete()
	if err != nil {
		repo.log.Error("DeleteOrder ", err)
		return 0, err
	}
	return res.RowsAffected, nil
}

// ListOrder 分页获取Order记录
func (repo *OrderRepo) ListOrder(ctx context.Context, m *model.Order, pageSize, page int) (result interface{}, counts int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	q := query.Use(repo.data.Db.GetDebugDB()).Order
	// 自定义查询条件
	result, counts, err = q.FindByPage(offset, limit)
	if err != nil {
		repo.log.Error("ListOrder ", err)
		return
	}
	return
}

// GetOrder 获取Order记录
func (repo *OrderRepo) GetOrder(ctx context.Context, m *model.Order) (result interface{}, err error) {
	u := query.Use(repo.data.Db.GetDebugDB()).Order
	res, err := u.Where(u.ID.Eq(m.ID)).Take()
	if err != nil {
		repo.log.Error("GetOrder ", err)
		return nil, err
	}
	return res, nil
}
