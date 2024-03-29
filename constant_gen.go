// THIS FILE IS A GENERATED CODE. DO NOT EDIT
// generated version: 0.4.0
package todolist

import "golang.org/x/xerrors"

// OpType - operator type
type OpType = string

const (
	OpTypeEqual              OpType = "=="
	OpTypeLessThan           OpType = "<"
	OpTypeLessThanOrEqual    OpType = "<="
	OpTypeGreaterThan        OpType = ">"
	OpTypeGreaterThanOrEqual OpType = ">="
	OpTypeIn                 OpType = "in"
	OpTypeArrayContains      OpType = "array-contains"
	OpTypeArrayContainsAny   OpType = "array-contains-any"
)

var (
	ErrAlreadyExists  = xerrors.New("already exists")
	ErrAlreadyDeleted = xerrors.New("already been deleted")
	ErrNotFound       = xerrors.New("not found")
)
