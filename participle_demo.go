package main

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/participle/v2"
	"log"
)

var parser = participle.MustBuild[Root]()

type Root struct {
	Left  *Operand   `"{" @@ "}"`
	Right []*Copilot `@@*`
}

type Copilot struct {
	Op    string   `@( "AND" | "OR" | "NOT")`
	Right *Operand `"{" @@ "}"`
}

type Operand struct {
	Field      string      ` "[" @Ident "]"`
	Sentence   string      `@"sentence"?`
	Unit       string      `@"unit"?`
	SubOperand *SubOperand `@@`
}

type SubOperand struct {
	Left  ValueWithInterval `@@`
	Right []*Expression     `@@*`
}

type ValueWithInterval struct {
	Left     Value `@@`
	Interval int   `( "-" @Int )?`
}

type Expression struct {
	Op    string            `@( "AND" | "OR" | "NOT")`
	Right ValueWithInterval `@@`
}

type Value struct {
	Str           string      `@Ident`
	SubExpression *SubOperand `| "(" @@ ")"`
}

func main() {
	tests := []string{
		//`{[bb] unit(((小A  AND 小B)-4 OR (小M AND 小N)) NOT (大C OR 大D))}`,
		//`{[aa]((小A AND 小B)-4 OR (小W AND 小N)) NOT (大C OR 大D)}`,
		//`{[aa] A} AND {[bb] B}`,
		//`{[标题] sentence((X AND Y)-5 OR (A AND B)-4)}`,
		//`[标题] sentence((X AND Y)-5 OR (A AND B)-4)`,
		//`[标题] sentence((X AND Y)-5 OR (A AND B))`,
		//`[标题] sentence((X AND Y)-5)`,
		//`[标题] sentence((X AND Y))`,
		//`[标题] sentence((X AND Y) NOT (A OR B))`,
		//`[标题] sentence(X AND Y)`,
		//`[标题] unit(X AND Y)`,
		//`[标题] X`,
		//`[标题] X AND Y`,
		//`[标题] X OR Y NOT Z`,
		//`[标题] X OR (Y NOT Z)`,
		//`[标题] ( X OR Y ) AND (A NOT B)`,
		//`[标题] X OR (Y AND A NOT B)`,
		`{[aa]((小A AND 小B)-4 OR (小W AND 小N)) NOT (大C OR 大D)} AND {[bb] unit(((小A  AND 小B)-4 OR (小M AND 小N)) NOT (大C OR 大D))}`,
	}

	for _, testString := range tests {
		expr, err := parser.ParseString("", testString)
		if err != nil {
			log.Fatalf("Error parsing: %v", err)
		}
		data, err := json.MarshalIndent(expr, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling to JSON: %v", err)
		}
		fmt.Println(string(data))
	}
}
