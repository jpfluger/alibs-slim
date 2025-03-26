package acron

import (
	"github.com/jpfluger/alibs-slim/areflect"
	"reflect"
)

const TYPEMANAGER_CRONTASKDATA = "crontaskdata"

func init() {
	_ = areflect.TypeManager().Register(TYPEMANAGER_CRONTASKDATA, "acron", returnTypeManagerCronJobData)
}

func returnTypeManagerCronJobData(typeName string) (reflect.Type, error) {
	var rtype reflect.Type // nil is the zero value for pointers, maps, slices, channels, and function types, interfaces, and other compound types.
	switch TaskType(typeName) {
	case TASKTYPE_SHELL:
		// Return the type of NoteFlag if typeName is "flag".
		rtype = reflect.TypeOf(TaskShell{})
	}
	// Return the determined reflect.Type and no error.
	return rtype, nil
}
