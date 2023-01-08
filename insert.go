package melody

import (
	"errors"
	"fmt"
	"strings"
)

type InsertBuilder struct {
	ctx insertContext
}

type insertContext struct {
	Table            string
	Value            []insertValue
	withDuplicateKey bool
}

type insertValue struct {
	Column           string
	Value            interface{}
	withDuplicateKey bool
}

func NewInsert(table string) *InsertBuilder {
	return &InsertBuilder{
		insertContext{
			Table: table,
		},
	}
}

func (i *InsertBuilder) Set(column string, value interface{}) *InsertBuilder {
	i.ctx.Value = append(i.ctx.Value, insertValue{Column: column, Value: value})
	return i
}

func (i *InsertBuilder) UpdateDuplicateKey() *InsertBuilder {
	i.ctx.withDuplicateKey = true
	i.ctx.Value[len(i.ctx.Value)-1].withDuplicateKey = true
	return i
}

func (i *InsertBuilder) Get() (query string, params []interface{}, err error) {
	return i.build()
}

func (i *InsertBuilder) build() (res string, params []interface{}, err error) {
	var result []string

	if i.ctx.Table == "" {
		return res, params, errors.New("one table need to be defined")
	}

	result = append(result, fmt.Sprintf("INSERT INTO %s SET", i.ctx.Table))

	var resultValue []string
	for _, v := range i.ctx.Value {
		resultValue = append(resultValue, fmt.Sprintf("%s = ?", v.Column))
		params = append(params, v.Value)
	}

	result = append(result, strings.Join(resultValue, ", "))

	if i.ctx.withDuplicateKey {
		var resultOnUpdate []string
		for _, v := range i.ctx.Value {
			if v.withDuplicateKey {
				resultOnUpdate = append(resultOnUpdate, fmt.Sprintf("%s = ?", v.Column))
				params = append(params, v.Value)
			}
		}

		result = append(result, "ON DUPLICATE KEY UPDATE")
		result = append(result, strings.Join(resultOnUpdate, ", "))
	}

	return strings.Join(result, " "), params, nil
}
