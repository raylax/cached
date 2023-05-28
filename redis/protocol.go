package redis

import (
	"fmt"
	"strconv"
)

type RespType byte

const (
	RespTypeInteger      RespType = ':'
	RespTypeSimpleString RespType = '+'
	RespTypeBulkString   RespType = '$'
	RespTypeArray        RespType = '*'
	RespTypeError        RespType = '-'
)

type resp struct {
	Type RespType
}

func (r *resp) RespType() RespType {
	return r.Type
}

type RespInteger struct {
	resp
	Value int
}

type RespString struct {
	resp
	Data []byte
}

type RespArray struct {
	resp
	Elements []any
}

type RespError struct {
	resp
	Message string
}

func ReadResp(r Reader) (any, error) {
	return readResp(&respReader{r: r})
}

func readResp(r *respReader) (any, error) {

	t, err := r.readType()
	if err != nil {
		return nil, err
	}

	switch t {
	case RespTypeInteger:
		i, err := r.readInt()
		if err != nil {
			return nil, err
		}
		return &RespInteger{resp: resp{Type: t}, Value: i}, nil
	case RespTypeSimpleString:
		b, err := r.readUntilCRLF()
		if err != nil {
			return nil, err
		}
		return &RespString{resp: resp{Type: t}, Data: b}, nil
	case RespTypeBulkString:
		n, err := r.readInt()
		if err != nil {
			return nil, err
		}

		if n == -1 {
			return &RespString{resp: resp{Type: t}, Data: nil}, nil
		}

		b, err := r.read(n)
		if err != nil {
			return nil, err
		}
		err = r.skipCRLF()
		if err != nil {
			return nil, err
		}
		return &RespString{resp: resp{Type: t}, Data: b}, nil
	case RespTypeArray:
		count, err := r.readInt()
		if err != nil {
			return nil, err
		}

		if count == -1 {
			return &RespArray{resp: resp{Type: t}, Elements: nil}, nil
		}

		elements := make([]any, count)
		for i := 0; i < count; i++ {
			elements[i], err = readResp(r)
			if err != nil {
				return nil, err
			}
		}
		return &RespArray{resp: resp{Type: t}, Elements: elements}, nil

	case RespTypeError:
		b, err := r.readUntilCRLF()
		if err != nil {
			return nil, err
		}
		return &RespError{resp: resp{Type: t}, Message: string(b)}, nil
	default:
		return nil, fmt.Errorf("unknown type %d", t)
	}

}

type respReader struct {
	r Reader
}

func (r *respReader) read(n int) ([]byte, error) {
	return r.r.Next(n)
}

func (r *respReader) skipCRLF() (err error) {
	return r.r.Skip(2)
}

func (r *respReader) readUntilCRLF() (b []byte, err error) {
	n := 0
	for {
		b, err = r.r.Peek(n + 2)
		if err != nil {
			return
		}
		if b[n] == '\r' && b[n+1] == '\n' {
			break
		}
		if b[n+1] == '\r' {
			n += 1
		} else {
			n += 2
		}
	}
	b, err = r.r.Next(n)
	if err != nil {
		return
	}
	err = r.skipCRLF()
	return
}

func (r *respReader) readType() (t RespType, err error) {
	b, err := r.read(1)
	if err != nil {
		return
	}
	t = RespType(b[0])
	return
}

func (r *respReader) readInt() (i int, err error) {
	b, err := r.readUntilCRLF()
	if err != nil {
		return
	}
	i, err = strconv.Atoi(string(b))
	return
}

type Reader interface {
	Next(n int) (p []byte, err error)
	Peek(n int) (buf []byte, err error)
	Skip(n int) (err error)
}
