package main
import(
    "encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/aberic/pkg/mod/github.com/hyperledger/fabric-contract-api-go/contractapi"
)
// 定义智能合约结构体
type SmartContract struct {
contractapi.Contract
}
// 定义查询结果集结构体
type QueryResult struct {
Key string `json:"Key"`
Record *QRcode
}
// 定义存储主体结构体
type QRcode struct {
//微信openid或访问ip
OpenId string `json:”open_id"`
//地理位置
Location string `json:"location"`
//QRcode文本
Text string `json:"text"`
//扫码时间
Time string `json:"time"`
//姓名
Name string `json:"name"`
//单位
Company string `json:"company"`
//是否发烧
Illness string `json:"illness"`
//旅游史
Travel string `json:"travel"`
}
//(s *SmartContract)意为为SmartContract实现接口功能
//即可通过contractapi.Contract.InitLedger访问此方法
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
qrcodes := []QRcode{
QRcode{OpenId: "127.0.0.1", Location: "HeNan,zhengzhou", Text: "Genesis block", Time: "2020-03-05", Name: "admin", Company: "郑州大学", Illness: "无", Travel: "无"},
}
for i, qrcode := range qrcodes {
//将数据编码成json字符串
QRcodeAsBytes, _ := json.Marshal(qrcode)
//将当前记录放入区块
err := ctx.GetStub().PutState("100000000"+strconv.Itoa(i), QRcodeAsBytes)
if err != nil {
return fmt.Errorf("Failed to put to world state. %s", err.Error())
}
}
fmt.Printf("Initialization successful！！！")
return nil
}
func (s *SmartContract) CreateQRcode(ctx contractapi.TransactionContextInterface, Key string, openId string, location string, text string, nowTime string, name string, company string,
illness string, travel string) error {
qrcode := QRcode{
OpenId: openId,
Location: location,
Text: text,
Time: nowTime,
Name: name,
Company: company,
Illness: illness,
Travel: travel,
}
QRcodeAsBytes, _ := json.Marshal(qrcode)
return ctx.GetStub().PutState(Key, QRcodeAsBytes)
}
func (s *SmartContract) QueryQRcode(ctx contractapi.TransactionContextInterface, Key string) (*QRcode, error) {
QRcodeAsBytes, err := ctx.GetStub().GetState(Key)
if err != nil {
return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
}
if QRcodeAsBytes == nil {
return nil, fmt.Errorf("%s does not exist", Key)
}
qrcode := new(QRcode)
//将json字符串解码到相应的数据结构
_ = json.Unmarshal(QRcodeAsBytes, qrcode)
return qrcode, nil
}
func (s *SmartContract) QueryAllQRcodes(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
startKey := "0"
endKey := "9999999999"
resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
if err != nil {
return nil, err
}
defer resultsIterator.Close()
results := []QueryResult{}
for resultsIterator.HasNext() {
queryResponse, err := resultsIterator.Next()
if err != nil {
return nil, err
}
qrcode := new(QRcode)
_ = json.Unmarshal(queryResponse.Value, qrcode)
queryResult := QueryResult{Key: queryResponse.Key, Record: qrcode}
results = append(results, queryResult)
}
return results, nil
}
func main() {
chaincode, err := contractapi.NewChaincode(new(SmartContract))
if err != nil {
fmt.Printf("Error create fabcar chaincode: %s", err.Error())
return
}
if err := chaincode.Start(); err != nil {
fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
}
}


