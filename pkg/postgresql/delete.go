package postgresql

import (
	"fmt"
)

func (t *transaction) delete() {
	t.WG.Add(1)
	defer t.WG.Done()
	for {
		record, more := <-t.DeleteChan
		if !more {
			break
		}
		if record != "" {
			fmt.Print("-")
		}
	}
}
