package utils

import "encoding/json"

//IsJSON error ij json invalid
func IsJSON(s string) error {
	var js struct{}

	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return err
	}
	return nil
}
