package storage

//func TestRead(t *testing.T) {
//
//	// Mock config
//	os.Args = []string{"cmd", "-config", "../testData/testConfig.yaml"}
//	config.Read(config.ParsedData)
//	logger.SetLogLevel(config.ParsedData.String("loglevel"))
//
//	err := config.ParsedData.Set("localStoragePath", "../testData/")
//	if err != nil {
//		t.Error(err.Error())
//	}
//
//	Configure()
//	result, err := Data.Read()
//	if err != nil {
//		t.Fatalf("Failed to read file" + config.ParsedData.String("localStoragePath") + localFileName)
//	} else {
//		t.Run("Read result exist", func(t *testing.T) {
//			if result.Data == nil {
//				t.Errorf("Expected data to be not nil, got nil")
//			}
//		})
//
//		t.Run("Read result region", func(t *testing.T) {
//			if result.Data[0].RegionName != "eu-west-1" {
//				t.Errorf("Expected region to be eu-west-1, got %s", result.AwsAccounts[0].RegionName)
//			}
//		})
//
//		t.Run("Read result ami", func(t *testing.T) {
//			if result.Data[0].VirtualMachines[0].ImageId != "ami-01c5300f289d64643" {
//				t.Errorf("Expected ami Id to be ami-01c5300f289d64643, got %s", result.Data[0].VirtualMachines[0].ImageId)
//			}
//		})
//	}
//
//}
