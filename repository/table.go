package repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kazu22002/dbRelation/util"
)

type DbInfo struct {
	Table DbTable
	Column []DbColumn
	TableNameId string
}

type DbTable struct {
	TableName string
}

type DbColumn struct {
	ColumnName string
}

func Table(db *gorm.DB) []DbInfo{
	var result []DbInfo
	var tables []DbTable

	// SELECT
	//rows, err := db.Query(`select relname as TABLE_NAME from pg_stat_user_tables`)
	db.Raw(`select relname as table_name from pg_stat_user_tables`).Find(&tables)
	for _, v := range tables {
		column := Column(db, v.TableName)

		r := DbInfo{
			Table: v,
			Column:column,
			TableNameId: util.SingleName(v.TableName)+"_id",
		}
		result = append(result, r)
	}
	return result
}

func Column(db *gorm.DB, table string) []DbColumn {
	var result []DbColumn

	db.Raw(
	`select column_name from information_schema.columns 
        where 
		table_name=?
		order by
		ordinal_position;
`, table).Find(&result)

	return result

}