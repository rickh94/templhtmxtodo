package components

import "fmt"

func HxCsrfHeader(csrf string) string {
	return fmt.Sprintf("{\"X-Csrf-Token\": \"%s\"}", csrf)
}
