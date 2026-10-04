package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const queryCancelKey = "gateforge:query_cancel"

// registerQueryDeadline gives every GORM statement a context deadline when the caller did not set one.
func registerQueryDeadline(db *gorm.DB, timeout time.Duration) {
	if db == nil || timeout <= 0 {
		return
	}
	before := func(tx *gorm.DB) {
		ctx := tx.Statement.Context
		if ctx == nil {
			ctx = context.Background()
		}
		if _, ok := ctx.Deadline(); ok {
			return
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		tx.Statement.Context = ctx
		tx.InstanceSet(queryCancelKey, cancel)
	}
	after := func(tx *gorm.DB) {
		v, ok := tx.InstanceGet(queryCancelKey)
		if !ok {
			return
		}
		cancel, ok := v.(context.CancelFunc)
		if ok {
			cancel()
		}
	}
	_ = db.Callback().Query().Before("gorm:query").Register("gateforge:query_deadline", before)
	_ = db.Callback().Query().After("gorm:query").Register("gateforge:query_deadline_done", after)
	_ = db.Callback().Create().Before("gorm:create").Register("gateforge:query_deadline", before)
	_ = db.Callback().Create().After("gorm:create").Register("gateforge:query_deadline_done", after)
	_ = db.Callback().Update().Before("gorm:update").Register("gateforge:query_deadline", before)
	_ = db.Callback().Update().After("gorm:update").Register("gateforge:query_deadline_done", after)
	_ = db.Callback().Delete().Before("gorm:delete").Register("gateforge:query_deadline", before)
	_ = db.Callback().Delete().After("gorm:delete").Register("gateforge:query_deadline_done", after)
	_ = db.Callback().Row().Before("gorm:row").Register("gateforge:query_deadline", before)
	_ = db.Callback().Row().After("gorm:row").Register("gateforge:query_deadline_done", after)
	_ = db.Callback().Raw().Before("gorm:raw").Register("gateforge:query_deadline", before)
	_ = db.Callback().Raw().After("gorm:raw").Register("gateforge:query_deadline_done", after)
}
