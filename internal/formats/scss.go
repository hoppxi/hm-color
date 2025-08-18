package formats

import "fmt"

func FormatSCSS(colors map[string]string) string {
	out := ""
	for k, v := range colors {
		out += fmt.Sprintf("$%s: %s;\n", k, v)
	}
	return out
}
