package comm

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/config"
)


type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

//生成优惠券凭证
func (backend *JWTAuthenticationBackend) GenerateCouponToken(openId string,couponCode string,trackCode string,orderNo string,couponAmount float64,appId string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["exp"] = time.Now().Add(time.Minute * time.Duration(config.GetValue("coupon_token_expiration").ToInt())).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["coupon_code"] = couponCode
	token.Claims["track_code"] = trackCode
	token.Claims["order_no"] = orderNo
	token.Claims["open_id"] = openId
	token.Claims["coupon_amount"] = couponAmount
	token.Claims["app_id"] = appId
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}



func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}




func getPrivateKey() *rsa.PrivateKey {


	privateKeyFile, err := os.Open(config.GetValue("privatekey_path").ToString())
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(config.GetValue("publickey_path").ToString())
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
