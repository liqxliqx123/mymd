package pariticiple_query

import (
	"github.com/alecthomas/participle/v2"
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
