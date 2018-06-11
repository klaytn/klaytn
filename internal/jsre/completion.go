package jsre

import (
	"github.com/robertkrimen/otto"
	"sort"
	"strings"
)

// CompleteKeywords returns potential continuations for the given line. Since line is
// evaluated, callers need to make sure that evaluating line does not have side effects.
func (jsre *JSRE) CompleteKeywords(line string) []string {
	var results []string
	jsre.Do(func(vm *otto.Otto) {
		results = getCompletions(vm, line)
	})
	return results
}

func getCompletions(vm *otto.Otto, line string) (results []string) {
	parts := strings.Split(line, ".")
	objRef := "this"
	prefix := line
	if len(parts) > 1 {
		objRef = strings.Join(parts[0:len(parts)-1], ".")
		prefix = parts[len(parts)-1]
	}

	obj, _ := vm.Object(objRef)
	if obj == nil {
		return nil
	}
	iterOwnAndConstructorKeys(vm, obj, func(k string) {
		if strings.HasPrefix(k, prefix) {
			if objRef == "this" {
				results = append(results, k)
			} else {
				results = append(results, strings.Join(parts[:len(parts)-1], ".")+"."+k)
			}
		}
	})

	// Append opening parenthesis (for functions) or dot (for objects)
	// if the line itself is the only completion.
	if len(results) == 1 && results[0] == line {
		obj, _ := vm.Object(line)
		if obj != nil {
			if obj.Class() == "Function" {
				results[0] += "("
			} else {
				results[0] += "."
			}
		}
	}

	sort.Strings(results)
	return results
}
