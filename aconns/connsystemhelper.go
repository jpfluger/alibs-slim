package aconns

import (
	"github.com/Masterminds/semver/v3"
)

type ConnSystemHelper struct {
	ActionMap   ConnActionMap   `json:"actionMap"`
	Action      ConnActionType  `json:"action"`
	FromVersion *semver.Version `json:"fromVersion,omitempty"`
	ToVersion   *semver.Version `json:"toVersion,omitempty"`
}

func (sb *ConnSystemHelper) GetMap() ConnActionMap {
	return sb.ActionMap
}

func (sb *ConnSystemHelper) GetAction() ConnActionType {
	return sb.Action
}

func (sb *ConnSystemHelper) IsCreate() bool {
	return sb.Action == CONNACTIONTYPE_CREATE
}

func (sb *ConnSystemHelper) IsUpgrade() bool {
	return sb.Action == CONNACTIONTYPE_UPGRADE
}

func (sb *ConnSystemHelper) IsDowngrade() bool {
	return sb.Action == CONNACTIONTYPE_DOWNGRADE
}

func (sb *ConnSystemHelper) IsDelete() bool {
	return sb.Action == CONNACTIONTYPE_DELETE
}

func (sb *ConnSystemHelper) HasFromVersion() bool {
	return sb.FromVersion != nil
}

func (sb *ConnSystemHelper) HasToVersion() bool {
	return sb.ToVersion != nil
}

func (sb *ConnSystemHelper) GetFromVersion() semver.Version {
	return *sb.FromVersion
}

func (sb *ConnSystemHelper) GetToVersion() semver.Version {
	return *sb.ToVersion
}

func (sb *ConnSystemHelper) MustGetByAction() string {
	return sb.MustGet(sb.Action)
}

func (sb *ConnSystemHelper) MustGetCreate() string {
	return sb.MustGet(CONNACTIONTYPE_CREATE)
}

func (sb *ConnSystemHelper) MustGetUpgrade() string {
	return sb.MustGet(CONNACTIONTYPE_UPGRADE)
}

func (sb *ConnSystemHelper) MustGetDowngrade() string {
	return sb.MustGet(CONNACTIONTYPE_DOWNGRADE)
}

func (sb *ConnSystemHelper) MustGetDelete() string {
	return sb.MustGet(CONNACTIONTYPE_DELETE)
}

func (sb *ConnSystemHelper) MustGet(action ConnActionType) string {
	if action.IsEmpty() {
		return ""
	}
	item := sb.ActionMap.GetItem(action)
	if item == nil {
		return ""
	}
	return item.Text
}
