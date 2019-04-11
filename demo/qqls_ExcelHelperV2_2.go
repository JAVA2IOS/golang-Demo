package main
/*
	生成windows可执行的文件 .exe CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build qqls_ExcelHelperV2_2.go

	window下的文件路径  ‘/’必须改为'\\' 反斜杠

	v2.0 自动写死输入的文件

	v2.1 自动读取特定的输入文件

	v2.2 从配置文件中获取参数
*/
import (
	"fmt"
	"os"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"time"
	"io/ioutil"
	"strings"
	"github.com/larspensjo/config"
)

// 配置文件变量值
const (
	configFile = "config.ini"	// 配置文件名字

	orderSectionName = "order"	// 订单section配置数据
	orderFileRegularParameter = "orderXlsxRegularString"	// 订单文件名匹配字符串参数键值
	orderFileRegularString = "orders_export"	// 订单文件名匹配字符串
	orderFileNeedCarriedParameter = "newSavedColumn"	// 订单文件需要保存到新文件的字段数据
	orderFileNeedCarriedString = "订单号"	// 订单文件需要保存到新文件的字段数据

	orderFileSaveColumnNewName = "newSavedNewColumn"	// 订单表中字段保存到新文件中显示的列名
	orderFileSaveColumnNewString = "订单号"	// 订单表中字段保存到新文件中显示的列名

	orderFileIndexColumnNameParameter = "orderIndexColumnName"	// 订单文件索引条件列名参数键值
	orderFileIndexColumnName = "收货人"	// 订单文件索引条件列名

	expressSectionName = "express"	// 快递单section配置数据
	expressFileRegularParameter = "expressXlsxRegularString"	// 订单文件名匹配字符串参数键值
	expressFileRegularString = "订单"	// 快递单文件名匹配字符串

	expressFileIndexColumnNameParameter = "expressIndexColumnName"	// 快递单文件索引条件列名参数键值
	expressFileIndexColumnName = "收货人"	// 快递单文件索引条件列名

	expressFileSavedColumnNameParameter = "expressNewXlsxFileSavedColumnName"
	expressFileSavedColumnNameString = "货运单号"

	expressFileNewSavedColumnNameParameter = "expressFileNewSavedColumnNameParameter"
	expressFileNewSavedColumnNameString = "快递单号"

	newFileSectionName = "newFile"	// 新文件section配置数据
	newFileNewColumnParameter = "newXlsxFileColumnName"
	newFileNewColumnString = "快递公司"

	newFileNewColumnValueParameter = "newXlsxFileColumnValue"
	newFileNewColumnValueString = "顺丰快递"

	newFileOrderFileColumnEnName = "newFileOrderFileColumnEnName"
	newFileOrderFileColumnEnValue = "order_sn"

	newFileNewColumnEnName = "newFileNewColumnEnName"
	newFileNewColumnEnValue = "shipping_name"

	newFileExpressFileColumnEnName = "newFileExpressFileColumnEnName"
	newFileExpressFileColumnEnValue = "shipping_sn"
)

// 正常参数
const (
	CodeZNotFound = -1	// 未找到
	savedSheetName = "newSheet"	// 保存的表单Sheet名称
	CodezXlsxFileSuffix = ".xlsx"	// xlsx文件后缀名
)

 
func main() {

	cfg, err := readConfigureFile()
 
	if err != nil {
		fmt.Println("读取配置文件失败(%s): %v", configFile, err)
	}else {
		fmt.Println("读取配置文件(", configFile, ")成功......")
	}

	CodeZOrderXlsxFileRegularName := getConfigValue(cfg, orderSectionName, orderFileRegularParameter, orderFileRegularString)
	// 货运单正则匹配字段
	CodeZExpressXlsxFileRegularName := getConfigValue(cfg, expressSectionName, expressFileRegularParameter, expressFileRegularString)


	fmt.Println(" \n.\n..\n...\n")
	fmt.Println("正在加载xlsx文件路径......")
	fmt.Println(".\n.\n.\n")

	dir, err := currentFilePath()
	if err == nil {
		fmt.Println("当前文件路径：", dir)
	}
	fmt.Println(".\n.\n.\n")
	fmt.Println("正在读取订单excel文件......")
	CodeZOrderXlsxFileName := readSpecifiedXlsxFile(CodeZOrderXlsxFileRegularName)	// 订单文件名
	if CodeZOrderXlsxFileName == "" {
		fmt.Println("订单excel文件(", CodeZOrderXlsxFileName, ")读取失败......")
		return
	}else {
		fmt.Println("订单excel文件(", CodeZOrderXlsxFileName, ")读取成功......")
	}

	fmt.Println("正在读取货运单excel文件......")
	CodeZExpressXlsxFileName := readSpecifiedXlsxFile(CodeZExpressXlsxFileRegularName) // 货运单文件名
	if CodeZExpressXlsxFileName == "" {
		fmt.Println("货运单excel文件(", CodeZExpressXlsxFileName, ")读取失败......")
		return
	}else {
		fmt.Println("货运单excel文件(", CodeZExpressXlsxFileName, ")读取成功......")
	}

	// 第一张表单（需要订单号和手机号）字段列名
	firSheetColumnName := getConfigValue(cfg, orderSectionName, orderFileIndexColumnNameParameter, orderFileIndexColumnName)
	//  "收货人"	第二张表单(待匹配的数据)列名
	secondSheetColumnName := getConfigValue(cfg, expressSectionName, expressFileIndexColumnNameParameter, expressFileIndexColumnName)
	// 第一张表单订单号列名 "订单号"	
	orderIdColumnName := getConfigValue(cfg, orderSectionName, orderFileNeedCarriedParameter, orderFileNeedCarriedString)
	// 需要获取的货运单号 "货运单号"	
	expressColumnName := getConfigValue(cfg, expressSectionName, expressFileSavedColumnNameParameter, expressFileSavedColumnNameString)

	// 保存的新文件需要
	//  "快递单号"	待保存的文件货运单号列名
	savedExpressColumnName := getConfigValue(cfg, expressSectionName, expressFileNewSavedColumnNameParameter, expressFileNewSavedColumnNameString)
	// "顺丰快递"	自动填充的字段值
	savedExpressDefaultFieldValue := getConfigValue(cfg, newFileSectionName, newFileNewColumnParameter, newFileNewColumnValueString)

	var CodezSlash string = "/"		// 斜杠
	if true {
		// windows下需要反斜杠
		CodezSlash = "\\"	// 反斜杠
	}

	// 订单文件路径
	CodeZOrderXlsxFileName = strings.Replace(CodeZOrderXlsxFileName, CodezXlsxFileSuffix, "", -1)
	var orderFilePath string = dir + CodezSlash + CodeZOrderXlsxFileName + CodezXlsxFileSuffix

	fmt.Println("\n.\n.\n.\n\n订单号文件：", orderFilePath, "\n=== 请确保表单内有\"", firSheetColumnName, "\", \"", orderIdColumnName, "\"字段 ===\n.\n.\n.\n\n")
	// 打开文件
	firstFile, err := openXlsxFile(orderFilePath)
	if firstFile == nil {
		fmt.Println("文件(", orderFilePath, ")打开失败:", err.Error())
		return
	}

	// 货运单文件路径
	CodeZExpressXlsxFileName = strings.Replace(CodeZExpressXlsxFileName, CodezXlsxFileSuffix, "", -1)
	var expressFilePath string = dir + CodezSlash + CodeZExpressXlsxFileName + CodezXlsxFileSuffix

	fmt.Println("快递单号文件：", expressFilePath, "\n=== 请确保表单内有\"", secondSheetColumnName, "\", \"", expressColumnName, "\"字段 ===\n.\n.\n.\n\n")

	secondFile, _ := openXlsxFile(expressFilePath)

	if secondFile == nil {
		fmt.Println("文件(", expressFilePath, ")打开失败:", err.Error())
		return
	}

	// excel文件中的表单sheet
	firstSheet, secondSheet := firstFile.Sheets[0], secondFile.Sheets[0]

 	fmt.Println("开始加载订单号表单(", orderFilePath, ")数据......\n\n")

 	firstTargetIndex, orderIdIndex, secondTargetIndex, expressIndex := CodeZNotFound, CodeZNotFound, CodeZNotFound, CodeZNotFound
	firstResult := []Orders{}

	// 获取收货人姓名和订单号
 	for rowIndex, row := range firstSheet.Rows {
 		if rowIndex == 0 {
 			// 跳过第一行，行名
 			var order Orders
 			for cellIndex, cell := range row.Cells {
				if cell.Type() == xlsx.CellTypeString {
					if cell.String() == firSheetColumnName {
						firstTargetIndex = cellIndex
						fmt.Println("订单号表单\"", firSheetColumnName, "\"所在列坐标：", firstTargetIndex)
						order.customerName = cell.String()
					}
					if cell.String() == orderIdColumnName {
						orderIdIndex = cellIndex
						fmt.Println("订单号表单\"", orderIdColumnName, "\"所在列坐标：", orderIdIndex)
						order.orderId = cell.String()
					}
				}
			}
			firstResult = append(firstResult, order)
			fmt.Println("")
			continue
 		}
		// 获取当前数据
		for cellIndex, cell := range row.Cells {
			if firstTargetIndex != CodeZNotFound {
				if firstTargetIndex == cellIndex {
					orderIdString := ""
					if orderIdIndex != CodeZNotFound {
						orderIdString = row.Cells[orderIdIndex].String()
					}
					order := Orders{customerName: cell.String(), orderId: orderIdString}
					fmt.Println("索引数据：", order)
					firstResult = append(firstResult, order)
				}
			}else {
				fmt.Println("索引数据失败：请确定订单号表单(", CodeZOrderXlsxFileName, CodezXlsxFileSuffix, ")内有\"", firSheetColumnName, "\"字段")
				break;
			}
		}
	}

	fmt.Println("索引完毕")
	fmt.Println(".")
	fmt.Println(".")
	fmt.Println(".")
	fmt.Println("")
 	fmt.Println("加载快递表单(", expressFilePath, ")数据......")
	fmt.Println("")

 	finalXlsx := xlsx.NewFile()
    finalSheet, err := finalXlsx.AddSheet(savedSheetName)
    if err != nil {
    	fmt.Println("\n.\n.\n.\n\n创建新表单失败: ", err)
    	return
    }

    // 写死
    newRow := finalSheet.AddRow()
	newCell := newRow.AddCell()
	// "订单号"
	newCell.Value = getConfigValue(cfg, orderSectionName, orderFileSaveColumnNewName, orderFileSaveColumnNewString)


	newCell = newRow.AddCell()
	// "快递公司"
	newCell.Value = savedExpressDefaultFieldValue

	newCell = newRow.AddCell()
	newCell.Value = savedExpressColumnName

	newRow = finalSheet.AddRow()
	newCell = newRow.AddCell()
	// "order_sn" 
	newCell.Value = getConfigValue(cfg, newFileSectionName, newFileOrderFileColumnEnName, newFileOrderFileColumnEnValue)

	newCell = newRow.AddCell()
	// "shipping_name" 
	newCell.Value = getConfigValue(cfg, newFileSectionName, newFileNewColumnEnName, newFileNewColumnEnValue)

	newCell = newRow.AddCell()
	// "shipping_sn"
	newCell.Value = getConfigValue(cfg, newFileSectionName, newFileExpressFileColumnEnName, newFileExpressFileColumnEnValue)


	// 匹配快递单号
 	for _, orderObj := range firstResult {

 		for rowIndex, row := range secondSheet.Rows {
 			// 列名
 			if rowIndex == 0 {
 				if secondTargetIndex == CodeZNotFound {
	 				for cellIndex, nameCell := range row.Cells {
	 					if nameCell.String() == expressColumnName {
							expressIndex = cellIndex
	 					}
	 					if nameCell.String() == secondSheetColumnName {
							secondTargetIndex = cellIndex
						}
	 				}
 				}
 			}else {
 				if secondTargetIndex != CodeZNotFound {
 					if expressIndex == CodeZNotFound {
 						fmt.Println("匹配快递单号失败：请确定表单(", CodeZExpressXlsxFileName, CodezXlsxFileSuffix, ")内有\"", expressColumnName, "\"字段")
 						break
 					}
	 				targetCell := row.Cells[secondTargetIndex]
					if targetCell.String() == orderObj.customerName {
						if row.Cells[expressIndex].String() != "" {
							fmt.Printf("匹配到第%4d 行:  ", rowIndex)
							newRow := finalSheet.AddRow()
	 						newCell := newRow.AddCell()
							newCell.Value = orderObj.orderId
							fmt.Printf("%-22s", orderObj.orderId)

	 						newCell = newRow.AddCell()
							newCell.Value =getConfigValue(cfg, newFileSectionName, newFileNewColumnValueParameter, newFileNewColumnValueString) 
							fmt.Printf("%8s", newCell.Value)

							newCell = newRow.AddCell()
							newCell.Value = row.Cells[expressIndex].String()
							fmt.Printf("%14s", newCell.Value)
							fmt.Printf("\n")
						}
		 			}
	 			}else {
	 				fmt.Println("匹配快递单号失败：请确定表单(", CodeZExpressXlsxFileName, CodezXlsxFileSuffix, ")内有\"", secondSheetColumnName, "\"字段")
	 				break
	 			}
 			}
 		}
 	}

 	str_time := time.Now().Format("2006_01_02_150405")
 	newSavedFileName := "时间-顺丰发货表" + str_time + CodezXlsxFileSuffix
 	newFilePath := dir + CodezSlash + newSavedFileName

	newXlsxFileName := newFilePath

 	newError := finalXlsx.Save(newXlsxFileName)

 	if newError != nil {
 		fmt.Println("\n.\n.\n.\n\n文件(", newSavedFileName, ")保存失败, 原因:", newError.Error())
 	}else {
 		fmt.Println("\n.\n.\n.\n\n加载完毕，匹配到的数据保存到(", newSavedFileName, "), 请查收")
 	}


 	time.Sleep(10 * time.Second)
}


/* ============================= 读取默认配置文件 ============================ */
func readConfigureFile() (*config.Config, error) {
	cfg, err := config.ReadDefault(configFile)
    if err != nil {
    	fmt.Println("读取配置文件失败，配置文件重新初始化......")
    	createDefaultConfigFile()
    	return config.ReadDefault(configFile)
    }

    return cfg, err
}


// 创建默认的配置文件
func createDefaultConfigFile() {
	cfg := config.NewDefault()
	// 订单配置数据
	cfg.AddSection(orderSectionName)
	cfg.AddOption(orderSectionName, orderFileRegularParameter, orderFileRegularString + "   #订单文件名匹配字符串")
	cfg.AddOption(orderSectionName, orderFileIndexColumnNameParameter, orderFileIndexColumnName + "   #订单文件索引条件列名")
	cfg.AddOption(orderSectionName, orderFileNeedCarriedParameter, orderFileNeedCarriedString + "   #订单文件需要保存到新文件的字段数据")
	cfg.AddOption(orderSectionName, orderFileSaveColumnNewName, orderFileSaveColumnNewString + "   #订单表中字段保存到新文件中显示的列名")

	// 快递单配置数据
	cfg.AddSection(expressSectionName)
	cfg.AddOption(expressSectionName, expressFileRegularParameter, expressFileRegularString + "   #快递单文件名匹配字符串")
	cfg.AddOption(expressSectionName, expressFileIndexColumnNameParameter, expressFileIndexColumnName + "   #快递单文件索引条件列名，匹配上方的orderIndexColumnName字段，相等匹配成功")
	cfg.AddOption(expressSectionName, expressFileSavedColumnNameParameter, expressFileSavedColumnNameString + "   #需要保存到新文件的字段名称")
	cfg.AddOption(expressSectionName, expressFileNewSavedColumnNameParameter, expressFileNewSavedColumnNameString + "   #保存到新文件上时显示的字段列名")

	// 新文件保存的字段
	cfg.AddSection(newFileSectionName)
	cfg.AddOption(newFileSectionName, newFileNewColumnParameter, newFileNewColumnString + "   #需要填充的字段名称")
	cfg.AddOption(newFileSectionName, newFileNewColumnValueParameter, newFileNewColumnValueString + "   #需要填充的字段值")
	cfg.AddOption(newFileSectionName, newFileOrderFileColumnEnName, newFileOrderFileColumnEnValue + "   # 订单号字段保存到当前文件的英文字段名称")
	cfg.AddOption(newFileSectionName, newFileNewColumnEnName, newFileNewColumnEnValue + "   # 新文件新增的字段英文名称")
	cfg.AddOption(newFileSectionName, newFileExpressFileColumnEnName, newFileExpressFileColumnEnValue + "   # 快递单号字段保存到当前文件的英文字段名称")

	cfg.WriteFile(configFile, 0644, "关于xlsx文件匹配参数定义")
}

// 读取配置文件数据
func getConfigValue(cfg *config.Config, section string, key string, defaultValue string) string {
	str, err := cfg.String(section, key)
	if err != nil {
		return defaultValue
	}

	return str
}



/* ============================= 读取xlsx文件 ============================ */
// 当前文件所在路径
func currentFilePath() (string, error) {
	return filepath.Abs("./")
}

// 读取特定excel文件名
func readSpecifiedXlsxFile(regularString string) string {
	fileName := ""
	currentDirPath, err := currentFilePath()
	if err == nil {
		files, _ := ioutil.ReadDir(currentDirPath)
	    for _, file := range files {
	        if file.IsDir() {
	            continue
	        } else {
	        	if strings.Contains(file.Name(), regularString)  {
	        		fileName = file.Name()
	        		break
	        	}
	        }
	    }
	}

	return fileName
}


// 打开文件获取订单信息
func openXlsxFile(filePath string) (*xlsx.File, error){
	xlsxFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Println("文件打开失败，原因:", err)
		os.Exit(1)
		return nil, err
	}

	return xlsxFile, nil
}

func getSheetByName(file *xlsx.File, sheetName string) *xlsx.Sheet {
	for _, sheet := range file.Sheets {
		if sheet.Name == sheetName {
			return sheet
			break
		}
	}

	return nil
}

// 订单号
type Orders struct {
	orderId string
	phoneNumber string
	customerName string
}




