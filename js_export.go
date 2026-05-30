package id

import "apigo.cc/go/jsmod"

func init() {
	jsmod.Register("id", map[string]any{
		"get8Bytes4KPerSecond":    Get8Bytes4KPerSecond,
		"get9Bytes90KPerSecond":    Get9Bytes90KPerSecond,
		"get10Bytes14MPerSecond":   Get10Bytes14MPerSecond,
		"get11Bytes900MPerSecond":  Get11Bytes900MPerSecond,
		"get12BytesUltraPerSecond": Get12BytesUltraPerSecond,
	})
}
