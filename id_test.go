package id_test

import (
	"testing"
	"apigo.cc/go/id"
)

func TestMakeID(t *testing.T) {
	uid := id.MakeID(10)
	if len(uid) != 10 {
		t.Errorf("expected length 10, got %d", len(uid))
	}
}

func TestGetForMysql(t *testing.T) {
	uid := id.DefaultIDMaker.GetForMysql(10)
	if len(uid) != 10 {
		t.Errorf("expected length 10, got %d", len(uid))
	}
}

func TestGetForPostgreSQL(t *testing.T) {
	uid := id.DefaultIDMaker.GetForPostgreSQL(10)
	if len(uid) != 10 {
		t.Errorf("expected length 10, got %d", len(uid))
	}
}

func BenchmarkIDMaker(b *testing.B) {
	b.Run("MakeID-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.MakeID(10)
		}
	})
	b.Run("GetForMysql-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = id.DefaultIDMaker.GetForMysql(10)
		}
	})
}
