package main

import (
  "bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	//"sort"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//TxEvent define transaction response
type TxEvent struct {
	Txid string
}

func init() {
	log.SetFlags(log.Ltime)
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	switch function {
	case "initCopyright":
		return t.initCopyright(stub, args)
	case "readCopyright":
		return t.readCopyright(stub, args)
	case "queryCopyrightByField":
		return t.queryCopyrightByField(stub, args)
	case "getHistoryForCopyright":
		return t.getHistoryForCopyright(stub, args)
	//case "modifyCopyright":
	//	return t.modifyCopyright(stub, args)
	case "getTransactionID":
		return t.getTransactionID(stub, args)
	default:
		fmt.Println("invoke did not find func: " + function) //error
		return shim.Error("Received unknown function invocation, func=" + function + "\targs=" + args[0])
	}
}

// ============================================================
// initCopyright - create a new record, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initCopyright(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting input 2 parameters key and values for initCopyright")
	}

	key := args[0] // 文档的唯一标识 常量
	val := args[1]
	/*
		// ==== Check if worksId already exists ====
		copyrightAsBytes, err := stub.GetState(key)
		if err != nil {
			return shim.Error("Failed to get document bytes: " + err.Error())
		} else if copyrightAsBytes != nil {
			fmt.Println("This copyright already exists: " + key)
			return shim.Error("This copyright already exists: " + key)
		}
	*/

	// === Save marble to state ===
	err := stub.PutState(key, []byte(val))
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Marble saved and indexed. Return success ====
	fmt.Println(time.Now(), " - end init document")

	data, _ := json.Marshal(TxEvent{
		Txid: stub.GetTxID(),
	})
	return shim.Success(data)
}

/*
modifyCopyright modify copyright
*/
func (t *SimpleChaincode) modifyCopyright(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(time.Now().String() + " Incorrect number of arguments. Expecting input 2 parameters key and values for initCopyright")
	}

	key := args[0] // 文档的唯一标识 常量
	val := args[1]
	data, err := stub.GetState(key)
	if err != nil {
		return shim.Error(time.Now().String() + " Failed to get document bytes: " + err.Error())
	}
	if data == nil {
		return shim.Error(time.Now().String() + " Copyright need init before modify")
	}
	err = stub.PutState(key, []byte(val))
	if err != nil {
		return shim.Error(err.Error())
	}
	data, _ = json.Marshal(TxEvent{
		Txid: stub.GetTxID(),
	})
	return shim.Success(data)
}

// ===============================================
// readMarble - read a COPYRIGHT from chaincode state
// ===============================================
func (t *SimpleChaincode) readCopyright(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the copyright to query")
	}

	worksId := args[0]
	valAsbytes, err := stub.GetState(worksId) //get the COPYRIGHT from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + worksId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"worksID does not exist: " + worksId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===== Example: Parameterized rich query =================================================
// queryCopyrightByField 通过输入符合couchdb的查询字符串进行查询
// key: 要查询的字段名；value: 字段名对应的值
// =========================================================================================
func (t *SimpleChaincode) queryCopyrightByField(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	limit,err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	queryResults, err := getQueryResultForQueryString(stub, args[0], limit)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string, limit int) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	index := 0
	
	for resultsIterator.HasNext() {
		//queryResultKey, queryResultRecord, err := resultsIterator.Next()
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
		
		index += 1
		if (index >= limit) {
		    break
		}
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *SimpleChaincode) getTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]
	results, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer results.Close()

	txids := make([]string, 0)

	for results.HasNext() {
		item, err := results.Next()
		if err != nil {

			log.Printf("%s", err.Error())
			return shim.Error(err.Error())
		}
		txids = append(txids, item.GetTxId())
	}
	return shim.Success([]byte(strings.Join(txids, ",")))
}

// =========================================================================================
// getHistoryForCopyright executes the passed in string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func (t *SimpleChaincode) getHistoryForCopyright(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	worksId := args[0]

	fmt.Printf(time.Now().String()+"- start getHistoryForCopyright: %s\n", worksId)

	resultsIterator, err := stub.GetHistoryForKey(worksId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the copyright
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		//txID, historicValue, err := resultsIterator.Next()
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// historicValue is a JSON copyright, so we write as-is
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf(time.Now().String()+"- getHistoryForCopyright returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

