package {{.Package}}

import (
	"context"
    "{{.Module}}/app/service/{{.ServerName}}/internal/biz"
    "{{.Module}}/app/service/{{.ServerName}}/internal/data/query"
	"{{.Module}}/app/service/{{.ServerName}}/internal/data"
    "{{.Module}}/app/service/{{.ServerName}}/internal/model"
    "github.com/go-kratos/kratos/v2/log"
)

var _ biz.{{.StructName}}Repo = (*{{.StructName}}Repo)(nil)

type {{.StructName}}Repo struct {
	data *data.Data
	log  *log.Helper
}

func New{{.StructName}}Repo(data *data.Data, logger log.Logger) biz.{{.StructName}}Repo {
	return &{{.StructName}}Repo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/{{.HumpPackageName}}")),
	}
}

// Create{{.StructName}} 写入{{.StructName}}
func ({{.Abbreviation}} *{{.StructName}}Repo) Create{{.StructName}}(ctx context.Context, m *model.{{.StructName}}) (int64, error) {
	err := query.Use(repo.data.Db.GetDebugDB()).{{.StructName}}.Create(m)
	if err != nil {
		repo.log.Error("Create{{.StructName}} Create: ", err)
		return 0, err
	}
	return 1, err
}

// Update{{.StructName}} 更新{{.StructName}}
func ({{.Abbreviation}} *{{.StructName}}Repo) Update{{.StructName}}(ctx context.Context, m *model.{{.StructName}}) (rows int64, error error) {
	u := query.Use(repo.data.Db.GetDebugDB()).{{.StructName}}
	res, err := u.Updates(m)
	if err != nil {
		repo.log.Error("UpdatePricingStrategy ", err)
		return 0, err
	}
	return res.RowsAffected, nil
}

// Delete{{.StructName}} 删除{{.StructName}}
func ({{.Abbreviation}} *{{.StructName}}Repo) Delete{{.StructName}}(ctx context.Context, id int64) (rows int64, error error) {
	u := query.Use(repo.data.Db.GetDebugDB()).{{.StructName}}
	res, err := u.Where(u.ID.Eq(id)).Delete()
	if err != nil {
		repo.log.Error("Delete{{.StructName}} ", err)
		return 0, err
	}
	return res.RowsAffected, nil
}


// List{{.StructName}} 分页获取{{.StructName}}记录
func ({{.Abbreviation}} *{{.StructName}}Repo) List{{.StructName}}(ctx context.Context, m *model.{{.StructName}}, pageSize,page int) (result interface{}, counts int64, err error) {
	limit := pageSize
	offset := pageSize * (page - 1)
	q := query.Use(repo.data.Db.GetDebugDB()).{{.StructName}}
	// 自定义查询条件
	result, counts, err = q.FindByPage(offset, limit)
	if err!=nil {
	    repo.log.Error("List{{.StructName}} ", err)
    	return
    }
	return
}

// Get{{.StructName}} 获取{{.StructName}}记录
func ({{.Abbreviation}} *{{.StructName}}Repo) Get{{.StructName}}(ctx context.Context, m *model.{{.StructName}}) (result interface{},err error) {
	u := query.Use(repo.data.Db.GetDebugDB()).{{.StructName}}
    res, err := u.Where(u.ID.Eq(m.ID)).Take()
    if err != nil {
    	repo.log.Error("Get{{.StructName}} ", err)
    	return nil, err
    }
    return res, nil
}