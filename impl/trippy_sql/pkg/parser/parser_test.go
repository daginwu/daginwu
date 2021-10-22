package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name     string
	SQL      string
	Expected AST
	Err      error
}

type output struct {
	NoErrorExamples []testCase
	ErrorExamples   []testCase
	Types           []string
	Operators       []string
}

func TestSQL(t *testing.T) {
	ts := []testCase{
		{
			Name:     "empty query fails",
			SQL:      "",
			Expected: AST{},
			Err:      fmt.Errorf("query type cannot be empty"),
		},
		{
			Name:     "SELECT without FROM fails",
			SQL:      "SELECT",
			Expected: AST{Type: Select},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "SELECT without fields fails",
			SQL:      "SELECT FROM 'a'",
			Expected: AST{Type: Select},
			Err:      fmt.Errorf("at SELECT: expected field to SELECT"),
		},
		{
			Name:     "SELECT with comma and empty field fails",
			SQL:      "SELECT b, FROM 'a'",
			Expected: AST{Type: Select},
			Err:      fmt.Errorf("at SELECT: expected field to SELECT"),
		},
		{
			Name:     "SELECT works",
			SQL:      "SELECT a FROM 'b'",
			Expected: AST{Type: Select, TableName: "b", Fields: []string{"a"}},
			Err:      nil,
		},
		{
			Name:     "SELECT works with lowercase",
			SQL:      "select a fRoM 'b'",
			Expected: AST{Type: Select, TableName: "b", Fields: []string{"a"}},
			Err:      nil,
		},
		{
			Name:     "SELECT many fields works",
			SQL:      "SELECT a, c, d FROM 'b'",
			Expected: AST{Type: Select, TableName: "b", Fields: []string{"a", "c", "d"}},
			Err:      nil,
		},
		{
			Name: "SELECT with alias works",
			SQL:  "SELECT a as z, b as y, c FROM 'b'",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "b", "c"},
				Aliases: map[string]string{
					"a": "z",
					"b": "y",
				},
			},
			Err: nil,
		},

		{
			Name:     "SELECT with empty WHERE fails",
			SQL:      "SELECT a, c, d FROM 'b' WHERE",
			Expected: AST{Type: Select, TableName: "b", Fields: []string{"a", "c", "d"}},
			Err:      fmt.Errorf("at WHERE: empty WHERE clause"),
		},
		{
			Name:     "SELECT with WHERE with only operand fails",
			SQL:      "SELECT a, c, d FROM 'b' WHERE a",
			Expected: AST{Type: Select, TableName: "b", Fields: []string{"a", "c", "d"}},
			Err:      fmt.Errorf("at WHERE: condition without operator"),
		},
		{
			Name: "SELECT with WHERE with = works",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a = ''",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Eq, Operand2: "", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with < works",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a < '1'",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Lt, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with <= works",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a <= '1'",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Lte, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with > works",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a > '1'",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Gt, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with >= works",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a >= '1'",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Gte, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with != works",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a != '1'",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Ne, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with != works (comparing field against another field)",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a != b",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Ne, Operand2: "b", Operand2IsField: true},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT * works",
			SQL:  "SELECT * FROM 'b'",
			Expected: AST{
				Type:       Select,
				TableName:  "b",
				Fields:     []string{"*"},
				Conditions: nil,
			},
			Err: nil,
		},
		{
			Name: "SELECT a, * works",
			SQL:  "SELECT a, * FROM 'b'",
			Expected: AST{
				Type:       Select,
				TableName:  "b",
				Fields:     []string{"a", "*"},
				Conditions: nil,
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with two conditions using AND works",
			SQL:  "SELECT a, c, d FROM 'b' WHERE a != '1' AND b = '2'",
			Expected: AST{
				Type:      Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Ne, Operand2: "1", Operand2IsField: false},
					{Operand1: "b", Operand1IsField: true, Operator: Eq, Operand2: "2", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name:     "Empty UPDATE fails",
			SQL:      "UPDATE",
			Expected: AST{},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "Incomplete UPDATE with table name fails",
			SQL:      "UPDATE 'a'",
			Expected: AST{},
			Err:      fmt.Errorf("at WHERE: WHERE clause is mandatory for UPDATE & DELETE"),
		},
		{
			Name:     "Incomplete UPDATE with table name and SET fails",
			SQL:      "UPDATE 'a' SET",
			Expected: AST{},
			Err:      fmt.Errorf("at WHERE: WHERE clause is mandatory for UPDATE & DELETE"),
		},
		{
			Name:     "Incomplete UPDATE with table name, SET with a field but no value and WHERE fails",
			SQL:      "UPDATE 'a' SET b WHERE",
			Expected: AST{},
			Err:      fmt.Errorf("at UPDATE: expected '='"),
		},
		{
			Name:     "Incomplete UPDATE with table name, SET with a field and = but no value and WHERE fails",
			SQL:      "UPDATE 'a' SET b = WHERE",
			Expected: AST{},
			Err:      fmt.Errorf("at UPDATE: expected quoted value"),
		},
		{
			Name:     "Incomplete UPDATE due to no WHERE clause fails",
			SQL:      "UPDATE 'a' SET b = 'hello' WHERE",
			Expected: AST{},
			Err:      fmt.Errorf("at WHERE: empty WHERE clause"),
		},
		{
			Name:     "Incomplete UPDATE due incomplete WHERE clause fails",
			SQL:      "UPDATE 'a' SET b = 'hello' WHERE a",
			Expected: AST{},
			Err:      fmt.Errorf("at WHERE: condition without operator"),
		},
		{
			Name: "UPDATE works",
			SQL:  "UPDATE 'a' SET b = 'hello' WHERE a = '1'",
			Expected: AST{
				Type:      Update,
				TableName: "a",
				Updates:   map[string]string{"b": "hello"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Eq, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "UPDATE works with simple quote inside",
			SQL:  "UPDATE 'a' SET b = 'hello\\'world' WHERE a = '1'",
			Expected: AST{
				Type:      Update,
				TableName: "a",
				Updates:   map[string]string{"b": "hello\\'world"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Eq, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "UPDATE with multiple SETs works",
			SQL:  "UPDATE 'a' SET b = 'hello', c = 'bye' WHERE a = '1'",
			Expected: AST{
				Type:      Update,
				TableName: "a",
				Updates:   map[string]string{"b": "hello", "c": "bye"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Eq, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name: "UPDATE with multiple SETs and multiple conditions works",
			SQL:  "UPDATE 'a' SET b = 'hello', c = 'bye' WHERE a = '1' AND b = '789'",
			Expected: AST{
				Type:      Update,
				TableName: "a",
				Updates:   map[string]string{"b": "hello", "c": "bye"},
				Conditions: []Condition{
					{Operand1: "a", Operand1IsField: true, Operator: Eq, Operand2: "1", Operand2IsField: false},
					{Operand1: "b", Operand1IsField: true, Operator: Eq, Operand2: "789", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name:     "Empty DELETE fails",
			SQL:      "DELETE FROM",
			Expected: AST{},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "DELETE without WHERE fails",
			SQL:      "DELETE FROM 'a'",
			Expected: AST{},
			Err:      fmt.Errorf("at WHERE: WHERE clause is mandatory for UPDATE & DELETE"),
		},
		{
			Name:     "DELETE with empty WHERE fails",
			SQL:      "DELETE FROM 'a' WHERE",
			Expected: AST{},
			Err:      fmt.Errorf("at WHERE: empty WHERE clause"),
		},
		{
			Name:     "DELETE with WHERE with field but no operator fails",
			SQL:      "DELETE FROM 'a' WHERE b",
			Expected: AST{},
			Err:      fmt.Errorf("at WHERE: condition without operator"),
		},
		{
			Name: "DELETE with WHERE works",
			SQL:  "DELETE FROM 'a' WHERE b = '1'",
			Expected: AST{
				Type:      Delete,
				TableName: "a",
				Conditions: []Condition{
					{Operand1: "b", Operand1IsField: true, Operator: Eq, Operand2: "1", Operand2IsField: false},
				},
			},
			Err: nil,
		},
		{
			Name:     "Empty INSERT fails",
			SQL:      "INSERT INTO",
			Expected: AST{},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "INSERT with no rows to insert fails",
			SQL:      "INSERT INTO 'a'",
			Expected: AST{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails",
			SQL:      "INSERT INTO 'a' (",
			Expected: AST{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails #2",
			SQL:      "INSERT INTO 'a' (b",
			Expected: AST{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails #3",
			SQL:      "INSERT INTO 'a' (b)",
			Expected: AST{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails #4",
			SQL:      "INSERT INTO 'a' (b) VALUES",
			Expected: AST{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete row fails",
			SQL:      "INSERT INTO 'a' (b) VALUES (",
			Expected: AST{},
			Err:      fmt.Errorf("at INSERT INTO: value count doesn't match field count"),
		},
		{
			Name: "INSERT works",
			SQL:  "INSERT INTO 'a' (b) VALUES ('1')",
			Expected: AST{
				Type:      Insert,
				TableName: "a",
				Fields:    []string{"b"},
				Inserts:   [][]string{{"1"}},
			},
			Err: nil,
		},
		{
			Name:     "INSERT * fails",
			SQL:      "INSERT INTO 'a' (*) VALUES ('1')",
			Expected: AST{},
			Err:      fmt.Errorf("at INSERT INTO: expected at least one field to insert"),
		},
		{
			Name: "INSERT with multiple fields works",
			SQL:  "INSERT INTO 'a' (b,c,    d) VALUES ('1','2' ,  '3' )",
			Expected: AST{
				Type:      Insert,
				TableName: "a",
				Fields:    []string{"b", "c", "d"},
				Inserts:   [][]string{{"1", "2", "3"}},
			},
			Err: nil,
		},
		{
			Name: "INSERT with multiple fields and multiple values works",
			SQL:  "INSERT INTO 'a' (b,c,    d) VALUES ('1','2' ,  '3' ),('4','5' ,'6' )",
			Expected: AST{
				Type:      Insert,
				TableName: "a",
				Fields:    []string{"b", "c", "d"},
				Inserts:   [][]string{{"1", "2", "3"}, {"4", "5", "6"}},
			},
			Err: nil,
		},
	}

	output := output{Types: TypeString, Operators: OperatorString}
	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := ParseMany([]string{tc.SQL})
			if tc.Err != nil && err == nil {
				t.Errorf("Error should have been %v", tc.Err)
			}
			if tc.Err == nil && err != nil {
				t.Errorf("Error should have been nil but was %v", err)
			}
			if tc.Err != nil && err != nil {
				require.Equal(t, tc.Err, err, "Unexpected error")
			}
			if len(actual) > 0 {
				require.Equal(t, tc.Expected, actual[0], "AST didn't match expectation")
			}
			if tc.Err != nil {
				output.ErrorExamples = append(output.ErrorExamples, tc)
			} else {
				output.NoErrorExamples = append(output.NoErrorExamples, tc)
			}
		})
	}

}
