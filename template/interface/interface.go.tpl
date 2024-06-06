
type {{.StructName}}Repo interface {
    // Create{{.StructName}} 写入{{.StructName}}
    Create{{.StructName}}(ctx context.Context, m *model.{{.StructName}}) (int64, error)
    // Update{{.StructName}} 更新{{.StructName}}
    Update{{.StructName}}(ctx context.Context, m *model.{{.StructName}}) (rows int64, error error)
    // Delete{{.StructName}} 删除{{.StructName}}
    Delete{{.StructName}}(ctx context.Context, id int64) (rows int64, error error)
    // Get{{.StructName}} 获取{{.StructName}}记录
    Get{{.StructName}}(ctx context.Context, m *model.{{.StructName}}) (result interface{},err error)
    // List{{.StructName}} 分页获取{{.StructName}}记录
    List{{.StructName}}(ctx context.Context, m *model.{{.StructName}}, pageSize,page int) (result interface{}, counts int64, err error)
}
