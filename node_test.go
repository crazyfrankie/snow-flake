package snowflake

import (
	"fmt"
	"testing"
)

func BenchmarkNode(b *testing.B) {
	node, err := NewNode(1)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		code := node.GenerateCode()
		fmt.Println(code)
	}
}
