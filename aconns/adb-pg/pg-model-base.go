package adb_pg

import (
	"database/sql"
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/uptrace/bun"
	"time"
)

type FNBunSelect func(q *bun.SelectQuery) *bun.SelectQuery

func GetClientByRI(ri *aconns.RI) (*ADBPG, error) {
	base := PGModelBase{}
	if err := base.SetClientByRI(ri); err != nil {
		return nil, err
	}
	return base.GetClient(), nil
}

type PGModelBase struct {
	cli *ADBPG
	ri  *aconns.RI

	// The last modified date of this record for import purposes. It may or may not have a valid time.
	RecImportModifiedDate time.Time `bun:"-" json:"recImportModifiedDate" bun:"-"`
	// When RecImportModifiedDate has a valid time, then RecImportOnExistsType is used.
	RecImportOnExistsType aconns.RecImportOnExistsType `bun:"-" json:"recImportOnExistsType" bun:"-"`

	irw aconns.IImportRecordWrapper
}

func (base *PGModelBase) IsErrNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func (base *PGModelBase) GetClient() *ADBPG {
	return base.cli
}

func (base *PGModelBase) SetClient(cli *ADBPG) {
	base.cli = cli
}

func (base *PGModelBase) GetRI() *aconns.RI {
	return base.ri
}

func (base *PGModelBase) SetClientByRI(ri *aconns.RI) error {
	var pgc *ADBPG
	if connId := ri.GetConnId(); !connId.IsNil() {
		pgc = PGCONNS().Get(connId)
	} else {
		pgc = PGADAPTERS().Get(ri.GetAdapterName())
	}
	if pgc == nil {
		return fmt.Errorf("failed to located postgres client for ConnId %s", ri.GetConnId().String())
	}

	base.SetClient(pgc)
	base.ri = ri

	// Getting record security should happen at point of select, insert, delete
	// if !ri.IsRecordSecurityInitialized() {
	// 	ri.GetRecordSecurity()
	// }

	return nil
}

func (base *PGModelBase) SetClientByRIOnce(ri *aconns.RI) error {
	if base.GetClient() != nil {
		return nil
	}

	return base.SetClientByRI(ri)
}

func (base *PGModelBase) GetImportRecordWrapper() aconns.IImportRecordWrapper {
	return base.irw
}

func (base *PGModelBase) OnImportRunAction(objImport aconns.IRecordModel, ri *aconns.RI, irw aconns.IImportRecordWrapper) (aconns.ImportRecordResults, error) {
	base.irw = irw
	recordExists := false
	var action aconns.RecordActionType
	if objFound, err := objImport.SelectIntoNewObject(ri); err != nil {
		action = aconns.REC_ACTION_INSERT
	} else {
		recordExists = true
		if objFound.GetRecordSecurity().Time.IsZero() {
			action = aconns.REC_ACTION_INSERT
		} else {
			modDate := irw.GetModifiedDate()
			existsType := irw.GetRecImportOnExistsType()
			if base.RecImportOnExistsType > 0 {
				// if base.RecImportOnExistsType is 1 and riw.GetRecImportOnExistsType() is 2
				// then force an update b/c that what the more granular setting is.
				existsType = base.RecImportOnExistsType
			}
			if !base.RecImportModifiedDate.IsZero() {
				modDate = base.RecImportModifiedDate
				//existsType = base.RecImportOnExistsType
			}
			switch existsType {
			case aconns.REC_IMPORT_ON_EXISTS_UPDATE:
				action = aconns.REC_ACTION_UPDATE
			case aconns.REC_IMPORT_ON_EXISTS_UPDATE_IF_NEWER:
				if modDate.After(objFound.GetRecordSecurity().Time) {
					action = aconns.REC_ACTION_UPDATE
				}
			}
		}
	}

	var err error = nil

	switch action {
	case aconns.REC_ACTION_INSERT:
		err = objImport.Insert(ri)
	case aconns.REC_ACTION_UPDATE:
		err = objImport.Update(ri)
	default:
		return aconns.NewImportRecordResults(objImport.GetPrimaryKeyAsString(), fmt.Errorf("skipped import recordExists=%t", recordExists)), nil
	}

	return aconns.NewImportRecordResults(objImport.GetPrimaryKeyAsString(), err), err
}
