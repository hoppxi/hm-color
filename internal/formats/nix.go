package formats

import "fmt"

func FormatNix(colors map[string]string) string {
	out := "{\n"
	for k, v := range colors {
		out += fmt.Sprintf("  %s = \"%s\";\n", k, v)
	}
	out += "}\n"
	return out
}
