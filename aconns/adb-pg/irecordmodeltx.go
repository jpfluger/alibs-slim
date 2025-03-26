package adb_pg

import (
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/uptrace/bun"
)

type IRecordModelTx interface {
	InsertTx(tx *bun.Tx, ri *aconns.RI) error
	UpdateTx(tx *bun.Tx, ri *aconns.RI) error
	DeleteTx(tx *bun.Tx, ri *aconns.RI) error
}
