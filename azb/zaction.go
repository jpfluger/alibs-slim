package azb

import (
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/autils"
	"strconv"
	"time"
)

// ZAction is used to simplify web-ui/api calls via JavaScript.
// It is embedded in structs that bind incoming data for API calls.
// It works in conjunction with https://github.com/jpfluger/zazzy-browser
//
// When creating a struct with which to bind incoming data, first embed the ZAction struct inside it:
//
//	type DIN struct {
//		ZAction ZAction `json:"zaction"`
//	}
//
// Complementary structs are built-into ZAction but these are invoked on a route-by-route basis. They include:
//
//	type MyDIN struct {
//	  DIN
//	  Data Contact `json:"data"`
//	}
//
// ExtractPayload is for custom Payload of type "interface{}"
// This is in conjunction with option (2) above but option (1) is better.
//
//	 func (din *MyDIN) ExtractPayload(b []byte) error {
//	   // unmarshal b to Data
//	}
type ZAction struct {
	Event     ZBType   `json:"event"`     // Event type, e.g., "zurl-dialog" for dialog events.
	Mod       ZBType   `json:"mod"`       // Module type, e.g., "paginate" for pagination.
	Id        string   `json:"zid"`       // Unique identifier for the action.
	IdParent  string   `json:"zidParent"` // Parent identifier for nested actions.
	Sequence  string   `json:"zseq"`      // Sequence number for ordering actions.
	LoopType  ZBType   `json:"loopType"`  // Type of loop, if applicable.
	PageOn    int      `json:"pageOn"`    // Current page number for pagination.
	PageLimit int      `json:"pageLimit"` // Number of items per page for pagination.
	ViewPort  ViewPort `json:"viewPort"`
}

const (
	ZMOD_PAGINATE = ZBType("paginate") // Pagination module type.
)

// IsDialog checks if the ZAction event is a dialog.
func (za *ZAction) IsDialog() bool {
	return za != nil && za.Event == "zurl-dialog"
}

// HasZId checks if the ZAction has a non-empty Id.
func (za *ZAction) HasZId() bool {
	return za.Id != ""
}

// IsPaginate checks if the ZAction module is for pagination.
func (za *ZAction) IsPaginate() bool {
	return za.Mod == ZMOD_PAGINATE
}

// HasZIdParent checks if the ZAction has a non-empty IdParent.
func (za *ZAction) HasZIdParent() bool {
	return za.IdParent != ""
}

// ToUUIDZId converts the ZAction Id to a UUID.
func (za *ZAction) ToUUIDZId() uuid.UUID {
	if za.Id == "" {
		return uuid.Nil
	}
	return autils.ParseUUID(za.Id)
}

// ToIntZId converts the ZAction Id to an integer.
func (za *ZAction) ToIntZId() int {
	if za.Id == "" {
		return 0
	}
	val, err := strconv.Atoi(za.Id)
	if err != nil {
		return 0
	}
	return val
}

// ToFloatZId converts the ZAction Id to a float64.
func (za *ZAction) ToFloatZId() float64 {
	if za.Id == "" {
		return 0
	}
	val, err := strconv.ParseFloat(za.Id, 64)
	if err != nil {
		return 0
	}
	return val
}

// ToTimeZId converts the ZAction Id to a time.Time.
func (za *ZAction) ToTimeZId() time.Time {
	if za.Id == "" {
		return time.Time{}
	}
	tm, err := time.Parse(time.RFC3339Nano, za.Id)
	if err != nil {
		return time.Time{}
	}
	return tm
}

// ToUUIDZIdParent converts the ZAction IdParent to a UUID.
func (za *ZAction) ToUUIDZIdParent() uuid.UUID {
	if za.IdParent == "" {
		return uuid.Nil
	}
	return autils.ParseUUID(za.IdParent)
}

// ToIntParent converts the ZAction IdParent to an integer.
func (za *ZAction) ToIntParent() int {
	if za.IdParent == "" {
		return 0
	}
	val, err := strconv.Atoi(za.IdParent)
	if err != nil {
		return 0
	}
	return val
}

// ToFloatZIdParent converts the ZAction IdParent to a float64.
func (za *ZAction) ToFloatZIdParent() float64 {
	if za.IdParent == "" {
		return 0
	}
	val, err := strconv.ParseFloat(za.IdParent, 64)
	if err != nil {
		return 0
	}
	return val
}

// ToTimeZIdParent converts the ZAction IdParent to a time.Time.
func (za *ZAction) ToTimeZIdParent() time.Time {
	if za.IdParent == "" {
		return time.Time{}
	}
	tm, err := time.Parse(time.RFC3339Nano, za.IdParent)
	if err != nil {
		return time.Time{}
	}
	return tm
}

// ToUUIDZSequence converts the ZAction Sequence to a UUID.
func (za *ZAction) ToUUIDZSequence() uuid.UUID {
	if za.Sequence == "" {
		return uuid.Nil
	}
	return autils.ParseUUID(za.Sequence)
}

// ToIntZSequence converts the ZAction Sequence to an integer.
func (za *ZAction) ToIntZSequence() int {
	if za.Sequence == "" {
		return 0
	}
	val, err := strconv.Atoi(za.Sequence)
	if err != nil {
		return 0
	}
	return val
}

// ToFloatZSequence converts the ZAction Sequence to a float64.
func (za *ZAction) ToFloatZSequence() float64 {
	if za.Sequence == "" {
		return 0
	}
	val, err := strconv.ParseFloat(za.Sequence, 64)
	if err != nil {
		return 0
	}
	return val
}

// ToTimeZSequence converts the ZAction Sequence to a time.Time.
func (za *ZAction) ToTimeZSequence() time.Time {
	if za.Sequence == "" {
		return time.Time{}
	}
	tm, err := time.Parse(time.RFC3339Nano, za.Sequence)
	if err != nil {
		return time.Time{}
	}
	return tm
}
