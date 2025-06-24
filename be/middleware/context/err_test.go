package context

import (
	"fmt"
	"strings"
	"testing"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/7/7
    @desc:

***************************/

func TestErr_Error(t *testing.T) {
	outputs := []string{
		"middleware/errors/error.go:",
		"middleware/errors/error.go:",
		"middleware/context/err.go:",
		"middleware/context/err_test.go:",
		"go/src/testing/testing.go:",
	}
	appErr := (&err{}).Error(-1, fmt.Errorf("app error"))
	callers := appErr.Caller()
	for idx, caller := range callers {
		if idx >= len(outputs) {
			break
		}
		if !strings.HasPrefix(caller, outputs[idx]) {
			t.Errorf("index: %d, atcual: %s, expect: %s", idx, caller, outputs[idx])
		}
	}

}
