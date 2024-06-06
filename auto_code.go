package autoGen

import (
	"bytes"
	"fmt"
	"github.com/Sanshix/Kratos-AutoGen/file"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	autoPath      = "autocode_template/"
	basePath      = "./template/doa"
	interfacePath = "./template/interface"
)

type AutoCodeStruct struct {
	StructName       string   `json:"structName"`       // Struct名称
	TableName        string   `json:"tableName"`        // 表名
	HumpPackageName  string   `json:"humpPackageName"`  // 驼峰文件名称
	Abbreviation     string   `json:"abbreviation"`     // Struct简称
	AutoMoveFile     bool     `json:"autoMoveFile"`     // 是否自动移动文件
	AutoMoveFilePath string   `json:"autoMoveFilePath"` // 自动移动文件路径
	Fields           []*Field `json:"fields"`
	Package          string   `json:"package"`       // 文件包名
	Module           string   `json:"module"`        // 项目模块名
	ServerName       string   `json:"serverName"`    // 服务名
	InterfacePath    string   `json:"interfacePath"` //	interface文件目录
}

type Field struct {
	FieldName string `json:"fieldName"` // Field名
	FieldDesc string `json:"fieldDesc"` // 中文名

	FieldType       string `json:"fieldType"`       // Field数据类型
	FieldJson       string `json:"fieldJson"`       // FieldJson
	DataTypeLong    string `json:"dataTypeLong"`    // 数据库字段长度
	Comment         string `json:"comment"`         // 数据库字段描述
	ColumnName      string `json:"columnName"`      // 数据库字段
	FieldSearchType string `json:"fieldSearchType"` // 搜索条件
	DictType        string `json:"dictType"`        // 字典
}

type SysAutoCode struct {
	PackageName string `json:"packageName" gorm:"comment:包名"`
	Label       string `json:"label" gorm:"comment:展示名"`
	Desc        string `json:"desc" gorm:"comment:描述"`
}

type tplData struct {
	template         *template.Template
	interfaceTmp     *template.Template
	locationPath     string
	autoCodePath     string
	autoMoveFilePath string
}

// AutoGenerateCRUD 自动构建数据表对应增删改查函数
func AutoGenerateCRUD(structName, tableName string, autoArg AutoCodeStruct) error {
	// 重定义结构体名称和表名
	autoArg.StructName = structName
	autoArg.TableName = tableName
	autoArg.HumpPackageName = tableName
	return saveAutoCode(autoArg)
}

func saveAutoCode(autoCode AutoCodeStruct) (err error) {
	err, ack := createTemp(autoCode)
	if err != nil {
		return err
	}
	if ack {
		return updateInterface(autoCode)
	}
	return
}

// 根据模板创建文件
func createTemp(autoCode AutoCodeStruct) (err error, ack bool) {
	// 获取需要新增的文件和文件夹路径以及模板文件路径
	dataList, _, needMkdir, err := getNeedList(&autoCode, basePath)
	if err != nil {
		return err, false
	}
	// 写入文件前，根据写入文件先创建文件夹
	if err = file.CreateDir(needMkdir...); err != nil {
		return err, false
	}

	// 生成文件
	for _, value := range dataList {
		f, err := os.OpenFile(value.autoCodePath, os.O_CREATE|os.O_WRONLY, 0o755)
		if err != nil {
			return err, false
		}
		if err = value.template.Execute(f, autoCode); err != nil {
			return err, false
		}
		_ = f.Close()
	}

	defer func() { // 移除中间文件
		if err := os.RemoveAll(autoPath); err != nil {
			return
		}
	}()

	// 最终需要生成的文件
	ackDataList := make([]tplData, 0, len(dataList))
	bf := strings.Builder{}
	if autoCode.AutoMoveFile { // 判断是否需要自动转移
		for index := range dataList {
			// 获取需要移动的路径
			addAutoMoveFile(&dataList[index])
		}
		// 校验待移动路径是否已存在
		for _, value := range dataList {
			if file.FileExist(value.autoMoveFilePath) {
				fmt.Println("目标文件已存在:", value.autoMoveFilePath)
				continue
			}
			ackDataList = append(ackDataList, value)
		}
		for _, value := range ackDataList { // 移动文件
			if err := file.FileMove(value.autoCodePath, value.autoMoveFilePath); err != nil {
				return err, false
			}
		}
		// 保存生成信息
		for _, data := range ackDataList {
			if len(data.autoMoveFilePath) != 0 {
				bf.WriteString(data.autoMoveFilePath)
				bf.WriteString(";")
			}
		}
	}
	if bf.String() != "" {
		fmt.Println("生成文件路径：", bf.String())
		return nil, true
	}
	return nil, false
}

// 更新interface文件
func updateInterface(autoCode AutoCodeStruct) (err error) {
	// 获取需要新增的文件和文件夹路径以及模板文件路径
	dataList, fileList, _, err := getNeedList(&autoCode, interfacePath)
	if err != nil {
		return err
	}
	// 更新文件
	for _, value := range dataList {
		f, err := os.OpenFile(value.autoCodePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o755)
		if err != nil {
			return err
		}
		defer f.Close()
		// 将模板数据转化并暂存到buffer中
		var buf bytes.Buffer
		if err = value.interfaceTmp.Execute(&buf, autoCode); err != nil {
			fmt.Println("interface模板数据转化失败：", err)
			return err
		}
		// 将buffer中数据追加到文件中
		if _, err = f.Write(buf.Bytes()); err != nil {
			fmt.Println("interface文件更新失败：", err)
			return err
		}
	}
	fmt.Println("interface文件已更新：", fileList)
	return nil
}

// 获取要转移的文件路径：根据路径获取文件名，拼接到文件路径到 autoMoveFilePath 字段
func addAutoMoveFile(data *tplData) {
	base := filepath.Base(data.autoCodePath)
	path := data.autoMoveFilePath + base
	data.autoMoveFilePath = filepath.Join(path)
}

// getAllTplFile 获取指定路径下所有tpl文件
func getAllTplFile(pathName string, fileList []string) ([]string, error) {
	files, err := os.ReadDir(pathName)
	if err != nil {
		return nil, err
	}
	for _, fi := range files {
		if fi.IsDir() {
			fileList, err = getAllTplFile(filepath.Join(pathName, fi.Name()), fileList)
			if err != nil {
				return nil, err
			}
		} else {
			if strings.HasSuffix(fi.Name(), ".tpl") {
				fileList = append(fileList, filepath.Join(pathName, fi.Name()))
			}
		}
	}
	return fileList, err
}

// 根据模板文件和参数创建文件夹、文件
func getNeedList(autoCode *AutoCodeStruct, path string) (dataList []tplData, fileList []string, needMkdir []string, err error) {
	// 去除所有空格
	file.TrimSpace(autoCode)
	for _, field := range autoCode.Fields {
		file.TrimSpace(field)
	}
	// 获取 basePath 文件夹下所有tpl文件
	tplFileList, err := getAllTplFile(path, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	dataList = make([]tplData, 0, len(tplFileList))
	fileList = make([]string, 0, len(tplFileList))
	needMkdir = make([]string, 0, len(tplFileList))
	// 根据文件路径生成 tplData 结构体，待填充数据
	for _, value := range tplFileList {
		dataList = append(dataList, tplData{locationPath: value})
	}
	// 读取template文件，生成 *Template, 填充 template 字段(模板文件)
	for index, value := range dataList {
		if path == basePath {
			dataList[index].template, err = template.ParseFiles(value.locationPath)
		} else {
			dataList[index].interfaceTmp, err = template.ParseFiles(value.locationPath)
		}
		if err != nil {
			return nil, nil, nil, err
		}
	}
	// 生成文件路径，填充 autoCodePath 字段(需要生成的文件路径)
	for index, value := range dataList {
		// 裁剪路径前缀
		trimBase := strings.TrimPrefix(value.locationPath, basePath+"/")
		// 判断是否有多级目录，获取第一个/的位置
		if lastSeparator := strings.LastIndex(trimBase, "/"); lastSeparator != -1 {
			// 裁剪目录和后缀，保留文件名
			origFileName := strings.TrimSuffix(trimBase[lastSeparator+1:], ".tpl")
			if origFileName == "interface.go" {
				dataList[index].autoCodePath = autoCode.InterfacePath
				continue
			}
			firstDot := strings.Index(origFileName, ".")
			if firstDot != -1 {
				fileName := autoCode.HumpPackageName + origFileName[firstDot:]
				// 拼接文件路径
				dataList[index].autoCodePath = filepath.Join(autoPath, trimBase[:lastSeparator], autoCode.StructName,
					origFileName[:firstDot], fileName)
				dataList[index].autoMoveFilePath = autoCode.AutoMoveFilePath
			}
		}

		if lastSeparator := strings.LastIndex(dataList[index].autoCodePath, string(os.PathSeparator)); lastSeparator != -1 {
			needMkdir = append(needMkdir, dataList[index].autoCodePath[:lastSeparator])
		}
	}
	for _, value := range dataList {
		fileList = append(fileList, value.autoCodePath)
	}
	return dataList, fileList, needMkdir, err
}
