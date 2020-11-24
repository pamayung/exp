package cache

import "galeri24/model"

var m = make(map[string]interface{})

func SetOtp(barcode string, o *model.OTP) {
	m[barcode] = o
}

func GetOtp(barcode string) interface{} {
	return m[barcode]
}
