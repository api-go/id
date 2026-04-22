package id_test

import (
	"testing"
	"apigo.cc/go/id"
)

func TestMakeId(t *testing.T) {
	uid := id.MakeId(10)
	if len(uid) != 10 {
		t.Errorf("expected length 10, got %d", len(uid))
	}
}

func TestGetForMysql(t *testing.T) {
	uid := id.DefaultIdMaker.GetForMysql(10)
	if len(uid) != 10 {
		t.Errorf("expected length 10, got %d", len(uid))
	}
}

func TestGetForPostgreSQL(t *testing.T) {
	uid := id.DefaultIdMaker.GetForPostgreSQL(10)
	if len(uid) != 10 {
		t.Errorf("expected length 10, got %d", len(uid))
	}
}

func BenchmarkIdMaker(b *testing.B) {
	b.Run("MakeId-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.MakeId(10)
		}
	})
	b.Run("GetForMysql-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.DefaultIdMaker.GetForMysql(10)
		}
	})
}
