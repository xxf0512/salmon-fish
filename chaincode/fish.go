package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type SalmonFish struct {
	FishId          string    `json:"fish_id"`
	ProductorId     string    `json:"productor_id"`
	ProductorName   string    `json:"productor_name"`
	ProductionPlace string    `json:"production_palce"`
	ProductionTime  time.Time `json:"production_time"`
	ExpirationDate  time.Time `json:"expiration_date"`
	Number          int       `json:"number"`
	Weight          float64   `json:"weight"`
}

type SalmonFishChaincode struct {
}

func (t *SalmonFishChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")
	return shim.Success(nil)
}

func (t *SalmonFishChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fun, args := stub.GetFunctionAndParameters()
	if fun == "addFish" {
		return t.addFish(stub, args) // 添加信息
	} else if fun == "queryFishByFishId" {
		return t.queryFishByFishId(stub, args) // 根据FishId查询信息
	} else if fun == "queryFishByProductorId" {
		return t.queryFishByProductorId(stub, args) // 根据ProductorId查询信息
	}

	return shim.Error("指定的函数名称错误")
}

func PutFish(stub shim.ChaincodeStubInterface, fish *SalmonFish) ([]byte, bool) {
	var (
		b   []byte
		err error
	)
	if b, err = json.Marshal(*fish); err != nil {
		return nil, false
	}

	if err = stub.PutState(fish.FishId, b); err != nil {
		return nil, false
	}

	return b, true
}

func GetFishInfo(stub shim.ChaincodeStubInterface, fishId string) (*SalmonFish, bool) {
	var fish *SalmonFish
	// 根据身份证号码查询信息状态
	b, err := stub.GetState(fishId)
	if err != nil || b == nil {
		return fish, false
	}

	if err = json.Unmarshal(b, fish); err != nil {
		return fish, false
	}

	// 返回结果
	return fish, true
}

// 根据指定的查询字符串实现富查询
func GetFishByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// 添加信息 args: salmonFish
// fishId key, SalmonFish 为 value
func (t *SalmonFishChaincode) addFish(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var fish SalmonFish
	if err := json.Unmarshal([]byte(args[0]), &fish); err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	if _, exist := GetFishInfo(stub, fish.FishId); exist {
		return shim.Error("要添加鱼的ID已存在")
	}

	if _, bl := PutFish(stub, &fish); !bl {
		return shim.Error("保存信息时发生错误")
	}

	if err := stub.SetEvent(args[1], []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

// 根据fishId查询信息 args: fishId
func (t *SalmonFishChaincode) queryFishByFishId(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	fishId := args[0]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"fishId\":\"%s\"}}", fishId)
	queryString := fmt.Sprintf("{\"selector\":{\"fishId\":\"%s\"}}", fishId)

	// 查询数据
	result, err := GetFishByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据fishId查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的fishId没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据productorId查询信息 args: productorId
func (t *SalmonFishChaincode) queryFishByProductorId(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	productorId := args[0]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"productorId\":\"%s\"}}", productorId)
	queryString := fmt.Sprintf("{\"selector\":{\"productorId\":\"%s\"}}", productorId)

	// 查询数据
	result, err := GetFishByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据productorId查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的productorId没有查询到相关的信息")
	}
	return shim.Success(result)
}

func main() {
	err := shim.Start(new(SalmonFishChaincode))
	if err != nil {
		fmt.Printf("启动SalmonFishChaincode时发生错误: %s", err)
	}
}
