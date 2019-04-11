package main
/*
	生成windows可执行的文件 .exe CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build qqls_ExcelHelper.go 

	window下的文件路径  ‘/’必须改为'\\' 反斜杠

	v2.0 自动写死输入的文件
*/
import (
	"fmt"
	"os"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"time"
)

const (
	CodeZNotFound = -1	// 未找到
	CodeZOrderXlsxFileName = "订单"	// 订单文件名
	CodeZExpressXlsxFileName = "货运单" // 货运单文件名
	savedSheetName = "newSheet"	// 保存的表单Sheet名称
	CodezXlsxFileSuffix = ".xlsx"	// xlsx文件后缀名
)

 
func main() {
	const firSheetColumnName string = "收货人"		// 第一张表单（需要订单号和手机号）字段列名
	const secondSheetColumnName string = "收货人"	// 第二张表单(待匹配的数据)列名
	const orderIdColumnName string = "订单号"	// 第一张表单订单号列名

	const expressColumnName string = "货运单号"	// 需要获取的货运单号

	// 保存的新文件需要
	const savedExpressColumnName string = "快递单号"	// 待保存的文件货运单号列名
	const savedExpressDefaultFieldValue string = "顺丰快递"	// 自动填充的字段值

	var CodezSlash string = "/"		// 斜杠
	if true {
		// windows下需要反斜杠
		CodezSlash = "\\"	// 反斜杠
	}

	dir, err := currentFilePath()
	if err == nil {
		fmt.Println("当前文件路径：", dir)
	}
	
	// 订单文件路径
	var orderFilePath string = dir + CodezSlash + CodeZOrderXlsxFileName + CodezXlsxFileSuffix

	fmt.Println("\n.\n.\n.\n\n=============订单号文件：", orderFilePath, "请确保表单内有\"", firSheetColumnName, "\", \"", orderIdColumnName, "\"字段==================\n.\n.\n.\n\n")
	// 打开文件
	firstFile, err := openXlsxFile(orderFilePath)
	if firstFile == nil {
		fmt.Println("文件(", orderFilePath, ")打开失败:", err.Error())
		return
	}

	// 货运单文件路径
	var expressFilePath string = dir + CodezSlash + CodeZExpressXlsxFileName + CodezXlsxFileSuffix

	fmt.Println("\n.\n.\n.\n\n=============快递单号文件：", expressFilePath, "请确保表单内有\"", secondSheetColumnName, "\", \"", expressColumnName, "\"字段==================\n.\n.\n.\n\n")

	secondFile, _ := openXlsxFile(expressFilePath)

	if secondFile == nil {
		fmt.Println("文件(", expressFilePath, ")打开失败:", err.Error())
		return
	}

	// excel文件中的表单sheet
	firstSheet, secondSheet := firstFile.Sheets[0], secondFile.Sheets[0]

 	fmt.Println("开始加载订单号表单(", orderFilePath, ")数据......")

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
	newCell.Value = "订单号"

	newCell = newRow.AddCell()
	newCell.Value = "快递公司"

	newCell = newRow.AddCell()
	newCell.Value = savedExpressColumnName

	newRow = finalSheet.AddRow()
	newCell = newRow.AddCell()
	newCell.Value = "order_sn"

	newCell = newRow.AddCell()
	newCell.Value = "shipping_name"

	newCell = newRow.AddCell()
	newCell.Value = "shipping_sn"


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
							newCell.Value = "顺丰快递"
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

 	str_time := time.Unix(1389058332, 0).Format("2006_01_02_150405")
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

// 当前文件所在路径
func currentFilePath() (string, error) {
	return filepath.Abs("./")
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




