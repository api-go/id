package id

import (
	"apigo.cc/go/rand"
	"sync"
	"time"

	"apigo.cc/go/encoding"
)

const Epoch = 946656000

type IDMaker struct {
	secCurrent   uint64
	secIndexNext uint64
	secIndexLock sync.Mutex
	Incr         func(sec uint64) uint64
}

func NewIDMaker(incr func(sec uint64) uint64) *IDMaker {
	return &IDMaker{Incr: incr, secIndexNext: 1}
}

var DefaultIDMaker = NewIDMaker(nil)

func MakeID(size int) string {
	return DefaultIDMaker.Get(size)
}

func (im *IDMaker) defaultIncr(sec uint64) uint64 {
	im.secIndexLock.Lock()
	defer im.secIndexLock.Unlock()

	if im.secCurrent == sec {
		im.secIndexNext++
		return im.secIndexNext
	}

	if im.secCurrent == 0 {
		im.secIndexNext = uint64(rand.FastInt(1000, 1999))
	} else {
		im.secIndexNext = 1
	}

	im.secCurrent = sec
	return im.secIndexNext
}

func (im *IDMaker) get(size int, ordered bool, hashToHead bool) string {
	tm := time.Now()

	nowSec := uint64(tm.Unix() - Epoch)
	var n, sec uint64
	const secCapacity = 901356495
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
		uid = encoding.ExchangeInt(uid)
		hashed := encoding.HashInt(uid, nil)
		uid = make([]byte, 0, size)
		for i := 0; i < size && i < len(hashed); i++ {
			uid = append(uid, intEncoder.EncodeInt(uint64(hashed[i]))[0])
		}
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

func (im *IDMaker) Get(size int) string {
	return im.get(size, false, false)
}

func (im *IDMaker) Get8Bytes4KPerSecond() string {
	return im.get(8, false, false)
}

func (im *IDMaker) Get9Bytes90KPerSecond() string {
	return im.get(9, false, false)
}

func (im *IDMaker) Get10Bytes14MPerSecond() string {
	return im.get(10, false, false)
}

func (im *IDMaker) Get11Bytes900MPerSecond() string {
	return im.get(11, false, false)
}

func (im *IDMaker) Get12BytesUltraPerSecond() string {
	return im.get(12, false, false)
}

func (im *IDMaker) GetForMysql(size int) string {
	return im.get(size, true, true)
}

func (im *IDMaker) GetForPostgreSQL(size int) string {
	return im.get(size, true, false)
}

func Get8Bytes4KPerSecond() string {
	return DefaultIDMaker.Get8Bytes4KPerSecond()
}

func Get9Bytes90KPerSecond() string {
	return DefaultIDMaker.Get9Bytes90KPerSecond()
}

func Get10Bytes14MPerSecond() string {
	return DefaultIDMaker.Get10Bytes14MPerSecond()
}

func Get11Bytes900MPerSecond() string {
	return DefaultIDMaker.Get11Bytes900MPerSecond()
}

func Get12BytesUltraPerSecond() string {
	return DefaultIDMaker.Get12BytesUltraPerSecond()
}
