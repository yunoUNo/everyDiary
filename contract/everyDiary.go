// 패키지 정의
package main

// 1. 외부 모듈 포함
import (
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pr "github.com/hyperledger/fabric/protos/peer"
)

// 2. 체인코드 클래스-구조체 정의 everyDiary
type everyDiary struct{

}

// JSON 으로 변환할 구조체 정의
type Diary struct{
	Key 	string `json:"key"`
	Value 	string `json:"value"`
}

// 3. Init 함수
func (t *everyDiary) Init(stub shim.ChaincodeStubInterface) pr.Response{
	return shim.Success(nil)
}

// 4. Invoke 함수
func (t *everyDiary) Invoke(stub shim.ChaincodeStubInterface) pr.Response{
	fn, args := stub.GetFunctionAndParameters()
	
	switch fn{
	case "set":
		return t.Set(stub, args)
	case "get":
		return t.Get(stub, args)
	case "del":
		return t.Del(stub, args)
	case "history":
		return t.History(stub, args)
	case "checkUser":
		return t.CheckUser(stub)
	}

	return shim.Error("plz check function name")
}

// 5. Set 함수
func (t *everyDiary) Set(stub shim.ChaincodeStubInterface, args []string) pr.Response{

	if len(args) !=2{
		return shim.Error("plz check your args. must have 2 arguments(key,value).")
	}
	// 오류체크 중복 키 검사 -> 덮어쓰기로 해결
	diary := Diary{Key: args[0], Value: args[1]}
	diaryAsBytes, err := json.Marshal(diary)
	if err != nil{
		shim.Error("Failed set marshal args: " + args[0] + " " + args[1])
	}
	err = stub.PutState(args[0], diaryAsBytes)
	if err != nil{
		return shim.Error("set Failed!!! : " + args[0])
	}

	return shim.Success([]byte(diaryAsBytes))
}
// 6. Get 함수
func (t *everyDiary) Get(stub shim.ChaincodeStubInterface, args []string) pr.Response{

	if len(args) != 1{
		return shim.Error("plz check your args. must have 1 arguments(key).")
	}

	value, err := stub.GetState(args[0])   // key 로 블록 검사 하네.
	if err != nil{
		shim.Error("get Faield!!! : "+ args[0] + " with error: " + err.Error())
	}
	if value == nil{
		shim.Error("Diary not found: "+ args[0])
	}

	return shim.Success([]byte(value))
}


func (t *everyDiary) CheckUser(stub shim.ChaincodeStubInterface) pr.Response{

	IterVal, err := stub.GetStateByRange("","")
	if err != nil{
		return shim.Error(err.Error())
	}
	defer IterVal.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	flag := false
	for IterVal.HasNext(){
		queryResponse, err := IterVal.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if flag == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString("}")
		flag = true
	}
	buffer.WriteString("]")

	fmt.Printf("- AllCheck:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
// 7. Del 함수
func (t *everyDiary) Del(stub shim.ChaincodeStubInterface, args []string) pr.Response{
	if len(args) != 1{
		return shim.Error("plz check your args. must have 1 arguments(key).")
	}

	value, err := stub.GetState(args[0])
	if err != nil{
		shim.Error("get Faield!!! : "+ args[0] + " with error: " + err.Error())
	}
	if value == nil{
		shim.Error("incorrect key. Diary not found: "+args[0])
	}

	err = stub.DelState(args[0])

	return shim.Success([]byte(args[0]))
}

//9. History 함수
func (t *everyDiary) History(stub shim.ChaincodeStubInterface, args []string) pr.Response{

	if len(args) <1{
		return shim.Error("plz check your arguments. must have more than 1")
	}

	diaryName := args[0]

	fmt.Printf("- start getHistoryForDiary: %s\n", diaryName)

	resultsIterator, err := stub.GetHistoryForKey(diaryName)
	if err != nil{
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	flag := false
	for resultsIterator.HasNext(){
		response, err := resultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}
		if flag == true{
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		flag = true
	}
	buffer.WriteString("]")

	fmt.Printf("- History retuning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
// main 함수
func main(){
	if err := shim.Start(new(everyDiary)); err!= nil{
		fmt.Printf("Error run ChainCode: %s", err)
	}
}