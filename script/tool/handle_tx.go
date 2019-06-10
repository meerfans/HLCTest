// Copyright 2017-2018 The nox developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package tool

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math"
	"os"
	"sort"
	"Nox-DAG-test/script/tool/json"
	"Nox-DAG-test/script/tool/hash"
	"Nox-DAG-test/script/tool/address"
	"Nox-DAG-test/script/tool/params"
	"Nox-DAG-test/script/tool/ecc"
	"Nox-DAG-test/script/tool/types"
	"Nox-DAG-test/script/tool/message"
	"Nox-DAG-test/script/tool/marshal"
	"Nox-DAG-test/script/tool/txscript"
)
type TxVersionFlag uint32
type TxLockTimeFlag uint32


type TxInputsFlag struct{
	inputs []TxInput
}
func (this *TxInputsFlag) SetFrom( txHash string ,index uint32 )  {
	tx := TxInput{}
	data, err :=hex.DecodeString(txHash)
	if err != nil{
		log.Fatalln("trans hash error",err.Error())
		os.Exit(0)
	}
	tx.txhash = data
	tx.index = index
	tx.sequence = uint32(math.MaxUint32)
	this.inputs = append(this.inputs,tx)
}

func (this *TxOutputsFlag) SetSmallOut( allCoinbase ,amount float64 , targets []string ,fromAdr string )  {
	tx := TxOutput{}
	//手续费
	tax := 0.004
	//按顺序将交易进行打包
	keys := make([]int,0)
	for k,_ := range targets{
		keys = append(keys,k)
	}
	sort.Ints(keys)
	for _,index := range keys{
		tx.amount = amount
		tx.target = targets[index]
		this.outputs = append(this.outputs,tx)
	}
	tx.amount = allCoinbase - tax - amount * float64(len(targets))
	tx.target = fromAdr
	this.outputs = append(this.outputs,tx)
}
type TxOutputsFlag struct{
	outputs []TxOutput
}
type TxInput struct {
	txhash []byte
	index uint32
	sequence uint32
}
type TxOutput struct {
	target string
	amount float64
}
func TxDecode(rawTxStr string) []byte{
	var param *params.Params
	param = &params.TestNetParams
	if len(rawTxStr)%2 != 0 {
		log.Fatalln(fmt.Errorf("invaild raw transaction : %s",rawTxStr))
	}
	serializedTx, err := hex.DecodeString(rawTxStr)
	if err != nil {
		log.Fatalln(err)
	}
	var tx types.Transaction
	err = tx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		log.Fatalln(err)
	}

	jsonTx := &json.OrderedResult{
		{"txid", tx.TxHash().String()},
		{"txhash", tx.TxHashFull().String()},
		{"version",  int32(tx.Version)},
		{"locktime", tx.LockTime},
		{"expire",tx.Expire},
		{"vin",      marshal.MarshJsonVin(&tx)},
		{"vout",     marshal.MarshJsonVout(&tx, nil,param)},
	}
	marshaledTx, err := jsonTx.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("\n%s\n",marshaledTx)
	return marshaledTx
}

func TxEncode(version TxVersionFlag, lockTime TxLockTimeFlag, txIn TxInputsFlag,txOut TxOutputsFlag) string{

	mtx := types.NewTransaction()

	mtx.Version = uint32(version)

	if lockTime!=0 {
		mtx.LockTime = uint32(lockTime)
	}

	for _, input := range txIn.inputs {
		txHash,err := hash.NewHashFromStr(hex.EncodeToString(input.txhash))
		if err!=nil{
			log.Fatalln(err)
		}
		prevOut := types.NewOutPoint(txHash, input.index)
		txIn := types.NewTxInput(prevOut, types.NullValueIn, []byte{})
		txIn.Sequence = input.sequence
		if lockTime != 0 {
			txIn.Sequence = types.MaxTxInSequenceNum - 1
		}
		mtx.AddTxIn(txIn)
	}

	//按顺序将交易进行打包
	keys := make([]int,0)
	for k, _:= range txOut.outputs{
		keys = append(keys,k)
	}
	sort.Ints(keys)
	for _, index:= range keys{
		output := txOut.outputs[index]
		// Decode the provided address.
		addr, err := address.DecodeAddress(output.target)
		if err != nil {
			log.Fatalln(errors.Wrapf(err,"fail to decode address %s",output.target))
		}

		// Ensure the address is one of the supported types and that
		// the network encoded with the address matches the network the
		// server is currently on.
		switch addr.(type) {
		case *address.PubKeyHashAddress:
		case *address.ScriptHashAddress:
		default:
			log.Fatalln(errors.Wrapf(err,"invalid type: %T", addr))
		}
		// Create a new script which pays to the provided address.
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			log.Fatalln(errors.Wrapf(err,"fail to create pk script for addr %s",addr))
		}

		atomic, err := types.NewAmount(output.amount)
		if err != nil {
			log.Fatalln(errors.Wrapf(err,"fail to create the currency amount from a " +
				"floating point value %f",output.amount))
		}
		//TODO fix type conversion
		txOut := types.NewTxOutput(uint64(atomic), pkScript)
		mtx.AddTxOut(txOut)
	}
	mtxHex, err := mtx.Serialize(types.TxSerializeNoWitness)
	if err != nil {
		log.Fatalln(err)
	}
	return hex.EncodeToString(mtxHex)
}

func TxSign(privkeyStr string, rawTxStr string) string{
	privkeyByte, err := hex.DecodeString(privkeyStr)
	if err!=nil {
		log.Fatalln("private key error",err,privkeyStr)
	}
	if len(privkeyByte) != 32 {
		log.Fatalln(fmt.Errorf("invaid ec private key bytes: %d",len(privkeyByte)))
	}
	privateKey, pubKey := ecc.Secp256k1.PrivKeyFromBytes(privkeyByte)
	h160 := hash.Hash160(pubKey.SerializeCompressed())

	var param *params.Params
	param = &params.TestNetParams
	addr,err := address.NewPubKeyHashAddress(h160,param,ecc.ECDSA_Secp256k1)
	if err!=nil {
		log.Fatalln(err)
	}
	// Create a new script which pays to the provided address.
	pkScript, err := txscript.PayToAddrScript(addr)
	if err!=nil {
		log.Fatalln(err)
	}

	if len(rawTxStr)%2 != 0 {
		log.Fatalln(fmt.Errorf("invaild raw transaction : %s",rawTxStr))
	}
	serializedTx, err := hex.DecodeString(rawTxStr)
	if err != nil {
		log.Fatalln(err)
	}

	var redeemTx types.Transaction
	err = redeemTx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		log.Fatalln(err)
	}
	var kdb txscript.KeyClosure
	kdb = func(types.Address) (ecc.PrivateKey, bool, error){
		return privateKey,true,nil // compressed is true
	}
	var sigScripts [][]byte
	for i,_:= range redeemTx.TxIn {
		sigScript,err := txscript.SignTxOutput(param,&redeemTx,i,pkScript,txscript.SigHashAll,kdb,nil,nil,ecc.ECDSA_Secp256k1)
		if err != nil {
			log.Fatalln(err)
		}
		sigScripts= append(sigScripts,sigScript)
	}

	for i2,_:=range sigScripts {
		redeemTx.TxIn[i2].SignScript = sigScripts[i2]
	}

	mtxHex, err := marshal.MessageToHex(&message.MsgTx{&redeemTx})
	if err != nil {
		log.Fatalln(err)
	}
	return mtxHex
	//fmt.Printf("%s\n",mtxHex)
}

