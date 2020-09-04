package main

import (
	"testing"
)

const tpv1 = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDHZ+SKx2uhfrai
EmQI8SNNzlWnTBSXMJzvMHBnogcNYTMsEUOa9hYWRrKvWXQMubfVW41vxhMblH9b
ozChsJKYhIXPmXxLMmhfTUbh97UsOxBBTu54ZP9bwz1SxHdpOjIDfuAetiblWLiL
omYuh9DUpOm3F+l1OStPDfTdQNYhXt5cRE//5y7DLuvSlingHsprTvO84ia/nIBP
5ost5hLQ80Jyaj3SBOOAjYEw/R0fIt9tl5e5kfSe5eMjzyX/v55HlAkDv8NK1gb4
bpURHLsC9UUHUxeKIRQsfXhzzTuL8YCyyV7KG3NgvW9ZbVnCZZ3nmELMsgad5pEx
Z4pUENL7AgMBAAECggEBALaZ7vEe+PLkRH5Z9P0zRK8FWe5ffyOMQsnOQ8DC4U5h
Sij6jjwjScqQZySn99uHXk6lDfnjGrBQ5eeWovwN49CC2r5mwSljOay76UMYQPIG
DDah/0KEykrPmSJoAyl7Pz1wO/Ajwa6X9jb4OjY17QgtFFC0Nvc/qOc10puhufTH
ePwhknHTOH+Y5XAyiMNp8gJoAXQ8NHh/6eJW3W0BIHf/HXKfsucLwVkz2412y/Jt
GrFZbRUUnHYSyJ/VjOfUwHAsmWR77YNh3EL12GKEF4u5OGA4zhwWeACIzj13wJa6
V0QgkpmYvXiTxtHwlLjAE3EznwRbs5YkR12h4/5UMmECgYEA++YFNkbOoTslH5F7
Y8NuzfjYzAxdZIKk1XsOOVpqkeas0euEvx539Qkgnt2RK6svtrFr5PCi19nf0Olc
csVDrzLnnv2S2kpI/JOwopxa0FFvdhHoe+wno0Rf8lnzY62jt61s66bRXJpk98IH
s7twYe2USwj8UD7ORLb8zpbbcIMCgYEAyqcRqPvhdMC9ewgBIhJV/c+FeComweTW
SFPUsd9dPsIg4OCXEX0Mq2xSJvRzAne95tZGrmFsr+pRQ3OXGaZdScYExD86OBUF
MpKY9T6J7ODVkUm31zjF631Bvp61py2OWbgTjguD73lLL7l9qVTnyb66BpeQyfon
E77Bk1M/mikCgYA3TEiqoKKtzGEa7AINZZLGjrFxIenCrddnsgruVkX834niz3Ql
zJeC6E0L8xHyZzMjRRGtgZIOFptGrmQIIfv40xD72yjI2PPq1rU5DV/2SVpRrh6+
TZpqAhGaD1sZ7714Dg9SMB3X2WD+7s5oC2bhaJlcW42gRBleBlm7NGzZ5wKBgD5R
G8wkEJNvhZTkxDxu+QSAoSFvjNWJAh/hr4E3F5xp4+RjC/Fzy8aXG7gg6ZDzs3Dd
qYSMLvj1jCG61NctYniCLQsQCl4ekKeZjvGzVoSCKwpvadoD+lDNBr+QXHnZN3H9
ef3vKpYkbWtyleLRWimeveOzDfIeO5AF087zBZbpAoGABFeS5KjjYI6C2Jsdhh+r
Z4qVKLTGcyeLL30lDmCGVZZvbSN2/3rW3/vaf5/Bm42wVHTPnBN2IBFTH5gv99qc
xUORO3LUFlDzX/WSDDpUF5RbT3eE8IFrLAH0ouazyfR9KhV9w5qzhy8cGGEE+3gw
J6msr7d0ebh9hhGZpdGzifA=
-----END PRIVATE KEY-----`

const tpv2 = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC42LgfNYpnmZlm
H1Sh4fld0+C0CSQqFaGxnBSfRb5QVEHXoW56nreBX9XQ8PZVd+9zWFcRQ9uJXK6t
D89ty7Bpkt4m25cOlvekCSakNwsWfgM9Q43OZcsmx0tVxKvWuWTRzLsTUjpNuyKz
rY5lj6Qva293aQcaqeIwK6mgV+9mt7fH1wL4ulaAv8fKA7DPRQGnOmBKBhjUPojp
8XXj/hnjcw7SZXziriO7kouOa2Qx/vHnUVWSluiRj8xI6dN2r2D8sORXRx6Ww7Kw
zPobj+WQo6n1gFieqd9mVXty1l1HROeT5JkaIo/M4tG0jzvt8UQ+1FtQ8dJtgKJL
WLEpyvujAgMBAAECggEAW9UXTDH+S8/sfObN1gr9J1cvCGKOO/Y5OZLQT/yHO0NQ
3r7Ns0+M3VJuWHqN9xI5vPVDdElhFRIdmc8D/4k1YI3dnjJ0DoSEDVOahfTupkL+
ev5HwiWiUuaqD7dAo9HBO9XZjxTH8HwrFlcAKoa9S+21F/7uz5lczQ+0Gvv07PM6
hN1Sqw+jNH68wJWgMzpKtb+wlEk1PHq9uUSUOBKVqzco6Okfqp7eEHT2jMpJNQFN
tGKB9MOs53DhHVmiFns5H1/MkTdpxsuZwFzQxirf6Q/EwsEvDloF9SrXmlrrmty6
VUG70itvbtorIW/rctWOfUlWvsQjOBVuZBz98tUgYQKBgQDlCQP6xjkRg2LFqbpu
ugBQnF0rq3DchECtCWYlgaXODVlbqMk4SyXlTKpt55K+/5BWZjz2/NyTLofeFEIe
TqgzVCIGChJj+9EGOai3/ndyLCGVcjN3cERMyVW7A9YEBDbpu+E9GvqyOaAa7beY
E3HVAa26Wqh13rhrn+TLzOmzUQKBgQDOm+JMDonumgrzpgDWQ4M3Qe5R6HRHdNAG
olckOO8SivNWRcw5vlSGnh8ngW3g30nUCAA/C6+spUaHTi3jX0K9tJChS5AVyWnJ
3TfLvdmiq9Uujse8KOBmjKt9TwEOf2nhfb/XmSrTM7+OqB9hODBnjOk0mybKI0Na
yVXsTI56swKBgQCA7oNT/5SOvFS1CygNPx4AQxXcCIXfTYAPKNRc1tAc37zm8Wxd
CUjK/U6P0iX06W86hBFbxNry6+XGaccSwprDUmBY4ACcUlzH0VueQFzDY/5/36sD
WKrKQyjEv5MR7cFv8LkKKg7ol7H+lsWckY2qKGjBGFnvCLLuuzMUW0VQEQKBgDhU
b+JkpF6VSR8cx2WjiobqRtu2EN3aj0z/vdp2W1gm4ilHZmLn7Yu2WLAgraB9wFc6
xzZpLUBY313Mht5S/pNSQ4x2WZZXD6ylz6yQ2mFrj/fdnb9DNcs/1xGXFKarPmbo
LgHOFMr4dOWkGMoc07WnyX06P90kuBxsgCyowr/ZAoGBALZLdV8JAc0AG4zmTa83
pI0hGtdCC9vwhNCNqJRHr4kelJcRKZ5RUhu/FvJn+JFfPtvjms5B1HqIgITqVbhr
aJmXqeETECSuiJAfHgDaUFZUf9NUOusKRTcDKJ5/korHhC/gZ7GhJ12CG5Rw9kQ5
M6wifG2NyRq1asPEPvuwe+QL
-----END PRIVATE KEY-----`

const tpv3 = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAxl34Ljf1HoXG4Du8As3QFXo6QW8Zk1TQjQUIXdCXVVzNqQOp
jHKEgjKlBpUoXh5ivCdXUpWaG6MFkgfkPMbTzf0rx0+CWtA42EddBpfUm6u0JvJN
9L+QVcNSaB0QxgOr7kRtWwVU4DAYDwH0TkzQcXp9mpgWTmhZux3LPyLBVUZVxOoM
XaRkGT5tdVNpqfXA/aJYpPdCDHgPMs2Ad9JueeWPGHrTiPM+VrTQuwh9jdx1ZrIp
L2yDJKVw9TvBywfdwx5TxG2dAlnJw7ZUG07z44yN7RHjq4Q1Ponh9PFpE4ahWoYy
DKtEVWiZ/gYau8lGYCZaCiBJv0XUyrjBlwJRlwIDAQABAoIBAC4S6XAiwyZBwD2Y
4kRsdWZnq1xDZ9dWndIDVzwjFONY2NPm69yZRLnQ9Y+f2W5y9E/re1bNNKwyozqV
5qdPLybrJN34U7iwIGPrv2mmnlVB/mgFY6HzVJz27w8LoTe85pRDBHtds6cWxJ/H
vmzpXojdAIuFw7iWeDED7I9PjUQ21qbTxTqDC3dnFjX6NX4dWtI+SgawVPbdlSTq
QG69eSto9bH4FUqPVDBv5uOirBhqZYcaAON2WMW33/l3PZXQV5lGh3CRZD5+HYxv
6pS7Qb+pyTQkncIhGQoE/n7vrXzrchcXlgzL+OBzAahLOJvodYboF5GqeKIuNz7k
XmhbjcECgYEA+6THZKaza9lpFseUB189gFbq0MgKIPhSqdy30RcOBvts/hEwURaE
sYuHPM/JDwqjoO7eO7+w/WhD37wXVkoNc3xEJw6i/pQqLMh7GmcDz4aK46YXeCUM
GHFPaywTvMIJTa1br/vn+rnZGd/ydXe/ySCoDDVQH2kQFuGuz5KRynECgYEAyc0V
CCSWqKBkQy+S3ANuipbBR83z1Bm4aqRumEDzF4cLvKv5bT7s5/MrTUrHsZjf/93I
9mR6O13aukUHX/L4X42vv4tc9BGDmbWdY/CJmFtVIko+bsZ/MKC6XwoD16F+JpiE
706m3vnXyArBUwIGtP4lXRscWEWdAzbpXUTRkIcCgYBhrrU3/P3o/5wrm416zx2w
lAzSvtQvuDVeeq9gGvL3AuJsPX/j+jnIMcFtebsye47JCfB6gQ7TT4YJc5obhONz
0OkjwCrFZ/53I9ulhBeWl0OS2waBPOBVHKcXkySWQTwbSxAsYDzMtxfvU19q+fEY
wfR5yLgxeTclqrWRHfQ6AQKBgHBed0CihxX8wfe7bO6AJrSbP6MJJqXLcKpJR6AW
QoauVzXHGUvgxzBdcpZGdq4I72pdiELTLlEScPJZ78JY3D7w+ZUSOD9b5UjZHXwB
+8xPxzch2mP6ueZNCZpUTFFtBn7dXOCYjkkJHEOy4XWkYjG0dv/CUeVBVi3tDMM3
x+3PAoGAOasPunUej3vm1JSITVridE9yVPeq/QrYL8h5bLxcj6EvDd4jHRm+LcRJ
sA/macDaeQIMZJIm9f0cGqptYZ8TnQUZ/gWozLoLiu2ol1PAQZ+MyCMJ0OKJJ+2J
4DYg+7ca7zb5uTdJD3ZBxTvs704eIro0QiHrxE0aYBQCI9ajaNc=
-----END RSA PRIVATE KEY-----`

func TestLoadPrivateKey(t *testing.T) {
	_, err := LoadPrivateKey([]byte(tpv1))
	if err != nil {
		t.Error(err)
	}

	_, err2 := LoadPrivateKey([]byte(""))
	if err2 == nil {
		t.Errorf("Must be an error")
	}

	_, err3 := LoadPrivateKey([]byte(tpv3))
	if err3 == nil {
		t.Errorf("Must be an error")
	}
}

func TestGenerateToken(t *testing.T) {
	rsa1, _ := LoadPrivateKey([]byte(tpv1))
	rsa2, _ := LoadPrivateKey([]byte(tpv2))
	subject := "sample"
	key := "foo"
	key2 := "foo1"
	token, err := GenerateToken(subject, 0, 2, key, rsa1)

	if err != nil {
		t.Error(err)
	}
	verify1, err := VerifyToken(token, key, rsa1)
	if err != nil {
		t.Error(err)
	}
	verify2, err := VerifyToken(token, key, rsa1)
	if err != nil {
		t.Error(err)
	}
	if verify1 != subject || verify2 != subject {
		t.Errorf("Verification failed")
	}

	_, err2 := VerifyToken(token, key, rsa2)
	if err2 == nil {
		t.Errorf("Error must be occurs")
	}
	_, err3 := VerifyToken(token, key2, rsa1)
	if err3 == nil {
		t.Errorf("Error must be occurs")
	}
	_, err4 := VerifyToken("somthing", key, rsa1)
	if err4 == nil {
		t.Errorf("Error must be occurs")
	}

}

func TestAuthToken(t *testing.T) {
	rsa1, _ := LoadPrivateKey([]byte(tpv1))

	authToken := AuthToken{}
	authToken.Checksum = "random-checksum"
	authToken.Type = ChallengeTypeJS

	salt := "foo"

	tokenValue, err := GenerateAuthToken(authToken, 0, salt, rsa1)
	if err != nil {
		t.Error(err)
	}

	authTokenOut, err := ValidateAuthToken(tokenValue, salt, rsa1)
	if err != nil {
		t.Error(err)
	}

	if authToken.Checksum != authTokenOut.Checksum {
		t.Errorf("Seems be not the same")
	}

	_, err2 := ValidateAuthToken(tokenValue, "what?", rsa1)
	if err2 == nil {
		t.Errorf("Error must be occurred")
	}
}

// func TestGenerateToken(t *testing.T) {
// 	subject := "1.1.1.1"
// 	key := "foo"
// 	token, err := GenerateToken(subject, 1, key)
// 	fmt.Println(token)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	verify, err := VerifyToken(token, key)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if verify != subject {
// 		t.Errorf("Must be same")
// 	}
// 	_, err2 := VerifyToken("sampleing", key)
// 	if err2 == nil {
// 		t.Errorf("Must be error")
// 	}
// 	if verify == "" {
// 		t.Errorf("Verification failed")
// 	}
// 	time.Sleep(3 * time.Second)
// 	verify, err = VerifyToken(token, key)
// 	if verify != "" {
// 		t.Error(err)
// 		t.Errorf("Verification must failed")
// 	}
// }

// func BenchmarkGenerateToken(b *testing.B) {
// 	subject := "1.1.1.1"
// 	key := "foo"
// 	for i := 0; i < b.N; i++ {
// 		GenerateToken(subject, 1, key)
// 	}
// }
// func BenchmarkVerifyToken(b *testing.B) {
// 	subject := "1.1.1.1"
// 	key := "foo"
// 	for i := 0; i < b.N; i++ {
// 		token, _ := GenerateToken(subject, 1, key)
// 		VerifyToken(token, key)
// 	}
// }
