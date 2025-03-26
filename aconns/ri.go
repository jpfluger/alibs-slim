package aconns

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/auser"
	"sync"
)

var ErrNoPK = fmt.Errorf("no primary key")

// RI is the Record Initializer structure.
type RI struct {
	// A ConnId or AdapterName is required.
	ConnId      ConnId
	AdapterName AdapterName

	// The user identifier to use when created a RecordSecurity for any given action.
	// This could be a username, email or id.
	User auser.RecordUserIdentity

	// If used, the event property describes the reason why the action took place.
	// For example, this value is set to "import" for bulk imports.
	// It could be made more granular (eg import from /path/to/file)
	Event string

	// This value can be nil.
	// DbModelManager importer will set this field otherwise is likely nil.
	sec *RecordSecurity

	// If true, wipes any history on save except for what sec contains.
	noApplyHistory bool

	// RI's typically do not survive beyond a single go routine
	mu sync.RWMutex
}

func NewRI(connId ConnId) *RI {
	return NewRIWithOptions(connId, "", auser.RecordUserIdentity{}, "")
}

func NewRIAdapter(adapterName AdapterName) *RI {
	return NewRIWithOptions(ConnId{}, adapterName, auser.RecordUserIdentity{}, "")
}

func NewRISEC(connId ConnId, user auser.RecordUserIdentity) *RI {
	return NewRIWithOptions(connId, "", user, "")
}

func NewRISECAdapter(adapterName AdapterName, user auser.RecordUserIdentity) *RI {
	return NewRIWithOptions(ConnId{}, adapterName, user, "")
}

func NewRIWithOptions(connId ConnId, adapterName AdapterName, user auser.RecordUserIdentity, event string) *RI {
	return &RI{
		ConnId:      connId,
		AdapterName: adapterName,
		User:        user,
		Event:       event,
	}
}

func (ri *RI) NewRI() *RI {
	return NewRIWithOptions(ri.ConnId, ri.AdapterName, auser.RecordUserIdentity{}, "")
}

func (ri *RI) NewRISEC() *RI {
	return NewRIWithOptions(ri.ConnId, ri.AdapterName, ri.User, "")
}

func (ri *RI) NewRISECUser(user auser.RecordUserIdentity) *RI {
	return NewRIWithOptions(ri.ConnId, ri.AdapterName, user, "")
}

func (ri *RI) GetConnId() ConnId {
	if ri == nil {
		return ConnId{}
	}
	ri.mu.RLock()
	defer ri.mu.RUnlock()
	return ri.ConnId
}

func (ri *RI) GetAdapterName() AdapterName {
	if ri == nil {
		return ""
	}
	ri.mu.RLock()
	defer ri.mu.RUnlock()
	return ri.AdapterName
}

func (ri *RI) SetRecordSecurity(sec *RecordSecurity) {
	ri.mu.Lock()
	defer ri.mu.Unlock()
	ri.sec = sec
}

func (ri *RI) SetNoApplyHistory(noApplyHistory bool) {
	ri.mu.Lock()
	defer ri.mu.Unlock()
	ri.noApplyHistory = noApplyHistory
}

func (ri *RI) GetRecordSecurity(action RecordActionType, target *RecordSecurity) (*RecordSecurity, error) {
	var rs *RecordSecurity

	if ri.sec != nil {
		// With this function the action is always applied.
		// Use the "event" to set additional information.

		if action.IsEmpty() {
			return nil, fmt.Errorf("ri.GetRecordSecurity action is empty")
		}

		ri.sec.Action = action

		// override -> such as DbModelManager Importer
		rs = ri.sec
	} else {

		if ri.User.IsEmpty() {
			return nil, fmt.Errorf("ri.GetRecordSecurity user is empty")
		}

		if action.IsEmpty() {
			return nil, fmt.Errorf("ri.GetRecordSecurity action is empty")
		}

		rs = NewRecordSecurity(ri.User, action, ri.Event)
	}

	if target != nil && target.IsValid() == nil && !ri.noApplyHistory {
		if err := rs.UpdateFrom(target); err != nil {
			return nil, fmt.Errorf("ri.GetRecordSecurity failed to update target; %v", err)
		}
	}

	return rs, nil
}

func (ri *RI) GetUser() auser.RecordUserIdentity {
	ri.mu.RLock()
	defer ri.mu.RUnlock()
	return ri.User
}

func (ri *RI) GetEvent() string {
	ri.mu.RLock()
	defer ri.mu.RUnlock()
	return ri.Event
}
