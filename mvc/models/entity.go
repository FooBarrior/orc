package models

import (
    "github.com/orc/db"
)

type ModelManager struct{}

type Entity struct {
    TableName string
    Caption   string
    Fields    []map[string]string
    Columns   []string
    ColNames  []string
    Ref       bool
    RefData   map[string]interface{}
    RefFields []string
    Sub       bool
    SubTable  []string
    SubField  string
}

func (this Entity) Create() {
    db.QueryCreateTable(this.TableName, this.Fields)
}

func (this Entity) Select(where []string, condition string, fields []string) ([]interface{}, map[string]interface{}) {
    result1 := db.Select(this.TableName, where, condition, fields)
    if this.Ref {
        result2 := this.RefData
        return result1, result2
    }
    return result1, nil
}

func (this Entity) Insert(fields []string, params []interface{}) {
    db.QueryInsert(this.TableName, fields, params)
}

func (this Entity) Update(fields []string, params []interface{}, where string) {
    db.QueryUpdate(this.TableName, where, fields, params)
}

func (this Entity) Delete(field string, params []interface{}) {
    db.QueryDelete(this.TableName, field, params)
}

func (this Entity) GetSubTable(index int) string {
    return this.SubTable[index]
}

func (this Entity) GetSubField() string {
    return this.SubField
}

func (this Entity) GetColumns() []string {
    return this.Columns
}

func (this Entity) GetColumnByIdx(index int) string {
    return this.Columns[index]
}

func (this Entity) GetColumnSlice(index int) []string {
    return this.Columns[index:]
}

func (this Entity) GetColNames() []string {
    return this.ColNames
}

func (this Entity) GetTableName() string {
    return this.TableName
}

func (this Entity) GetCaption() string {
    return this.Caption
}

func (this Entity) GetRefFields() []string {
    return this.RefFields
}

func (this Entity) GetSub() bool {
    return this.Sub
}

type VirtEntity interface {
    GetTableName() string
    GetCaption() string
    GetSub() bool
    GetSubTable(index int) string
    GetSubField() string
    GetColumns() []string
    GetColumnByIdx(index int) string
    GetColumnSlice(index int) []string
    GetColNames() []string
    GetRefFields() []string
    Create()
    Select(where []string, condition string, fields []string) ([]interface{}, map[string]interface{})
    Insert(fields []string, params []interface{})
    Update(fields []string, params []interface{}, where string)
    Delete(field string, params []interface{})
}