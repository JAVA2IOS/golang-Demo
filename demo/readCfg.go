package main 

import (
	"fmt"
	"log"
	"github.com/larspensjo/config"
)


// 配置文件
const (
	configFile = "config.ini"	// 配置文件名字

	orderSectionName = "order"	// 订单section配置数据
	orderFileRegularParameter = "orderXlsxRegularString"	// 订单文件名匹配字符串参数键值
	orderFileRegularString = "orders_export   #订单文件名匹配字符串"	// 订单文件名匹配字符串
	orderFileNeedCarriedParameter = "newSavedColumn"	// 订单文件需要保存到新文件的字段数据
	orderFileNeedCarriedString = "订单号   #订单文件需要保存到新文件的字段数据"	// 订单文件需要保存到新文件的字段数据

	orderFileIndexColumnNameParameter = "orderIndexColumnName"	// 订单文件索引条件列名参数键值
	orderFileIndexColumnName = "收货人   #订单文件索引条件列名"	// 订单文件索引条件列名

	expressSectionName = "express"	// 快递单section配置数据
	expressFileRegularParameter = "expressXlsxRegularString"	// 订单文件名匹配字符串参数键值
	expressFileRegularString = "订单   #快递单文件名匹配字符串"	// 快递单文件名匹配字符串

	expressFileIndexColumnNameParameter = "expressIndexColumnName"	// 快递单文件索引条件列名参数键值
	expressFileIndexColumnName = "收货人   #快递单文件索引条件列名，匹配上方的orderIndexColumnName字段，相等匹配成功"	// 快递单文件索引条件列名

	expressFileSavedColumnNameParameter = "expressNewXlsxFileSavedColumnName"
	expressFileSavedColumnNameString = "货运单号   #需要保存到新文件的字段名称"

	newFileSectionName = "newFile"	// 新文件section配置数据
	newFileNewColumnParameter = "newXlsxFileColumnName"
	newFileNewColumnString = "快递公司   #需要填充的字段名称"

	newFileNewColumnValueParameter = "newXlsxFileColumnValue"
	newFileNewColumnValueString = "顺丰快递   #需要填充的字段值"

	newFileOrderFileColumnEnName = "newFileOrderFileColumnEnName"
	newFileOrderFileColumnEnValue = "order_sn   # 订单号字段保存到当前文件的英文字段名称"

	newFileNewColumnEnName = "newFileNewColumnEnName"
	newFileNewColumnEnValue = "shipping_name   # 新文件新增的字段英文名称"

	newFileExpressFileColumnEnName = "newFileExpressFileColumnEnName"
	newFileExpressFileColumnEnValue = "shipping_sn   # 快递单号字段保存到当前文件的英文字段名称"
)

var (
	TOPIC = make(map[string]string)
)

func main() {
	cfg, err := readConfigureFile()
 
	if err != nil {
		log.Fatalf("读取配置文件失败(%s): %v", configFile, err)
	}else {
		fmt.Println("读取配置文件成功......")
	}
 	
	if cfg.HasSection(orderSectionName) {   //判断配置文件中是否有section（一级标签）
		options,err := cfg.SectionOptions(orderSectionName)    //获取一级标签的所有子标签options（只有标签没有值）
		if err == nil {
			for _,v := range options{
				optionValue, err := cfg.String(orderSectionName, v)  //根据一级标签section和option获取对应的值
				if err == nil {
					TOPIC[v] = optionValue
					fmt.Println("当前值:", optionValue)
				}
			}
		}
	}
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
	cfg.AddOption(orderSectionName, orderFileRegularParameter, orderFileRegularString)
	cfg.AddOption(orderSectionName, orderFileIndexColumnNameParameter, orderFileIndexColumnName)
	cfg.AddOption(orderSectionName, orderFileNeedCarriedParameter, orderFileNeedCarriedString)

	// 快递单配置数据
	cfg.AddSection(expressSectionName)
	cfg.AddOption(expressSectionName, expressFileRegularParameter, expressFileRegularString)
	cfg.AddOption(expressSectionName, expressFileIndexColumnNameParameter, expressFileIndexColumnName)
	cfg.AddOption(expressSectionName, expressFileSavedColumnNameParameter, expressFileSavedColumnNameString)

	// 新文件保存的字段
	cfg.AddSection(newFileSectionName)
	cfg.AddOption(newFileSectionName, newFileNewColumnParameter, newFileNewColumnString)
	cfg.AddOption(newFileSectionName, newFileNewColumnValueParameter, newFileNewColumnValueString)
	cfg.AddOption(newFileSectionName, newFileOrderFileColumnEnName, newFileOrderFileColumnEnValue)
	cfg.AddOption(newFileSectionName, newFileNewColumnEnName, newFileNewColumnEnValue)
	cfg.AddOption(newFileSectionName, newFileExpressFileColumnEnName, newFileExpressFileColumnEnValue)

	cfg.WriteFile(configFile, 0644, "关于xlsx文件匹配参数定义")
}



/*======================================*/