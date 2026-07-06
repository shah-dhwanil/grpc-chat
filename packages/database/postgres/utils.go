package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shah-dhwanil/grpc-chat/packages/pkgerror"
)
func ExtractOrderParam(orderBy string) (string,string) {
	if strings.HasPrefix(orderBy, "-") {
		return orderBy[1:], "DESC"
	}
	if strings.HasPrefix(orderBy, "+") {
		return orderBy[1:], "ASC"
	}
	return orderBy, "ASC"
}

func ConstructWhereClause(conditions []string) string {
	if len(conditions) == 0 {
		return "1=1" // No conditions, so we return a tautology
	}
	return strings.Join(conditions, " AND ")
}

func ConstructOrderByClause(orderBy []string) string {
	if len(orderBy) == 0 {
		return "created_at DESC" // Default order by
	}
	return strings.Join(orderBy, ", ")
}

func ConstructSetClause(fields []string) string {
	return strings.Join(fields, ", ")
}

func QueryInTransaction[T any](ctx context.Context, executor DBTX, fn func(Tx)(T,error)) (T,error) {
	var zero T
	if executor == nil {
		return zero, pkgerror.NewUnknownError(nil, "DATABSE_ERROR", "Database connection is not initialized", nil)
	}
	txn,err:=executor.Begin(ctx)
	if err!=nil {
		dbError,ok := ConvertPgError(err)
		if ok {
			return zero,dbError
		}
		return zero,pkgerror.NewUnknownError(err,"DATABSE_ERROR","Unknown Error while starting transaction",nil)
	}
	rows,err:=fn(txn)
	if err!=nil {
		if rbErr:=txn.Rollback(ctx);rbErr!=nil {
			return zero,fmt.Errorf("transaction failed: %v, rollback failed: %v", err, rbErr)
		}
		dbError,ok :=ConvertPgError(err)
		if ok {
			return zero,dbError
		}
		return zero,pkgerror.NewUnknownError(err,"DATABSE_ERROR","Unknown Error while executing transaction",nil)
	}
	if err = txn.Commit(ctx); err != nil {
		dbError,ok :=ConvertPgError(err)
		if ok {
			return zero,dbError
		}
		return zero,pkgerror.NewUnknownError(err,"DATABSE_ERROR","Unknown Error while committing transaction",nil)
	}
	return rows,nil
}

var ExecuteInTransaction func(ctx context.Context, executor DBTX, fn func(Tx) (pgconn.CommandTag,error)) (pgconn.CommandTag,error) = QueryInTransaction