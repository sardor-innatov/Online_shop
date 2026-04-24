package helper

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonObject map[string]any

// Value вызывается при записи в БД (INSERT/UPDATE)
func (j JsonObject) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil // Сохранит NULL в базе
	}
	return json.Marshal(j)
}

// Scan вызывается при чтении из БД (SELECT)
func (j *JsonObject) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported data type: %T", value)
	}

	// Инициализируем карту перед заполнением
	if *j == nil {
		*j = make(JsonObject)
	}
	
	return json.Unmarshal(bytes, j)
}
