package sql

import (
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestLongStr(t *testing.T) {
	str := `asdjasdjash ajshd hajshd jashd ahsdh ajshhajsh hjashd hjkasdhj as
            yqwue uasdh asd jqweh jhqwjeh jqhwje hjqhwjeh jhqwje hjqhw ehqwe io`
	fmt.Println(str)
}
