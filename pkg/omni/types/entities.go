// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Event struct {
	Name      string
	Anonymous bool
	Fields    []*Field
	Logs      map[int64]Log // Map of VulcanizeIdLog to parsed event log
}

type Method struct {
	Name    string
	Const   bool
	Inputs  []*Field
	Outputs []*Field
	Results []*Result
}

type Field struct {
	abi.Argument        // Name, Type, Indexed
	PgType       string // Holds type used when committing data held in this field to postgres
}

// Struct to hold results from method call with given inputs across different blocks
type Result struct {
	Inputs  []interface{} // Will only use addresses
	Outputs map[int64]interface{}
	PgType  string // Holds output pg type
}

// Struct to hold event log data data
type Log struct {
	Id     int64             // VulcanizeIdLog
	Values map[string]string // Map of event input names to their values
	Block  int64
	Tx     string
}

// Unpack abi.Event into our custom Event struct
func NewEvent(e abi.Event) *Event {
	fields := make([]*Field, len(e.Inputs))
	for i, input := range e.Inputs {
		fields[i] = &Field{}
		fields[i].Name = input.Name
		fields[i].Type = input.Type
		fields[i].Indexed = input.Indexed
		// Fill in pg type based on abi type
		switch fields[i].Type.T {
		case abi.StringTy, abi.HashTy, abi.AddressTy:
			fields[i].PgType = "CHARACTER VARYING(66)"
		case abi.IntTy, abi.UintTy:
			fields[i].PgType = "DECIMAL"
		case abi.BoolTy:
			fields[i].PgType = "BOOLEAN"
		case abi.BytesTy, abi.FixedBytesTy:
			fields[i].PgType = "BYTEA"
		case abi.ArrayTy:
			fields[i].PgType = "TEXT[]"
		case abi.FixedPointTy:
			fields[i].PgType = "MONEY" // use shopspring/decimal for fixed point numbers in go and money type in postgres?
		case abi.FunctionTy:
			fields[i].PgType = "TEXT"
		default:
			fields[i].PgType = "TEXT"
		}
	}

	return &Event{
		Name:      e.Name,
		Anonymous: e.Anonymous,
		Fields:    fields,
		Logs:      map[int64]Log{},
	}
}

// Unpack abi.Method into our custom Method struct
func NewMethod(m abi.Method) *Method {
	inputs := make([]*Field, len(m.Inputs))
	for i, input := range m.Inputs {
		inputs[i] = &Field{}
		inputs[i].Name = input.Name
		inputs[i].Type = input.Type
		inputs[i].Indexed = input.Indexed
		switch inputs[i].Type.T {
		case abi.StringTy, abi.HashTy, abi.AddressTy:
			inputs[i].PgType = "CHARACTER VARYING(66)"
		case abi.IntTy, abi.UintTy:
			inputs[i].PgType = "DECIMAL"
		case abi.BoolTy:
			inputs[i].PgType = "BOOLEAN"
		case abi.BytesTy, abi.FixedBytesTy:
			inputs[i].PgType = "BYTEA"
		case abi.ArrayTy:
			inputs[i].PgType = "TEXT[]"
		case abi.FixedPointTy:
			inputs[i].PgType = "MONEY" // use shopspring/decimal for fixed point numbers in go and money type in postgres?
		case abi.FunctionTy:
			inputs[i].PgType = "TEXT"
		default:
			inputs[i].PgType = "TEXT"
		}
	}

	outputs := make([]*Field, len(m.Outputs))
	for i, output := range m.Outputs {
		outputs[i] = &Field{}
		outputs[i].Name = output.Name
		outputs[i].Type = output.Type
		outputs[i].Indexed = output.Indexed
		switch outputs[i].Type.T {
		case abi.StringTy, abi.HashTy, abi.AddressTy:
			outputs[i].PgType = "CHARACTER VARYING(66)"
		case abi.IntTy, abi.UintTy:
			outputs[i].PgType = "DECIMAL"
		case abi.BoolTy:
			outputs[i].PgType = "BOOLEAN"
		case abi.BytesTy, abi.FixedBytesTy:
			outputs[i].PgType = "BYTEA"
		case abi.ArrayTy:
			outputs[i].PgType = "TEXT[]"
		case abi.FixedPointTy:
			outputs[i].PgType = "MONEY" // use shopspring/decimal for fixed point numbers in go and money type in postgres?
		case abi.FunctionTy:
			outputs[i].PgType = "TEXT"
		default:
			outputs[i].PgType = "TEXT"
		}
	}

	return &Method{
		Name:    m.Name,
		Const:   m.Const,
		Inputs:  inputs,
		Outputs: outputs,
		Results: make([]*Result, 0),
	}
}

func (e Event) Sig() string {
	types := make([]string, len(e.Fields))

	for i, input := range e.Fields {
		types[i] = input.Type.String()
	}

	return fmt.Sprintf("%v(%v)", e.Name, strings.Join(types, ","))
}

func (m Method) Sig() string {
	types := make([]string, len(m.Inputs))
	i := 0
	for _, input := range m.Inputs {
		types[i] = input.Type.String()
		i++
	}

	return fmt.Sprintf("%v(%v)", m.Name, strings.Join(types, ","))
}