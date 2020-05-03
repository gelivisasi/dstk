package core

import dstk "github.com/anujga/dstk/pkg/api/proto"

const (
	//note: for backward compatibility, we can only add entries in the end
	ErrUnknown               = -1
	ErrKeyNotFound           = 1
	ErrInvalidPartition      = 2
	ErrChannelFull           = iota
	ErrCantCreateSnapshotDir = iota
	ErrPeerAlreadyExists     = iota
	ErrNodeAlreadyExists     = iota

	//...
	ErrMaxErrorCode = iota
)

type Errr struct {
	dstk.Ex
}

func NewErr(id int64, msg string) *Errr {
	return &Errr{dstk.Ex{Id: id, Msg: msg}}
}

func FromErr(err error) *Errr {
	return NewErr(ErrUnknown, err.Error())
}

func FromRErr(err dstk.Ex) *Errr {
	return &Errr{err}
}

func (m *Errr) Error() string {
	return m.Msg
}
