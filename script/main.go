package main

import (
	"log"
	"dag-bench-test-script/tool"
	"encoding/json"
	"sort"
	"strings"
	"os"
	"fmt"
)

type SendRawResult struct {
	JsonRpc string `json:"json_rpc"`
	ID int `json:"id"`
	Result string `json:"result"`
}

var TxID string
var CoinbaseAmount float64
var TxRAW []byte

func main()  {
	cfg, _, err := tool.LoadConfig()
	if err != nil {
		log.Fatal("Config file error,please check it.【",err,"】")
		return
	}

	rpc := &tool.RpcClient{}
	rpc.Cfg = cfg
	switch cfg.Action {
	case "generate-new-address":
		log.Println("Create new HLC Address")
		pk,addr := tool.CreateNoxAddr(cfg.Network)
		log.Println("【HLC private key】",pk)
		log.Println("【HLC base58 address】",addr)
		break
	case "batch-generate-address-signed-transactions":

	default:
		if cfg.FromAddress == ""{
			log.Fatalln("please set the FromAddress ! --faddress")
			os.Exit(0)
		}
		if cfg.FromPrivateKey == ""{
			log.Fatalln("please set the FromAddress private key!--privkey")
			os.Exit(0)
		}
		if cfg.RPCServer == ""{
			log.Fatalln("please set the RpcServer !-s")
			os.Exit(0)
		}
		log.Println("Batch generate addresses and signed transactions")
		//get the txid by heigth witch can spend
		GetTxID(cfg,rpc)
		//create 999 rand hlc address
		CreateAddresses(cfg)
		//spend coinbase money
		CreateMoneyAccounts(cfg,rpc)
		//create 999 tx
		SendMoneyAccounts(cfg,rpc)
	}
	
}

func GetTxID(cfg *tool.Config,rpc *tool.RpcClient){
	body := rpc.RpcResult("getBlockByOrder",[]interface{}{cfg.Height,true})
	//log.Println(string(body))
	if strings.Contains(string(body),"error"){
		log.Fatalln("not exist height block!",cfg.Height)
		os.Exit(0)
	}
	var res map[string]interface{}
	json.Unmarshal(body,&res)
	if res == nil{
		log.Fatalln("rpc error,please check!",cfg.Height)
		os.Exit(0)
	}
	result := res["result"].(map[string]interface{})
	transactions := result["transactions"].([]interface{})
	if len(transactions)<1{
		log.Fatalln("no transactions")
		os.Exit(0)
	}
	tx := transactions[0].(map[string]interface{})
	outs := tx["vout"].([]interface{})
	if len(outs) < 3{
		log.Fatalln("not have coinbase")
		os.Exit(0)
	}
	coinbase := outs[2].(map[string]interface{})
	scriptPubKey := coinbase["scriptPubKey"].(map[string]interface{})
	//log.Println(scriptPubKey["addresses"])
	addrs := scriptPubKey["addresses"].([]interface{})
	if !tool.InArrayString(cfg.FromAddress,addrs){
		log.Fatalln("the coinbase is belong to the account")
		os.Exit(0)
	}
	cfg.FromTransactionHash = tx["txid"].(string)
	CoinbaseAmount = coinbase["amount"].(float64)
	cfg.AddressFile = fmt.Sprintf("%s%d",cfg.AddressFile,cfg.Height)
	cfg.TXFile = fmt.Sprintf("%s%d",cfg.TXFile,cfg.Height)
}

func CreateAddresses(cfg *tool.Config){
	addresses := make([]string,0)
	adrContent := make([][]string,0)
	for j:=0;j<=998;j++ {
		pk,addr := tool.CreateNoxAddr(cfg.Network)
		addresses = append(addresses,addr)
		adrContent = append(adrContent,[]string{pk,addr})
	}
	tool.WriteCsv(cfg.AddressFile,adrContent)
	log.Println("999 address success！")
}
func CreateMoneyAccounts(cfg *tool.Config,rpc *tool.RpcClient){
	addresses := make([]string,0)
	csvContent := tool.ReadCsv(cfg.AddressFile,0,999)
	for _,v:=range csvContent{
		addresses = append(addresses,v[1])
	}
	log.Println(CoinbaseAmount/100000000.00)
	SendRawTxHash(CoinbaseAmount/100000000.00,cfg.FromPrivateKey,cfg.FromTransactionHash,cfg.FromAddress,addresses,rpc,0.02)
}
//spend coinbase
func SendRawTxHash(allCoinbase float64,fromPK string ,fromTxHash string, fromAddr string ,toAddrs []string ,rpc *tool.RpcClient,amout float64) string {
	version := tool.TxVersionFlag(1)
	locktime := tool.TxLockTimeFlag(0)
	//build tx in out
	txLn := tool.TxInputsFlag{}
	//coinbase trx
	txLn.SetFrom(fromTxHash,2)

	txOut := tool.TxOutputsFlag{}
	txOut.SetSmallOut(allCoinbase,amout,toAddrs,fromAddr)

	//raw txhash
	newTxRaw := tool.TxEncode(version,locktime,txLn,txOut)
	TxRAW = tool.TxDecode(newTxRaw)
	//sign
	signHash := tool.TxSign(fromPK,newTxRaw)

	//send tx
	result := rpc.RpcResult("sendRawTransaction",[]interface{}{signHash})
	log.Println(string(result))
	if strings.Contains(string(result),"error"){
		log.Fatalln("generate account failed！")
		os.Exit(0)
	}
	var res SendRawResult
	json.Unmarshal(result,&res)
	TxID = res.Result
	return signHash
}


// 0.01 amount ,  hlc 0 addr => 1 addr 1 addr => 2 addr ... 998 addr => 0 addr
func SendMoneyAccounts(cfg *tool.Config,rpc *tool.RpcClient){

	csvContent := tool.ReadCsv(cfg.AddressFile,0,999)
	csvContent1 := make([][]string,0)
	keys := make([]int,0)
	for k,_:=range csvContent{
		keys = append(keys,k)
	}
	sort.Ints(keys)
	for _,index := range keys{
		pk := csvContent[index][0]
		fromAddr := csvContent[index][1]
		addresses := []string{}
		if index < 998{
			addresses = []string{csvContent[index+1][1]}
		} else{
			addresses = []string{csvContent[0][1]}
		}
		signRawhash := SendTransaction(cfg,0.02,pk,TxID,fromAddr,addresses,rpc,0.01,uint32(index))
		csvContent1 = append(csvContent1,[]string{signRawhash})
	}
	tool.WriteCsv(cfg.TXFile,csvContent1)
	log.Println("complete 999 txs:",cfg.TXFile)
}

func SendTransaction(cfg *tool.Config,allCoinbase float64,fromPK string ,fromTxHash string, fromAddr string ,toAddrs []string ,rpc *tool.RpcClient,amout float64,txIndex uint32) string {
	fromPK = strings.Replace(fromPK, " ", "", -1)
	fromPK = strings.Replace(fromPK, "\xEF\xBB\xBF", "", -1)
	fromPK = strings.Replace(fromPK, "\n", "", -1)
	fromAddr = strings.Replace(fromAddr, " ", "", -1)
	fromAddr = strings.Replace(fromAddr, "\n", "", -1)
	version := tool.TxVersionFlag(1)
	locktime := tool.TxLockTimeFlag(0)
	//build tx in out
	txLn := tool.TxInputsFlag{}
	txLn.SetFrom(fromTxHash,txIndex)

	txOut := tool.TxOutputsFlag{}
	txOut.SetSmallOut(allCoinbase,amout,toAddrs,fromAddr)

	//raw txhash
	newTxRaw := tool.TxEncode(version,locktime,txLn,txOut)
	//sign
	signHash := tool.TxSign(fromPK,newTxRaw)
	//send tx
	if cfg.Send{
		result := rpc.RpcResult("sendRawTransaction",[]interface{}{signHash})
		log.Println(string(result))
		var res SendRawResult
		json.Unmarshal(result,&res)
	}

	return signHash
}


