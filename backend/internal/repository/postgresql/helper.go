package postgresql

import (
	"fmt"
	"strings"
)

func updateQuery(table, setColumn string, columns []string) *strings.Builder {
	var q strings.Builder

	fmt.Fprintf(&q, "update %s set ", table)

	cnt := 1
	for i, c := range columns {
		fmt.Fprintf(&q, "%s = $%d", c, cnt)
		cnt++

		if i != len(columns)-1 {
			q.WriteString(", ")
		}
	}

	fmt.Fprintf(&q, " where %s = $%d", setColumn, cnt)

	return &q
}
