package pkg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, s); err != nil {
			return fmt.Errorf("failed to unmarshal string array: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(v), s); err != nil {
			return fmt.Errorf("failed to unmarshal string array: %w", err)
		}
	default:
		return fmt.Errorf("failed to scan string array: unsupported type %T", value)
	}
	return nil
}

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal string array: %w", err)
	}
	return string(jsonBytes), nil
}
