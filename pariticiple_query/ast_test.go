package pariticiple_query

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParticiple(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    string
		wantErr error
	}{
		{
			name:    "{[account_name] 星}",
			expr:    "{[account_name] 星}",
			want:    `{"query":{"bool":{"must":[{"wildcard":{"account_name":{"value":"*星*"}}}]}}}`,
			wantErr: nil,
		},
		{
			name:    "{[account_name] 星 AND 光}",
			expr:    "{[account_name] 星 AND 光}",
			want:    `{"query":{"bool":{"must":[{"bool":{"must":[{"wildcard":{"account_name":{"value":"*星*"}}},{"wildcard":{"account_name":{"value":"*光*"}}}]}}]}}}`,
			wantErr: nil,
		},
		{
			name:    "{[account_name] 星 OR 光}",
			expr:    "{[account_name] 星 OR 光}",
			want:    `{"query":{"bool":{"should":[{"match":{"account_name":"星"}},{"match":{"account_name":"光"}}]}}}`,
			wantErr: nil,
		},
		{
			name:    "{[account_name] 星 NOT 光}",
			expr:    "{[account_name] 星 NOT 光}",
			want:    `{"query":{"bool":{"must":[{"match":{"account_name":"星"}}],"must_not":[{"match":{"account_name":"光"}}]}}}`,
			wantErr: nil,
		},
		{
			name:    "{[account_name] 星 OR 光 NOT 复}",
			expr:    "{[account_name] 星 OR 光 NOT 复}",
			want:    `{"query":{"bool":{"should":[{"match":{"account_name":"星"}},{"match":{"account_name":"光"}}],"must_not":[{"match":{"account_name":"复"}}]}}}`,
			wantErr: nil,
		},
	}
	//tests := []string{
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

	//`{[account_name] 星}`,
	//`{[account_name] 星 AND 光}`,
	//`{[account_name] 星 OR 光}`,
	//`{[account_name] 星 NOT 光}`,
	//`{[account_name] 星 OR 光 NOT 复}`,
	//`[标题] X OR (Y NOT Z)`,
	//`[标题] ( X OR Y ) AND (A NOT B)`,
	//`[标题] X OR (Y AND A NOT B)`,
	//`{[aa]((小A AND 小B)-4 OR (小W AND 小N)) NOT (大C OR 大D)} AND {[bb] unit(((小A  AND 小B)-4 OR (小M AND 小N)) NOT (大C OR 大D))}`,
	//`{[headline.kw] 星}`,
	//}
	for _, tt := range tests {
		_, err := parser.ParseString("", tt.expr)
		if err != nil {
			log.Fatalf("Error parsing: %v", err)
		}
		//dsl := expr.ToESQuery()
		//marshal, err := json.Marshal(dsl)
		//if err != nil {
		//	log.Fatalf(err.Error())
		//}
		//fmt.Println(string(marshal))
		assert.NoError(t, err)
		//assert.Equal(t, tt.want, string(marshal))
	}
}
