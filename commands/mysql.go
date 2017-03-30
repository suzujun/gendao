package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type (
	MysqlConnection struct {
		user     string
		password string
		dbname   string
		db       *sql.DB
	}
	MysqlTableJSON struct {
		Catalog string        `json:"catalog"`
		Schema  string        `json:"schema"`
		Name    string        `json:"name"`
		Columns []MysqlColumn `json:"columns"`
		Indexs  []MysqlIndex  `json:"indexs"`
	}
	MysqlColumn struct {
		TableCatalog           string      `json:"tableCatalog" db:"TABLE_CATALOG"`
		TableSchema            string      `json:"tableSchema" db:"TABLE_SCHEMA"`
		TableName              string      `json:"tableName" db:"TABLE_NAME"`
		ColumnName             string      `json:"columnName" db:"COLUMN_NAME"`
		OrdinalPosition        uint        `json:"ordinalPosition" db:"ORDINAL_POSITION"`
		ColumnDefault          interface{} `json:"columnDefault" db:"COLUMN_DEFAULT"`
		IsNullable             bool        `json:"isNullable" db:"IS_NULLABLE"`
		DataType               string      `json:"dataType" db:"DATA_TYPE"`
		CharacterMaximumLength *uint       `json:"characterMaximumLength" db:"CHARACTER_MAXIMUM_LENGTH"`
		CharacterOctetLength   *uint       `json:"characterOctetLength" db:"CHARACTER_OCTET_LENGTH"`
		NumericPrecision       *uint       `json:"numericPrecision" db:"NUMERIC_PRECISION"`
		NumericScale           *uint       `json:"numericScale" db:"NUMERIC_SCALE"`
		DatetimePrecision      *uint       `json:"datetimePrecision" db:"DATETIME_PRECISION"`
		CharacterSetName       *string     `json:"characterSetName" db:"CHARACTER_SET_NAME"`
		CollationName          *string     `json:"collationName" db:"COLLATION_NAME"`
		ColumnType             string      `json:"columnType" db:"COLUMN_TYPE"`
		ColumnKey              string      `json:"columnKey" db:"COLUMN_KEY"`
		Extra                  string      `json:"extra" db:"EXTRA"`
		Privileges             string      `json:"privileges" db:"PRIVILEGES"`
		ColumnComment          string      `json:"columnComment" db:"COLUMN_COMMENT"`
	}
	MysqlIndex struct {
		TableCatalog string  `db:"TABLE_CATALOG"`
		TableSchema  string  `db:"TABLE_SCHEMA"`
		TableName    string  `db:"TABLE_NAME"`
		NonUnique    uint    `db:"NON_UNIQUE"`
		IndexSchema  string  `db:"INDEX_SCHEMA"`
		IndexName    string  `db:"INDEX_NAME"`
		SeqInIndex   uint    `db:"SEQ_IN_INDEX"`
		ColumnName   string  `db:"COLUMN_NAME"`
		Collation    string  `db:"COLLATION"`
		Cardinality  uint    `db:"CARDINALITY"`
		SubPart      *string `db:"SUB_PART"`
		Packed       *string `db:"PACKED"`
		Nullable     bool    `db:"NULLABLE"`
		IndexType    string  `db:"INDEX_TYPE"`
		Comment      string  `db:"COMMENT"`
		IndexComment string  `db:"INDEX_COMMENT"`
	}
	dataTypeRange struct {
		Min  int64
		Max  uint64
		Null bool
	}
)

func NewConnection(user, password, dbname string, close bool) (*MysqlConnection, error) {
	con := MysqlConnection{
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

func (con *MysqlConnection) Open() error {
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

func (con *MysqlConnection) Close() error {
	if con.db == nil {
		return nil
	}
	return con.db.Close()
}

func (con *MysqlConnection) GetTables() ([]string, error) {

	if con.db == nil {
		return nil, errors.New("database is closed")
	}
	rows, err := con.db.Query("SHOW TABLES")
	if err != nil {
		panic(err.Error())
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

func (con *MysqlConnection) GetColumns(tname string) ([]MysqlColumn, error) {

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
		panic(err.Error())
	}

	var tableCatalog, tableSchema, tableName, columnName, isNullable, dataType,
		columnType, columnKey, extra, privileges,
		columnComment string
	var characterSetName, collationName sql.NullString
	var ordinalPosition uint
	var characterMaximumLength, characterOctetLength, numericPrecision,
		numericScale, datetimePrecision sql.NullInt64
	var columnDefault interface{}

	result := []MysqlColumn{}
	for rows.Next() {
		err = rows.Scan(
			&tableCatalog, &tableSchema, &tableName, &columnName, &ordinalPosition,
			&columnDefault, &isNullable, &dataType, &characterMaximumLength,
			&characterOctetLength, &numericPrecision, &numericScale,
			&datetimePrecision, &characterSetName, &collationName, &columnType,
			&columnKey, &extra, &privileges, &columnComment,
		)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, MysqlColumn{
			TableCatalog:           tableCatalog,
			TableSchema:            tableSchema,
			TableName:              tableName,
			ColumnName:             columnName,
			OrdinalPosition:        ordinalPosition,
			ColumnDefault:          columnDefault,
			IsNullable:             isNullable == "YES",
			DataType:               dataType,
			CharacterMaximumLength: parseIntPointer(&characterMaximumLength),
			CharacterOctetLength:   parseIntPointer(&characterOctetLength),
			NumericPrecision:       parseIntPointer(&numericPrecision),
			NumericScale:           parseIntPointer(&numericScale),
			DatetimePrecision:      parseIntPointer(&datetimePrecision),
			CharacterSetName:       parseStringPointer(&characterSetName),
			CollationName:          parseStringPointer(&collationName),
			ColumnType:             columnType,
			ColumnKey:              columnKey,
			Extra:                  extra,
			Privileges:             privileges,
			ColumnComment:          columnComment,
		})
	}
	return result, nil
}

func (con *MysqlConnection) GetIndexs(tname string) ([]MysqlIndex, error) {

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
		panic(err.Error())
	}

	var tableCatalog, tableSchema, tableName, indexSchema, indexName, columnName, collation, nullable,
		indexType, comment, indexComment string
	var subPart, packed sql.NullString
	var nonUnique, seqInIndex, cardinality, sortNumber uint

	result := []MysqlIndex{}
	for rows.Next() {
		err = rows.Scan(
			&tableCatalog, &tableSchema, &tableName, &nonUnique, &indexSchema, &indexName, &seqInIndex,
			&columnName, &collation, &cardinality, &subPart, &packed, &nullable, &indexType, &comment,
			&indexComment, &sortNumber,
		)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, MysqlIndex{
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
			SubPart:      parseStringPointer(&subPart),
			Packed:       parseStringPointer(&packed),
			Nullable:     nullable == "YES",
			IndexType:    indexType,
			Comment:      comment,
			IndexComment: indexComment,
		})
	}
	return result, nil
}

func (mc MysqlColumn) Unsigned() bool {
	return strings.Index(mc.ColumnType, "unsigned") >= 0
}

func (mc MysqlColumn) AutoIncrement() bool {
	return strings.Index(mc.Extra, "auto_increment") >= 0
}

func (mc MysqlColumn) Primary() bool {
	return mc.ColumnKey == "PRI"
}

func (mc MysqlColumn) Unique() bool {
	return mc.ColumnKey == "UNI"
}

func (mc MysqlColumn) DataTypeRange() dataTypeRange {
	unsigned := mc.Unsigned()
	switch mc.DataType {
	case "char", "varchar", "enum", "set":
		var max uint64
		if mc.CharacterMaximumLength != nil {
			max = uint64(*mc.CharacterMaximumLength)
		}
		return dataTypeRange{Min: 0, Max: max, Null: mc.IsNullable}
	case "tinyint":
		if unsigned {
			return dataTypeRange{Min: 0, Max: 255, Null: mc.IsNullable} // "uint8" // 0 to 255
		}
		return dataTypeRange{Min: -128, Max: 127, Null: mc.IsNullable} // "int8" // -128 to 127
	case "smallint":
		if unsigned {
			return dataTypeRange{Min: 0, Max: 65535, Null: mc.IsNullable} // "uint16" // 0 to 65535
		}
		return dataTypeRange{Min: -32768, Max: 32767, Null: mc.IsNullable} // "int16" // -32768 to 32767
	case "mediumint":
		if unsigned {
			return dataTypeRange{Min: 0, Max: 16777215, Null: mc.IsNullable} // "uint32" // unsigned: 0 to 16777215
		}
		return dataTypeRange{Min: -8388608, Max: 8388607, Null: mc.IsNullable} // "int32" // signed: -8388608 to 8388607
	case "bigint":
		if unsigned {
			return dataTypeRange{Min: 0, Max: 18446744073709551615, Null: mc.IsNullable} // "uint64" // unsigned: 0 to 16777215
		}
		return dataTypeRange{Min: -9223372036854775808, Max: 9223372036854775807, Null: mc.IsNullable} // "int64" // signed: -8388608 to 8388607
	case "int", "integer":
		if unsigned {
			return dataTypeRange{Min: 0, Max: 4294967295, Null: mc.IsNullable} // "uint32" // 0 to 4294967295
		}
		return dataTypeRange{Min: -2147483648, Max: 2147483647, Null: mc.IsNullable} // "int32" // -2147483648 to 2147483647
	case "float":
		return dataTypeRange{Null: mc.IsNullable}
	case "double", "decimal", "dec":
		return dataTypeRange{Null: mc.IsNullable}
	case "date", "datetime", "timestamp", "time":
		return dataTypeRange{Null: mc.IsNullable}
	default:
		return dataTypeRange{Null: mc.IsNullable}
	}
}
