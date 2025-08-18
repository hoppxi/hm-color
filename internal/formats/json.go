package formats

import (
	"encoding/json"
)

func FormatJSON(colors map[string]string) string {
	data, _ := json.MarshalIndent(colors, "", "  ")
	return string(data)
}
