package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/suzujun/gendao/helper"
)

type (
	Connection struct {
		user     string
		password string
		dbname   string
		db       *sql.DB
	}
)

func NewConnection(user, password, dbname string, close bool) (*Connection, error) {
	con := Connection{
		user:     user,
		password: password,
		dbname:   dbname,
	}
	if err := con.Open(); err != nil {
		return nil, err
	}
	if close {
		defer con.Close()
	}
	return &con, nil
}

func (con *Connection) Open() error {
	if con.db == nil {
		if con.dbname == "" {
			return errors.New("No database name selected in config")
		}
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", con.user, con.password, con.dbname))
		if err != nil {
			return err
		}
		if err = db.Ping(); err != nil {
			return err
		}
		con.db = db
	}
	return nil
}

func (con *Connection) Close() error {
	if con.db == nil {
		return nil
	}
	return con.db.Close()
}

func (con *Connection) GetTableNames() ([]string, error) {

	if con.db == nil {
		return nil, errors.New("database is closed")
	}
	rows, err := con.db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}

	tnames := []string{}
	for rows.Next() {
		var tname string
		err := rows.Scan(&tname)
		if err != nil {
			return nil, err
		}
		tnames = append(tnames, tname)
	}
	return tnames, nil
}

func (con *Connection) GetTable(tableName string) (*Table, error) {
	columns, err := con.GetColumns(tableName)
	if err != nil {
		return nil, err
	}
	indexes, err := con.GetIndexes(tableName)
	if err != nil {
		return nil, err
	}
	mt := Table{}
	mt.Columns = columns
	mt.Indexes = indexes
	if len(columns) > 0 {
		mt.Catalog = columns[0].TableCatalog
		mt.Schema = columns[0].TableSchema
		mt.Name = columns[0].TableName
	}
	return &mt, nil
}

func (con *Connection) GetColumns(tname string) ([]Column, error) {

	if con.db == nil {
		return nil, errors.New("database is closed")
	}
	rows, err := con.db.Query(fmt.Sprintf(`
select TABLE_CATALOG, TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION,
 COLUMN_DEFAULT, IS_NULLABLE, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH,
 CHARACTER_OCTET_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, DATETIME_PRECISION,
 CHARACTER_SET_NAME, COLLATION_NAME, COLUMN_TYPE, COLUMN_KEY, EXTRA, PRIVILEGES,
 COLUMN_COMMENT
from INFORMATION_SCHEMA.COLUMNS
where TABLE_SCHEMA = '%s'
and TABLE_NAME = '%s'
order by TABLE_NAME, ORDINAL_POSITION
`, con.dbname, tname))
	if err != nil {
		panic(err)
	}

	var tableCatalog, tableSchema, tableName, columnName, isNullable, dataType,
		columnType, columnKey, extra, privileges,
		columnComment string
	var characterSetName, collationName sql.NullString
	var ordinalPosition uint
	var characterMaximumLength, characterOctetLength, numericPrecision,
		numericScale, datetimePrecision sql.NullInt64
	var columnDefault interface{}

	result := []Column{}
	for rows.Next() {
		err = rows.Scan(
			&tableCatalog, &tableSchema, &tableName, &columnName, &ordinalPosition,
			&columnDefault, &isNullable, &dataType, &characterMaximumLength,
			&characterOctetLength, &numericPrecision, &numericScale,
			&datetimePrecision, &characterSetName, &collationName, &columnType,
			&columnKey, &extra, &privileges, &columnComment,
		)
		if err != nil {
			panic(err)
		}
		result = append(result, Column{
			TableCatalog:           tableCatalog,
			TableSchema:            tableSchema,
			TableName:              tableName,
			ColumnName:             columnName,
			OrdinalPosition:        ordinalPosition,
			ColumnDefault:          columnDefault,
			IsNullable:             isNullable == "YES",
			DataType:               dataType,
			CharacterMaximumLength: helper.ParseIntPointer(&characterMaximumLength),
			CharacterOctetLength:   helper.ParseIntPointer(&characterOctetLength),
			NumericPrecision:       helper.ParseIntPointer(&numericPrecision),
			NumericScale:           helper.ParseIntPointer(&numericScale),
			DatetimePrecision:      helper.ParseIntPointer(&datetimePrecision),
			CharacterSetName:       helper.ParseStringPointer(&characterSetName),
			CollationName:          helper.ParseStringPointer(&collationName),
			ColumnType:             columnType,
			ColumnKey:              columnKey,
			Extra:                  extra,
			Privileges:             privileges,
			ColumnComment:          columnComment,
		})
	}
	return result, nil
}

func (con *Connection) GetIndexes(tname string) ([]Index, error) {

	if con.db == nil {
		return nil, errors.New("database is closed")
	}
	rows, err := con.db.Query(fmt.Sprintf(`
select TABLE_CATALOG, TABLE_SCHEMA, TABLE_NAME, NON_UNIQUE, INDEX_SCHEMA,
  INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME, COLLATION, CARDINALITY, SUB_PART, PACKED,
  NULLABLE, INDEX_TYPE, COMMENT, INDEX_COMMENT,
  CASE
    WHEN INDEX_NAME = "PRIMARY" THEN 1
    WHEN NON_UNIQUE= 0 THEN 2
    ELSE 3
  END SORT_NUMBER
from information_schema.statistics
where table_schema = "%s"
and TABLE_NAME = "%s"
order by SORT_NUMBER, INDEX_NAME, SEQ_IN_INDEX
`, con.dbname, tname))
	if err != nil {
		panic(err)
	}

	var tableCatalog, tableSchema, tableName, indexSchema, indexName, columnName, collation, nullable,
		indexType, comment, indexComment string
	var subPart, packed sql.NullString
	var nonUnique, seqInIndex, cardinality, sortNumber uint

	result := []Index{}
	for rows.Next() {
		err = rows.Scan(
			&tableCatalog, &tableSchema, &tableName, &nonUnique, &indexSchema, &indexName, &seqInIndex,
			&columnName, &collation, &cardinality, &subPart, &packed, &nullable, &indexType, &comment,
			&indexComment, &sortNumber,
		)
		if err != nil {
			panic(err)
		}
		result = append(result, Index{
			TableCatalog: tableCatalog,
			TableSchema:  tableSchema,
			TableName:    tableName,
			NonUnique:    nonUnique,
			IndexSchema:  indexSchema,
			IndexName:    indexName,
			SeqInIndex:   seqInIndex,
			ColumnName:   columnName,
			Collation:    collation,
			Cardinality:  cardinality,
			SubPart:      helper.ParseStringPointer(&subPart),
			Packed:       helper.ParseStringPointer(&packed),
			Nullable:     nullable == "YES",
			IndexType:    indexType,
			Comment:      comment,
			IndexComment: indexComment,
		})
	}
	return result, nil
}
