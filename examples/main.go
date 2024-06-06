package main

import (
	autoGen "github.com/Sanshix/Kratos-AutoGen"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const (
	Module     = "autoGenExamples"
	ServerName = "order"
)

// generate code
func main() {
	g := gen.NewGenerator(gen.Config{
		// 自动构建curd代码的输出目录
		OutPath: "./app/service/" + ServerName + "/internal/data/query",
		// 自动构建model代码的输出目录
		ModelPkgPath: "./app/service/" + ServerName + "/internal/model",
		// 生成 model 时，是否生成 gorm tag
		FieldWithTypeTag: true,
		// without context
		Mode: gen.WithoutContext,
	})
	// 链接数据库
	dsn := "dbAccount:dbPassword@tcp(127.0.0.1:3306)/dbName?parseTime=true&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(db)

	// 自定义字段的数据类型,统一数字类型为int64,兼容protobuf
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"smallint": func(columnType gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"mediumint": func(columnType gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"int": func(columnType gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"bigint": func(columnType gorm.ColumnType) (dataType string) {
			return "int64"
		},
	}
	g.WithDataTypeMap(dataMap)

	// 设置生成model的同时构建gen封装代码
	g.ApplyBasic(
		g.GenerateModel("orders"),
	)
	// 执行生成model代码
	g.Execute()

	// 定义配置信息，指定文件目录和文件名
	autoArg := autoGen.AutoCodeStruct{
		Abbreviation:     "repo",
		AutoMoveFile:     true,
		AutoMoveFilePath: "./app/service/" + ServerName + "/internal/data/",
		Fields:           nil,
		Package:          "data",
		Module:           Module,
		ServerName:       ServerName,
		InterfacePath:    "./app/service/" + ServerName + "/internal/biz/interface.go",
	}
	// 自动生成repo层的CRUD实现和更新interface接口定义
	// 根据order表自动生成CRUD代码
	err := autoGen.AutoGenerateCRUD("Order", "order", autoArg)
	if err != nil {
		panic(err)
	}
}
