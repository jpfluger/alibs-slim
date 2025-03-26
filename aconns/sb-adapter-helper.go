package aconns

import "github.com/Masterminds/semver/v3"

type ISBAdapterHelper interface {
	GetMap() ConnSystemHelper
	CheckVersionState() ConnActionType
	IsCreate() bool
	IsUpgrade() bool
	IsDowngrade() bool
	IsDelete() bool
	HasFromVersion() bool
	HasToVersion() bool
	GetFromVersion() semver.Version
	GetToVersion() semver.Version
	MustGetByAction() string
	MustGetCreate() string
	MustGetUpgrade() string
	MustGetDowngrade() string
	MustGetDelete() string
	MustGet(action ConnActionType) string
}
