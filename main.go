package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/xerrors"
)

// Dao 层并不知道 sql.ErrNoRows 对于业务意味着什么也无法处理，所以需要上抛给业务层
// Dao 层可以自行定义一个错误类型，除了 err 成员变量外，还可以包含 sql query 语句等

var (
	_ error           = &sqlError{}
	_ xerrors.Wrapper = &sqlError{}
)

type sqlError struct {
	query string
	err   error
}

func (e *sqlError) Unwrap() error {
	return e.err
}

func (e *sqlError) Error() string {
	return fmt.Sprintf("err: %s, query: %s", e.err.Error(), e.query)
}

// sqlQueryReturnInvalidQuery query 语句非法场景
func sqlQueryReturnInvalidQuery() error {
	// 模拟 query 语句非法
	query := "sqlQueryReturnInvalidQuery"
	// err := checkQuery(query)
	// 一般判断 query 是否合法都是 dao 层的逻辑，所以这里 err string 前添加包名 dao
	return &sqlError{query: query, err: errors.New("dao: invalid query")}
}

// sqlQueryReturnNoRows sql.ErrNoRows 场景
func sqlQueryReturnNoRows() error {
	// 模拟 Dao sql 操作返回了 sql.ErrNoRows
	// sql 语句
	query := "sqlQueryReturnNoRows"
	//_, err := db.Query(query)
	err := sql.ErrNoRows
	return &sqlError{query: query, err: err}
}

// sqlQueryReturnNoRowsWithWrap sql.ErrNoRows 场景且 wrap 错误
func sqlQueryReturnNoRowsWithWrap() error {
	// 模拟 Dao sql 操作返回了 sql.ErrNoRows
	// sql 语句
	query := "sqlQueryReturnNoRowsWithWrap"
	//_, err := db.Query(query)
	err := sql.ErrNoRows
	return errors.Wrap(&sqlError{query: query, err: err}, "wrapped err")
}

// 对于业务层来说，Dao 层返回的 sql.ErrNoRows 可能可以忽略（比如：账号没有订单）也可能不能忽略（比如：根据账号 ID 获取不到账号基本信息）
// 所以业务层必须要对 Dao 层返回的错误进行校验
// 如果是 sql.ErrNoRows，可以忽略的情况下继续正常流程，不能忽略的情况下进入异常流程

func isSqlErrNoRows(err error) bool {
	fmt.Printf("this err is %T %+v.\n\n", err, err)
	if errors.Is(err, sql.ErrNoRows) {
		return true
	} else {
		return false
	}
}

// wrap 错误可以帮助开发者快速定位报错地点。
// 如果业务层有多个地方可能会返回 sql.ErrNoRows 或者其他错误，上抛错误时最好 wrap 一下
// 否则也可以不添加，便于日志的浏览

func main() {
	err1 := sqlQueryReturnNoRows()
	isSqlErrNoRows(err1)

	//if isSqlErrNoRows {
	//	if errNoRowsNoMatter() {
	//		// 正常流程
	//	} else {
	//		// 异常流程
	//	}
	//}

	err2 := sqlQueryReturnInvalidQuery()
	isSqlErrNoRows(err2)
	
	err3 := sqlQueryReturnNoRowsWithWrap()
	isSqlErrNoRows(err3)
}
