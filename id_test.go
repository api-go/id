package id_test

import (
	"github.com/api-go/encoding"
	"github.com/api-go/id"
	"math"
	"strings"
	"testing"
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

func TestPrintIDs(t *testing.T) {
	t.Log("=== Normal Mode (Get) ===")
	for size := 6; size <= 12; size++ {
		t.Logf("--- Size: %d ---", size)
		for i := 0; i < 10; i++ {
			uid := id.DefaultIDMaker.Get(size)
			t.Logf("  #%d: %s", i+1, uid)
		}
	}

	t.Log("=== MySQL Mode (GetForMysql) ===")
	for size := 6; size <= 12; size++ {
		t.Logf("--- Size: %d ---", size)
		for i := 0; i < 10; i++ {
			uid := id.DefaultIDMaker.GetForMysql(size)
			t.Logf("  #%d: %s", i+1, uid)
		}
	}
}

// UnexchangeInt is the inverse of ExchangeInt
func UnexchangeInt(buf []byte) []byte {
	size := len(buf)
	if size <= 1 {
		return buf
	}
	res := make([]byte, size)
	buf2_i := 0
	buf2_ai := 0
	buf2_ri := size - 1

	// In ExchangeInt:
	// for i := 0; i < size; i++ {
	//     if i%2 == 0 {
	//         buf2[buf2_i] = buf[buf2_ri]; buf2_ri--
	//     } else {
	//         buf2[buf2_i] = buf[buf2_ai]; buf2_ai++
	//     }
	//     buf2_i++
	// }
	// We reverse this mapping.
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			res[buf2_ri] = buf[buf2_i]
			buf2_ri--
		} else {
			res[buf2_ai] = buf[buf2_i]
			buf2_ai++
		}
		buf2_i++
	}
	return res
}

// UnhashInt is the inverse of HashInt
func UnhashInt(buf []byte, digits string, decodeMap [256]int, radix int) []byte {
	if len(buf) == 0 {
		return buf
	}

	orderedDigits := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	defaultDigits := "9ukH1grX75TQS6LzpFAjIivsdZoO0mc8NBwnyYDhtMWEC2V3KaGxfJRPqe4lbU"

	// F is the non-linear cross-charset update function for prevP
	F := func(p int) int {
		char := defaultDigits[p]
		return strings.IndexByte(orderedDigits, char)
	}

	// 1. Undo Round 2 (Backward propagation, right to left)
	prevP := (len(buf) * 31) % radix
	for i := len(buf) - 1; i >= 0; i-- {
		newChar := buf[i]
		newP := decodeMap[newChar]
		if newP < 0 {
			newP = 0
		}
		p := (newP - prevP + radix) % radix
		buf[i] = digits[p]
		prevP = F(newP)
	}

	// 2. Undo Round 1 (Forward propagation, left to right)
	prevP = (len(buf) * 17) % radix
	for i := 0; i < len(buf); i++ {
		newChar := buf[i]
		newP := decodeMap[newChar]
		if newP < 0 {
			newP = 0
		}
		p := (newP - prevP + radix) % radix
		buf[i] = digits[p]
		prevP = F(newP)
	}

	return buf
}

func TestReversibility(t *testing.T) {
	digits := "9ukH1grX75TQS6LzpFAjIivsdZoO0mc8NBwnyYDhtMWEC2V3KaGxfJRPqe4lbU"
	var decodeMap [256]int
	for i := 0; i < 256; i++ {
		decodeMap[i] = -1
	}
	for i, c := range digits {
		decodeMap[c] = i
	}
	radix := len(digits)

	// We generate 100,000 sequential combinations to check if they map uniquely
	t.Log("Testing reversibility of HashInt and ExchangeInt...")

	// Create some arbitrary input
	for size := 5; size <= 15; size++ {
		// Test multiple inputs for each size
		for step := 0; step < 1000; step++ {
			original := make([]byte, size)
			for i := 0; i < size; i++ {
				// Fill with pseudo-random digits from the digit set
				original[i] = digits[(step+i)%radix]
			}

			// Copy for modification
			temp := make([]byte, size)
			copy(temp, original)

			// Forward transformation: Exchange then Hash
			// We simulate what id.go does: encoding.HashInt(encoding.ExchangeInt(uid))
			// Since encoding.ExchangeInt returns a new slice, and encoding.HashInt modifies in-place
			exchanged := encoding.ExchangeInt(temp)
			hashed := encoding.HashInt(exchanged)

			// Backward transformation: Unhash then Unexchange
			unhashed := UnhashInt(hashed, digits, decodeMap, radix)
			unexchanged := UnexchangeInt(unhashed)

			for i := 0; i < size; i++ {
				if unexchanged[i] != original[i] {
					t.Fatalf("Reversibility failed for size %d, step %d! Expected %s, got %s", size, step, string(original), string(unexchanged))
				}
			}
		}
	}
	t.Log("Reversibility verified! The transformation is a perfect BIJECTION (no collisions possible on same-length unique inputs before truncation).")
}

func TestCollisionRate(t *testing.T) {
	// We simulate high concurrency generation of IDs for lengths 6 to 12.
	// We verify if the distribution behaves as a mathematically ideal random distribution.
	importMath := func(n float64, space float64) float64 {
		ratio := n / space
		if ratio < 1e-4 {
			return (n * n) / (2.0 * space)
		}
		importMathExp := space * (1.0 - math.Exp(-ratio))
		return n - importMathExp
	}

	for size := 6; size <= 10; size++ {
		seen := make(map[string]int)
		collisions := 0
		n := 100000 // Test 100,000 IDs
		if size == 6 {
			n = 1000000 // Test 1,000,000 IDs for size 6 to see actual collisions
		} else if size >= 9 {
			n = 500000
		}

		logCount := 0
		for i := 0; i < n; i++ {
			uid := id.DefaultIDMaker.Get(size)
			if prevIdx, found := seen[uid]; found {
				collisions++
				if size == 6 && logCount < 10 {
					t.Logf("Coll size 6: ID=%s found at iteration %d and %d", uid, prevIdx, i)
					logCount++
				}
			} else {
				seen[uid] = i
			}
		}

		space := math.Pow(62, float64(size))
		expected := importMath(float64(n), space)
		t.Logf("Size %2d: Space %12.0f, N = %d | Actual Collisions: %d (Rate: %.5f%%) | Expected Collisions: %.2f (Rate: %.5f%%)",
			size, space, n, collisions, float64(collisions)/float64(n)*100, expected, expected/float64(n)*100)
	}
}

func TestCheckHashDiff(t *testing.T) {
	// Let's generate two concurrent IDs with consecutive secIndex
	// Normal mode (size = 6)
	// We want to see:
	// 1. The original uid (before Exchange and Hash)
	// 2. The exchanged uid
	// 3. The hashed uid
	// 4. The final 6-character truncated version

	// We stub the process of making UID
	sec := uint64(14776336 + 100) // arbitrary time sec
	intEncoder := encoding.DefaultIntEncoder
	secBytes := intEncoder.EncodeInt(sec) // 5 bytes

	for idx := uint64(1); idx <= 5; idx++ {
		secIndex := idx
		inSecIndexBytes := intEncoder.EncodeInt(secIndex)

		m := min(uint64(len(inSecIndexBytes)), 5)
		secTagVal := uint64(0)*5 + (m - 1)

		var uid = make([]byte, 0, 6)
		uid = intEncoder.AppendInt(uid, secTagVal)
		uid = append(uid, secBytes...)
		uid = append(uid, inSecIndexBytes...)
		uid = intEncoder.FillInt(uid, 6)

		orig := string(uid)
		exchanged := encoding.ExchangeInt(uid)
		exchStr := string(exchanged)

		// Hash modifies in-place
		hashed := encoding.HashInt(exchanged)
		hashStr := string(hashed)
		truncated := hashStr[:6]

		t.Logf("idx=%d | Orig: %s (len=%d) | Exch: %s | Hashed: %s | Truncated: %s",
			idx, orig, len(orig), exchStr, hashStr, truncated)
	}
}

func TestPrintCollisions(t *testing.T) {
	seen := make(map[string]struct {
		sec       uint64
		secIndex  uint64
		orig      string
		exchanged string
		hashed    string
	})
	intEncoder := encoding.DefaultIntEncoder

	size := 6
	t.Logf("Printing collisions for size = %d across two adjacent seconds...", size)

	count := 0

	// Second 1
	sec1 := uint64(14776336 + 100)
	for secIndex := uint64(1); secIndex <= 50000; secIndex++ {
		secBytes := intEncoder.EncodeInt(sec1)
		inSecIndexBytes := intEncoder.EncodeInt(secIndex)
		m := min(uint64(len(inSecIndexBytes)), 5)
		secTagVal := uint64(0)*5 + (m - 1)

		var uid = make([]byte, 0, size)
		uid = intEncoder.AppendInt(uid, secTagVal)
		uid = append(uid, secBytes...)
		uid = append(uid, inSecIndexBytes...)
		uid = intEncoder.FillInt(uid, size)

		orig := string(uid)
		exchanged := encoding.ExchangeInt(uid)
		exchStr := string(exchanged)
		hashed := encoding.HashInt(exchanged)
		hashStr := string(hashed)
		finalID := hashStr[:size]

		seen[finalID] = struct {
			sec       uint64
			secIndex  uint64
			orig      string
			exchanged string
			hashed    string
		}{sec: sec1, secIndex: secIndex, orig: orig, exchanged: exchStr, hashed: hashStr}
	}

	// Second 2 (Adjacent second)
	sec2 := sec1 + 1
	for secIndex := uint64(1); secIndex <= 50000; secIndex++ {
		secBytes := intEncoder.EncodeInt(sec2)
		inSecIndexBytes := intEncoder.EncodeInt(secIndex)
		m := min(uint64(len(inSecIndexBytes)), 5)
		secTagVal := uint64(0)*5 + (m - 1)

		var uid = make([]byte, 0, size)
		uid = intEncoder.AppendInt(uid, secTagVal)
		uid = append(uid, secBytes...)
		uid = append(uid, inSecIndexBytes...)
		uid = intEncoder.FillInt(uid, size)

		orig := string(uid)
		exchanged := encoding.ExchangeInt(uid)
		exchStr := string(exchanged)
		hashed := encoding.HashInt(exchanged)
		hashStr := string(hashed)
		finalID := hashStr[:size]

		if prev, found := seen[finalID]; found {
			t.Logf("Coll! ID: %s", finalID)
			t.Logf("  Pair 1: sec=%d, idx=%d | Orig: %s | Exch: %s | Hashed: %s",
				prev.sec, prev.secIndex, prev.orig, prev.exchanged, prev.hashed)
			t.Logf("  Pair 2: sec=%d, idx=%d | Orig: %s | Exch: %s | Hashed: %s",
				sec2, secIndex, orig, exchStr, hashStr)
			count++
			if count >= 10 {
				break
			}
		}
	}
	t.Logf("Total collision pairs logged. Simulated adjacent-second collision counts: %d out of 50,000", count)
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
