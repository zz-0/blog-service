package model

//model层(模型层)连接数据库并从数据库中获取对应的字段

import (
	"blog-service/global"
	"blog-service/pkg/setting"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	STATE_OPEN  = 1
	STATE_CLOSE = 0
)

//三张表共用的字段
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"` //文章ID
	CreatedBy  string `json:"created_by"`            //创建人
	ModifiedBy string `json:"modified_by"`           //修改人
	CreatedOn  uint32 `json:"created_on"`            //创建时间
	ModifiedOn uint32 `json:"modified_on"`           //修改时间
	DeletedOn  uint32 `json:"deleted_on"`            //删除时间
	IsDel      uint8  `json:"is_del"`                //删除确认
}

//创建DB实例(连接数据库的语句),引入gorm开源库和mysql驱动库
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	//获取数据库的连接
	db, err := gorm.Open(databaseSetting.DBType,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
			databaseSetting.UserName,
			databaseSetting.PassWord,
			databaseSetting.Host,
			databaseSetting.DBName,
			databaseSetting.Charset,
			databaseSetting.ParseTime,
		))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}

	//在grom中表名是结构体名的复数形式，此句作用是让grom转义struct名字的时候不用加上s
	db.SingularTable(true)

	//注册下列的回调行为
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	//至此公共字段的处理完成

	//设置最大连接数
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}

//为什么要编写回调函数？

//新增行为的回调
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		//通过scope.FieldByName方法，获取当前是否包含所需的字段
		if createTimeField, ok := scope.FieldByName("CreateOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			//判断modifyTimeField.IsBlank的值，得知该字段的值是否为空
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

//更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	//调用scope.Get("gorm:update_column")来获取当前设置的标识gorm:update_column的字段属性
	//若不存在，即没有自定义设置update_column,则在更新回调内设置默认字段ModifiedOn的值为当前的时间戳
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

//删除行为的回调
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string

		//获取当前设置的标识gorm:delete_option的字段属性
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")

		//判断是否存在deleteon和isdel字段。若存在则执行update操作进行软删除(修改deletedon和isdel的值),否则执行delete进行硬删除
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",

				//获取当前引用的表名,在对sql语句的组成部分进行处理和转移
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),

				//调用scope.CombinedConditionSql完成sql语句的组装
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return "" + str
	}
	return " "
}
