package tool

import (
	"fmt"
	"log"
	"encoding/hex"
	"Nox-DAG-test/script/tool/seed"
	"Nox-DAG-test/script/tool/bip32"
	"Nox-DAG-test/script/tool/ecc"
	"Nox-DAG-test/script/tool/address"
	"Nox-DAG-test/script/tool/hash"
	"Nox-DAG-test/script/tool/base58"
	"Nox-DAG-test/script/tool/params"
)

//create address
func CreateNoxAddr(network string) (priKey string ,base58Addr string ){
	seed1 := NewEntropy(32)
	//log.Println("【rand seed】",seed)
	privateKey := EcNew("secp256k1",seed1)
	log.Println("【HLC private key】",privateKey)
	publicKey := EcPrivateKeyToEcPublicKey(false,privateKey)
	//log.Println("【public key】",publicKey)
	param := params.PrivNetParams
	switch network {
	case "private":
		break
	case "test":
		param = params.TestNetParams
		break
	case "main":
		break
	}

	//param := params.MainNetParams
	addr := EcPubKeyToAddress(param.PubKeyHashAddrID[:],publicKey)
	addres,err := address.DecodeAddress(addr)
	if err != nil{
		log.Fatalln("【verify failed】",err)
		return
	}
	log.Println("【HLC base58 address】",addres)
	return privateKey,addres.String()
}
//generate seed
func NewEntropy(size uint) string{
	s,err :=seed.GenerateSeed(uint16(size))
	if err!=nil {
		log.Fatal(err)
		return ""
	}
	return fmt.Sprintf("%x",s)
}
//secp256k1 generate private key
func EcNew(curve string, entropyStr string) string{
	entropy, err := hex.DecodeString(entropyStr)
	if err!=nil {
		log.Fatalln("【error】",entropyStr,err)
		return ""
	}
	switch curve {
	case "secp256k1":
		masterKey,err := bip32.NewMasterKey(entropy)
		if err!=nil {
			log.Fatalln(err)
			return ""
		}
		return fmt.Sprintf("%x",masterKey.Key[:])
	default:
	}
	return ""
}

//from private key to public key
func EcPrivateKeyToEcPublicKey(uncompressed bool, privateKeyStr string) string{
	data, err := hex.DecodeString(privateKeyStr)
	if err!=nil {
		log.Fatalln(err)
		return ""
	}
	_, pubKey := ecc.Secp256k1.PrivKeyFromBytes(data)
	var key []byte
	if uncompressed {
		key = pubKey.SerializeUncompressed()
	}else {
		key = pubKey.SerializeCompressed()
	}
	return fmt.Sprintf("%x",key[:])
}

// public key to bas58 address
func EcPubKeyToAddress(version []byte, pubkey string) string{
	data, err :=hex.DecodeString(pubkey)
	if err != nil {
		log.Println(err)
		return ""
	}
	h := hash.Hash160(data)

	addr := base58.NoxCheckEncode(h, version[:])
	return fmt.Sprintf("%s",addr)
}

