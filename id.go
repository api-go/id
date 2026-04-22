package id

import (
	"apigo.cc/go/rand"
	"sync"
	"time"

	"apigo.cc/go/encoding"
)

type IdMaker struct {
	secCurrent   uint64
	secIndexNext uint64
	secIndexLock sync.Mutex
	Incr         func(sec uint64) uint64
}

func NewIdMaker(incr func(sec uint64) uint64) *IdMaker {
	return &IdMaker{Incr: incr, secIndexNext: 1}
}

var DefaultIdMaker = NewIdMaker(nil)

func MakeId(size int) string {
	return DefaultIdMaker.Get(size)
}

func (im *IdMaker) defaultIncr(sec uint64) uint64 {
	im.secIndexLock.Lock()
	defer im.secIndexLock.Unlock()
	if im.secCurrent == sec {
		im.secIndexNext++
	} else {
		if im.secCurrent == 0 {
			im.secIndexNext = rand.FastInt(uint64(1000), uint64(1999))
		} else {
			im.secIndexNext = 1
		}
		im.secCurrent = sec
	}
	return im.secIndexNext
}

func (im *IdMaker) get(size int, ordered bool, hashToHead bool) string {
	tm := time.Now()

	nowSec := uint64(tm.Unix() - 946656000)
	var n, sec uint64
	secCapacity := uint64(901356495)
	if nowSec < 11*secCapacity {
		n = nowSec / secCapacity
		sec = (nowSec % secCapacity) + 14776336
	} else {
		n = 11
		sec = nowSec
	}

	var secIndex uint64
	if im.Incr != nil {
		secIndex = im.Incr(sec)
	}
	if secIndex == 0 {
		secIndex = im.defaultIncr(sec)
	}

	intEncoder := encoding.DefaultIntEncoder
	if ordered {
		intEncoder = encoding.OrderedIntEncoder
	}
	secBytes := intEncoder.EncodeInt(sec)
	secLen := len(secBytes)
	inSecIndexBytes := intEncoder.EncodeInt(secIndex)

	m := min(uint64(len(inSecIndexBytes)), 5)

	secTagVal := n*5 + (m - 1)
	var uid = make([]byte, 0, size)
	uid = intEncoder.AppendInt(uid, secTagVal)

	uid = append(uid, secBytes...)
	uid = append(uid, inSecIndexBytes...)
	uid = intEncoder.FillInt(uid, size)
	if !ordered {
		encoding.HashInt(encoding.ExchangeInt(uid), nil)
	} else {
		encoding.HashInt(encoding.ExchangeInt(uid[secLen+1:]), nil)
		if hashToHead {
			size = len(uid)
			lastByte := uid[size-1]
			copy(uid[1:], uid[:size-1])
			uid[0] = lastByte
		}
	}
	return string(uid)
}

func (im *IdMaker) Get(size int) string {
	return im.get(size, false, false)
}

func (im *IdMaker) GetForMysql(size int) string {
	return im.get(size, true, true)
}

func (im *IdMaker) GetForPostgreSQL(size int) string {
	return im.get(size, true, false)
}
