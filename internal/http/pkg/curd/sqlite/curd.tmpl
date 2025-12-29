package sqlite

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/noble-gase/ne/sqls"
)

// M 用于 INSERT & UPDATE
type M map[Column]any

func (m M) Split() (cols ColumnList, vals []any) {
	cap := len(m)

	cols = make(ColumnList, 0, cap)
	vals = make([]any, 0, cap)

	for k, v := range m {
		cols = append(cols, k)
		vals = append(vals, v)
	}
	return
}

// Create 创建记录
//
//	// 导入模块
//	import . "github.com/go-jet/jet/v2/sqlite"
//
//	// 语句示例
//	table.Demo.INSERT(table.Demo.Name).VALUES("hello")
//	// or
//	table.Demo.INSERT(table.Demo.Name).MODEL(model.Demo{Name: "hello"})
//
//	// 批量插入
//	table.Demo.INSERT(table.Demo.Name).
//		VALUES("hello").
//		VALUES("world")
//	// or
//	table.Demo.INSERT(table.Demo.Name).MODELS([]model.Demo{
//		{Name: "hello"},
//		{Name: "world"},
//	})
//
//	// 创建方法
//	sqlite.Create(ctx, db.DB(), stmt)
func Create(ctx context.Context, db qrm.DB, stmt InsertStatement) (int64, error) {
	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	ret, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}

	id, _ := ret.LastInsertId()
	return id, nil
}

// Update 更新记录
//
//	// 导入模块
//	import . "github.com/go-jet/jet/v2/sqlite"
//
//	// 语句示例
//	table.Demo.UPDATE(table.Demo.Name).SET("hello").WHERE(table.Demo.ID.EQ(Int64(1)))
//	// or
//	table.Demo.UPDATE(table.Demo.Name).MODEL(model.Demo{Name: "hello"}).WHERE(table.Demo.ID.EQ(Int64(1)))
//
//	// 更新方法
//	sqlite.Update(ctx, db.DB(), stmt)
func Update(ctx context.Context, db qrm.DB, stmt UpdateStatement) (int64, error) {
	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	ret, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}

	rows, _ := ret.RowsAffected()
	return rows, nil
}

// Delete 删除记录
//
//	// 导入模块
//	import . "github.com/go-jet/jet/v2/sqlite"
//
//	// 语句示例
//	table.Demo.DELETE().WHERE(table.Demo.ID.EQ(Int64(1)))
//
//	// 删除方法
//	sqlite.Delete(ctx, db.DB(), stmt)
func Delete(ctx context.Context, db qrm.DB, stmt DeleteStatement) (int64, error) {
	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	ret, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}

	rows, _ := ret.RowsAffected()
	return rows, nil
}

// FindOne 查询一条记录
//
//	// 导入模块
//	import . "github.com/go-jet/jet/v2/sqlite"
//
//	// 语句示例
//	table.Demo.SELECT(table.Demo.AllColumns).WHERE(table.Demo.ID.EQ(Int64(1)))
//	// or
//	SELECT(table.Demo.AllColumns).FROM(table.Demo).WHERE(table.Demo.ID.EQ(Int64(1)))
//
//	// 查询方法
//	sqlite.FindOne[model.Demo](ctx, db.DB(), stmt)
func FindOne[T any](ctx context.Context, db qrm.DB, stmt SelectStatement) (*T, error) {
	stmt = stmt.LIMIT(1)

	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	var dest T
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dest, nil
}

// FindAll 查询多条记录
//
//	// 导入模块
//	import . "github.com/go-jet/jet/v2/sqlite"
//
//	// 语句示例
//	table.Demo.SELECT(table.Demo.AllColumns).WHERE(table.Demo.Name.LIKE(String("%hello%")))
//	// or
//	SELECT(table.Demo.AllColumns).FROM(table.Demo).WHERE(table.Demo.Name.LIKE(String("%hello%")))
//
//	// 查询方法
//	sqlite.FindAll[model.Demo](ctx, db.DB(), stmt)
func FindAll[T any](ctx context.Context, db qrm.DB, stmt SelectStatement) ([]T, error) {
	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	var dest []T
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		return nil, err
	}
	return dest, nil
}

// Count 返回记录数
//
//	// 导入模块
//	import . "github.com/go-jet/jet/v2/sqlite"
//
//	// 查询方法
//	sqlite.Count(ctx, db.DB(), func(count SelectStatement) SelectStatement {
//		return count.FROM(table.Demo.Table).WHERE(table.Demo.Name.LIKE(String("%hello%")))
//	})
func Count(ctx context.Context, db qrm.DB, fn func(count SelectStatement) SelectStatement) (int64, error) {
	stmt := fn(SELECT(COUNT(STAR).AS("count")))

	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	var total struct {
		Count int64
	}
	if err := stmt.QueryContext(ctx, db, &total); err != nil {
		return 0, err
	}
	return total.Count, nil
}

// Paginate 分页查询
//
//	// 导入模块
//	import . "github.com/go-jet/jet/v2/sqlite"
//
//	// 查询方法
//	sqlite.Paginate[model.Demo](ctx, db.DB(), func(query SelectStatement) SelectStatement {
//		return query.FROM(table.Demo.Table).WHERE(table.Demo.Name.LIKE(String("%hello%")))
//	}, page, size, table.Demo.AllColumns, table.Demo.ID.DESC())
func Paginate[T any](ctx context.Context, db qrm.DB, fn func(query SelectStatement) SelectStatement, page, size int, cols ColumnList, orderBy ...OrderByClause) ([]T, int64, error) {
	stmt := fn(SELECT(COUNT(STAR).AS("count")))

	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	var total struct {
		Count int64
	}
	if err := stmt.QueryContext(ctx, db, &total); err != nil {
		return nil, 0, err
	}
	if total.Count == 0 {
		return []T{}, 0, nil
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size

	stmt = fn(SELECT(cols)).ORDER_BY(orderBy...).LIMIT(int64(size)).OFFSET(int64(offset))

	// SQL日志
	slog.InfoContext(ctx, sqls.Minify(stmt.DebugSql()))

	var dest []T
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		return nil, 0, err
	}
	return dest, total.Count, nil
}
