package myTest

import (
	"encoding/json"
	"io/ioutil"
	"myproject/public/fmt"
	"myproject/public/msg"
	"myproject/src/bitcoin/public"
	"net/http"
	"strconv"
	"time"
)

type AddressBalance struct {
	Balance float64 `json:"balance"`
}

type AddressInner struct {
	Transactions []interface{} `json:"transactions"`
	Address AddressBalance `json:"address"`
}

type AddressData struct {
	Data map[string]interface{} `json:"data"`
}

func runFun(num int,address string) (bool,float64) {
	fmt.Println("num : ",num)
	if num==0{
		return vitifyAddress0(address)
	}else if num==1{
		return vitifyAddress1(address)
	}else if num==2{
		return vitifyAddress2(address)
	}else{
		return vitifyAddress3(address)
	}
}

func Test()  {
	i :=0
	for{
		if i>100000000{
			i=1
		}
		i++

		f:=i%4
		pri,pub:= public.NewKey()
		address:= public.PubkeyToAddress(pub)
		b,v:=vitifyAddress1(address)
		fmt.Println(time.Now())
		fmt.Println(pri)
		fmt.Println(address)

		b,v=runFun(f,address)

		if b{
			fmt.Println("-----",v,b)
			fmt.Println(msg.SendLarkNotify("a5383493-cbd9-4a56-a797-45f0390a1c75","momey",address))
		}
		if v==-1{
			fmt.Println("qie huan")
			f++
			b,v=runFun(f%4,address)
			if b{
				fmt.Println("-----",v,b)
				fmt.Println(msg.SendLarkNotify("a5383493-cbd9-4a56-a797-45f0390a1c75","momey",address))
			}
		}
	}
}

func vitifyAddress0(address string) (bool,float64) {
	resp,err:=http.Get("https://api.blockchair.com/bitcoin/dashboards/address/"+address)
	if err!=nil{
		fmt.Println("1 err:",err)
		return false,-1
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("2 err:",err)
		return false,-1
	}
	var addressData AddressData
	json.Unmarshal(body,&addressData)

	addressInner,ok:=addressData.Data[address].(map[string]interface{})
	if !ok{
		fmt.Println("3 err")
		return false,-1
	}
	addressBalance,ok:=addressInner["address"].(map[string]interface{})
	if !ok{
		fmt.Println("4 err")
		return false,-1
	}
	balance,ok:=addressBalance["balance"].(float64)
	if !ok{
		fmt.Println("5 err")
		return false,-1
	}
	return balance!=0.0,balance
}

func vitifyAddress1(address string) (bool,float64) {
	resp,err:=http.Get("https://api.blockcypher.com/v1/btc/main/addrs/"+address+"/balance")
	if err!=nil{
		fmt.Println("1 err:",err)
		return false,-1
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("2 err:",err)
		return false,-1
	}
	var addressData map[string]interface{}
	json.Unmarshal(body,&addressData)

	balance,ok:=addressData["balance"].(float64)
	if !ok{
		fmt.Println("3 err")
		return false,-1
	}
	return balance!=0.0,balance
}

// https://blockchain.coinmarketcap.com/api/address?address=1Cs72gQRPp1aGYdJCXfVSnQknpch6p12iK&symbol=BTC&start=1&limit=10
func vitifyAddress2(address string) (bool,float64) {
	resp,err:=http.Get("https://blockchain.coinmarketcap.com/api/address?address="+address+"&symbol=BTC&start=1&limit=1")
	if err!=nil{
		fmt.Println("1 err:",err)
		return false,-1
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("2 err:",err)
		return false,-1
	}
	var addressData map[string]interface{}
	json.Unmarshal(body,&addressData)

	b,ok:=addressData["balance"].(string)
	if !ok{
		fmt.Println("3 err")
		return false,-1
	}
	balance, err := strconv.ParseFloat(b, 64)
	if err!=nil{
		fmt.Println("4 err")
		return false,-1
	}
	return balance!=0.0,balance
}

// https://btc5.trezor.io/api/v2/address/1PQwtwajfHWyAkedss5utwBvULqbGocRpu
func vitifyAddress3(address string) (bool,float64) {
	resp,err:=http.Get("https://btc5.trezor.io/api/v2/address/"+address)
	if err!=nil{
		fmt.Println("1 err:",err)
		return false,-1
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("2 err:",err)
		return false,-1
	}
	var addressData map[string]interface{}
	json.Unmarshal(body,&addressData)

	b,ok:=addressData["balance"].(string)
	if !ok{
		fmt.Println("3 err")
		return false,-1
	}
	balance, err := strconv.ParseFloat(b, 64)
	if err!=nil{
		fmt.Println("4 err")
		return false,-1
	}
	return balance!=0.0,balance
}