package tarantool

import (
	"errors"

	"github.com/tinylib/msgp/msgp"
)

type Insert struct {
	Space interface{}
	Tuple []interface{}
}

var _ Query = (*Insert)(nil)

func (q *Insert) GetCommandID() uint {
	return InsertCommand
}

func (q *Insert) packMsg(data *packData, b []byte) (o []byte, err error) {
	if q.Tuple == nil {
		return o, errors.New("Tuple can not be nil")
	}

	o = b
	o = msgp.AppendMapHeader(o, 2)

	if o, err = data.packSpace(q.Space, o); err != nil {
		return o, err
	}

	o = msgp.AppendUint(o, KeyTuple)
	return msgp.AppendIntf(o, q.Tuple)
}

// MarshalMsg implements msgp.Marshaler
func (q *Insert) MarshalMsg(b []byte) (data []byte, err error) {
	return q.packMsg(defaultPackData, b)
}

// UnmarshalMsg implements msgp.Unmarshaler
func (q *Insert) UnmarshalMsg(data []byte) (buf []byte, err error) {
	var i uint32
	var k uint
	var t interface{}

	q.Space = nil
	q.Tuple = nil

	buf = data
	if i, buf, err = msgp.ReadMapHeaderBytes(buf); err != nil {
		return
	}

	for ; i > 0; i-- {
		if k, buf, err = msgp.ReadUintBytes(buf); err != nil {
			return
		}

		switch k {
		case KeySpaceNo:
			if q.Space, buf, err = msgp.ReadUintBytes(buf); err != nil {
				return
			}
		case KeyTuple:
			t, buf, err = msgp.ReadIntfBytes(buf)
			if q.Tuple = t.([]interface{}); q.Tuple == nil {
				return buf, errors.New("interface type is not []interface{}")
			}
		}
	}

	if q.Space == nil {
		return buf, errors.New("Insert.Unpack: no space specified")
	}
	if q.Tuple == nil {
		return buf, errors.New("Insert.Unpack: no tuple specified")
	}

	return
}
