package main

import (
	"testing"
	"dag-bench-test-script/tool"
	"log"
	"encoding/json"
	"fmt"
	"Nox-DAG-test/script/tool/types"
)

func TestDecode(t *testing.T)  {
	var res map[string]interface{}
	rawHash := "01000000016fe96dd5f988ab6a749e76dd844d292102ed61a1ceaf0ff2442e8ba832d6d00b02000000ffffffff0500a3e111000000001976a9142d5c76e7d22193f4c66a7a3f07ac9020891be69988ac00a3e111000000001976a9147b2d2ae5cd6566ce2dc1a98bb8c9e26797c796cb88ac00a3e111000000001976a914e3c2784f6efe4532eb57b1fc8e4ae9457345931488ac00a3e111000000001976a914895f7ada8141b7dccc467e75901a73bb5e859b3988ac0036773e000000001976a914864c051cdb39c31f21924a5ac88b4cf82124d2c188ac00000000000000000180461c860000000013000000000000006b4830450221008d299a077b8d779295a314b92bfc3c4899ee112b23e4ecb7ff2d9e8c8d4b9ea702201d87c4a05d2fb19ba84d75fd118a442cc21b6eac8f00da62167e8b816850c42d012102b3e7c21a906433171cad38589335002c34a6928e19b7798224077c30f03e835e"
	s := tool.TxDecode(rawHash)
	json.Unmarshal(s,&res)
	log.Println(res["txid"])
	for index,trans := range res["vout"].([]interface{}){
		log.Println(index,trans)
	}
}

func TestCreateAddresses(t *testing.T){
	addresses := make([]string,0)
	adrContent := make([][]string,0)
	for j:=0;j<=30000;j++ {
		pk,addr := tool.CreateNoxAddr("private")
		addresses = append(addresses,addr)
		adrContent = append(adrContent,[]string{pk,addr})
	}
	tool.WriteCsv("/tmp/address.csv",adrContent)
}


func TestReadCsv(t *testing.T){
	result := tool.ReadCsv("/tmp/address.csv",0,1000)
	fmt.Println(result)
}



func TestMax(t *testing.T){
	fmt.Println(types.MaxAmount)
}


