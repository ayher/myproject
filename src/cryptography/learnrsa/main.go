package main

import (
	"myproject/src/cryptography/learnrsa/src"
	"myproject/public/fmt"
)



func main()  {
	var err error
	//_,_,err=src.GetRsa()
	//if err!=nil {
	//	fmt.Error(err)
	//}
	err=src.Vitify(
		`-----BEGIN RSA PRIVATE KEY-----
MC0CAQACBQDiV2x5AgMBAAECBHwV3gECAwDtWQIDAPQhAgMAvnECAg4hAgMA7TY=
-----END RSA PRIVATE KEY-----`,
		`-----BEGIN RSA PUBLIC KEY-----
MCAwDQYJKoZIhvcNAQEBBQADDwAwDAIFAOJXbHkCAwEAAQ==
-----END RSA PUBLIC KEY-----`)
	if err!=nil {
		fmt.Error(err)
	}
}
