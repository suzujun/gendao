package mysql

type (
	Index struct {
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
)
