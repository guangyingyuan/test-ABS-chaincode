// 债券基础信息相关的链码操作

/*

*/

package main


import (
	"fmt"
	"encoding/json"
  "time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("BusinessPartnerInfo")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// BusinessPartnerInfo struct
// ============================================================================================================================
type BusinessPartnerInfoStruct struct {
	UserName             string `json:"UserName"`
	Organization         string `json:"Organization"`
	Company              string `json:"Company"`
	Account              string `json:"Account"`
	CreatedTime          time.Time `json:"CreatedTime"`
	OperateLog           []string `json:"OperateLog"`
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	logger.Info("########### BusinessPartnerInfo Init ###########")
	return shim.Success(nil)


}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### BusinessPartnerInfo Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
  if function == "add" {
		// Deletes an entity from its state
		return t.add(stub, args)
	}
	if function == "addOperateLog" {
		// Deletes an entity from its state
		return t.addOperateLog(stub, args)
	}
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if function == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	if function == "update" {
		// Deletes an entity from its state
		return t.update(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2. ")
	}
  var PartnerInfoObj BusinessPartnerInfoStruct
	PartnerInfo :=args[0]
  err = json.Unmarshal([]byte(PartnerInfo),&PartnerInfoObj)
	if err != nil {
	fmt.Println("error:", err)
	return shim.Error(err.Error())
	 }
	UserName := PartnerInfoObj.UserName
	UserTest, _ := stub.GetState(UserName)
	if UserTest != nil {
		return shim.Error("the user is existed")
	}
	timestamp, _:= stub.GetTxTimestamp()
	PartnerInfoObj.CreatedTime = time.Unix(timestamp.Seconds, int64(timestamp.Nanos))
	var OperateLog []string
	TxID := stub.GetTxID()
	OperateLog =append(OperateLog,TxID)
	PartnerInfoObj.OperateLog = OperateLog
	jsonAsBytes,_:= json.Marshal(PartnerInfoObj)
	err = stub.PutState(UserName,[]byte(jsonAsBytes))
	if err != nil {
		return shim.Error(err.Error())
	}
  return shim.Success(nil);
}

func (t *SimpleChaincode) addOperateLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. ")
	}
	var PartnerInfoObj BusinessPartnerInfoStruct
	UserName := args[0]
	NewOperateLog := args[1]
	OldUserInfo, _ := stub.GetState(UserName)
	if OldUserInfo == nil {
		return shim.Error("the user is not exist!!")
	}
	err = json.Unmarshal([]byte(OldUserInfo),&PartnerInfoObj)
	if err != nil {
		return shim.Error(err.Error())
	}
	PartnerInfoObj.OperateLog = append(PartnerInfoObj.OperateLog,NewOperateLog)
	jsonAsBytes,_:= json.Marshal(PartnerInfoObj)
	err = stub.PutState(UserName,[]byte(jsonAsBytes))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil);
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!= 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 ")
		}
	UserName := args[0]
	UserInfo, err := stub.GetState(UserName)
	//test if the user has been existed
	if err != nil {
		return shim.Error("The user never been exited")
	}
	if UserInfo == nil {
	return shim.Error("The user`s information is empty!")
 }
 err = stub.DelState(UserName) //remove the key from chaincode state
 	if err != nil {
 		return shim.Error("Failed to delete the user. ")
 	}
 return shim.Success(nil)

}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	 }
	 UserName := args[0]
 	UserInfo, err := stub.GetState(UserName)
 if err != nil {
 	jsonResp := "{\"Error\":\"Failed to get state for " + UserName + "\"}"
 	return shim.Error(jsonResp)
  }
 if UserInfo == nil {
 	jsonResp := "{\"Error\":\"Nil content for " + UserName + "\"}"
 	return shim.Error(jsonResp)
  }
 return shim.Success(UserInfo)

}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2. ")
	}
	var NewPartnerInfoObj BusinessPartnerInfoStruct
	var OldPartnerInfoObj BusinessPartnerInfoStruct
	NewPartnerInfo :=args[0]
	err = json.Unmarshal([]byte(NewPartnerInfo),&NewPartnerInfoObj)
	if err != nil {
	fmt.Println("error:", err)
	return shim.Error(err.Error())
	 }
	UserName := NewPartnerInfoObj.UserName
	OldUserInfo, _ := stub.GetState(UserName)
	if OldUserInfo == nil {
		return shim.Error("the user is not exist!!")
	}
	err = json.Unmarshal([]byte(OldUserInfo),&OldPartnerInfoObj)
	NewPartnerInfoObj.CreatedTime = OldPartnerInfoObj.CreatedTime
	NewPartnerInfoObj.OperateLog = OldPartnerInfoObj.OperateLog
	jsonAsBytes,_:= json.Marshal(NewPartnerInfoObj)
	err = stub.PutState(UserName,[]byte(jsonAsBytes))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil);
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
