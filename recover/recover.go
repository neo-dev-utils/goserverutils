package recover

import (
	// "common/log"
	"fmt"
)

// Recover 回收
func RecoverHandle(v ...interface{}) {
	if err := recover(); err != nil {
		if len(v) > 0 {
			s := fmt.Sprintln(v...)
			fmt.Println("s : ", s, " err : ", err)
			// log.Error("recover err %s : ", err)
			// log.Error("%s\n%s\n", s, string(debug.Stack()))
		} else {
			fmt.Println("v : ", v)
			// log.Error("%s\n", string(debug.Stack()))
		}
	}
}
