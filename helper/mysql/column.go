package mysql

import (
	"strings"
)

type (
	Column struct {
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
	dataTypeRange struct {
		Min  int64
		Max  uint64
		Null bool
	}
)

func (mc Column) Unsigned() bool {
	return strings.Contains(mc.ColumnType, "unsigned")
}

func (mc Column) AutoIncrement() bool {
	return strings.Contains(mc.Extra, "auto_increment")
}

func (mc Column) Primary() bool {
	return mc.ColumnKey == "PRI"
}

func (mc Column) Unique() bool {
	return mc.ColumnKey == "UNI"
}

func (mc Column) DataTypeRange() dataTypeRange {
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
