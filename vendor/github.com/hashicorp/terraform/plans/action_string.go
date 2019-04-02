// Code generated by "stringer -type Action"; DO NOT EDIT.

package plans

import "strconv"

const (
	_Action_name_0 = "NoOp"
	_Action_name_1 = "Create"
	_Action_name_2 = "Delete"
	_Action_name_3 = "Update"
	_Action_name_4 = "CreateThenDelete"
	_Action_name_5 = "Read"
	_Action_name_6 = "DeleteThenCreate"
)

func (i Action) String() string {
	switch {
	case i == 0:
		return _Action_name_0
	case i == 43:
		return _Action_name_1
	case i == 45:
		return _Action_name_2
	case i == 126:
		return _Action_name_3
	case i == 177:
		return _Action_name_4
	case i == 8592:
		return _Action_name_5
	case i == 8723:
		return _Action_name_6
	default:
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
