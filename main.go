package main

import (
	"encoding/json"
	"fmt"
	"os"
	"salmon-fish/sdkInit"
	"salmon-fish/service"
	"salmon-fish/web"
	"salmon-fish/web/controller"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
	ROOTPATH   = "/home/pc433/xxf"
)

func main() {
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: ROOTPATH + "/salmon-fish/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: ROOTPATH + "/salmon-fish/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    ROOTPATH + "/salmon-fish/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    ROOTPATH + "/salmon-fish/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Printf(">> create chaincode lifecycle error: %v\n", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	fish := service.SalmonFish{
		FishId:          "123",
		ProductorId:     "abc",
		ProductorName:   "test",
		ProductionPlace: "beijing",
		Number:          10,
		Weight:          20,
	}

	serviceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	msg, err := serviceSetup.SaveFish(fish)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	result, err := serviceSetup.FindFishByFishId("123")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var fish service.SalmonFish
		json.Unmarshal(result, &fish)
		fmt.Println("根据FishId查询信息成功：")
		fmt.Println(fish)
	}

	app := controller.Application{
		Setup: serviceSetup,
	}
	web.WebStart(app)
}
