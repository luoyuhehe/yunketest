package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// SetCreateCallback is set create callback
func SetCreateCallback(db *gorm.DB, callback func(scope *gorm.Scope)) {
	//DefaultCallback.Create().Register("gorm:begin_transaction", beginTransactionCallback)
	//DefaultCallback.Create().Register("gorm:before_create", beforeCreateCallback)
	//DefaultCallback.Create().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	//DefaultCallback.Create().Register("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//DefaultCallback.Create().Register("gorm:create", createCallback)
	//DefaultCallback.Create().Register("gorm:force_reload_after_create", forceReloadAfterCreateCallback)
	//DefaultCallback.Create().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	//DefaultCallback.Create().Register("gorm:after_create", afterCreateCallback)
	//DefaultCallback.Create().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
	db.Callback().Create().Replace("gorm:update_time_stamp", callback)
}

// SetUpdateCallback is set update callback
func SetUpdateCallback(db *gorm.DB, callback func(scope *gorm.Scope)) {
	db.Callback().Update().Replace("gorm:update_time_stamp", callback)
}

// SetDeleteCallback is set delete callback
func SetDeleteCallback(db *gorm.DB, callback func(scope *gorm.Scope)) {
	db.Callback().Delete().Replace("gorm:delete", callback)
}

// UpdateTimeStampForCreateCallback sets `CreateTime`, `ModifyTime` when creating.
func UpdateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := time.Now().Unix()

		if createTimeField, ok := scope.FieldByName("CreateTime"); ok {
			if createTimeField.IsBlank {
				if err := createTimeField.Set(now); err != nil {
					panic(err)
				}
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifyTime"); ok {
			if modifyTimeField.IsBlank {
				if err := modifyTimeField.Set(now); err != nil {
					panic(err)
				}
			}
		}
	}
}

// UpdateTimeStampForUpdateCallback sets `ModifyTime` when updating.
func UpdateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		now := time.Now().Unix()
		if err := scope.SetColumn("ModifyTime", now); err != nil {
			panic(err)
		}
	}
}

// DeleteCallback used to delete data from database or set DeletedTime to
// current time and is_del = 1(when using with soft delete).
func DeleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedTimeField, hasDeletedAtField := scope.FieldByName("DeleteTime")
		isDelField, hasIsDelField := scope.FieldByName("IsDelete")

		if !scope.Search.Unscoped && hasDeletedAtField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedTimeField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
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

// addExtraSpaceIfExist used to add extra space if exists str.
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
