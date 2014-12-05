package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
)

func createSshKey() string {

	println("Creating SSH Keys")

	privatekey, err := rsa.GenerateKey(rand.Reader, 2014)

	if err != nil {
		panic(err)
	}

	publicKey := &privatekey.PublicKey

	pkey := x509.MarshalPKCS1PrivateKey(privatekey)
	ioutil.WriteFile("private.key", pkey, 0777)

	pubkey, _ := x509.MarshalPKIXPublicKey(publicKey)
	ioutil.WriteFile("public.key", pubkey, 0777)

	return string(pubkey)

}
