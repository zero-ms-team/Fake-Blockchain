package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func NewWallet() *Wallet {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return &Wallet{privateKey, publicKey}
}

func (w *Wallet) GetAddress() string {
	publicRIPEMD160 := HashPubKey(w.PublicKey)
	version := byte(0x00)

	return base58.CheckEncode(publicRIPEMD160, version)
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	return RIPEMD160Hasher.Sum(nil)
}
