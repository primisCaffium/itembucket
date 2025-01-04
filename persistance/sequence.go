package persistance

import (
	"itembucket/common"
	"sync"
)

type Sequence struct {
	Id  *int64
	Mux *sync.Mutex
}

func NewSequence(initialValue *int64) *Sequence {
	iv := initialValue
	if iv == nil {
		iv = common.PInt64(0)
	}

	return &Sequence{
		Id:  iv,
		Mux: &sync.Mutex{},
	}
}

func (o *Sequence) Next() *int64 {
	defer func() {
		o.Mux.Unlock()
	}()

	o.Mux.Lock()
	o.Id = common.PInt64(*o.Id + 1)
	return o.Id
}
