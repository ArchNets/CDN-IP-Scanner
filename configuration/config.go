package config

import (
	"CFScanner/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	PROGRAMDIR, _          = filepath.Abs(filepath.Dir(os.Args[0]))
	DIR                    = filepath.Join(PROGRAMDIR, "config")
	RESULTDIR              = filepath.Join(PROGRAMDIR, "result")
	StartDtStr             = time.Now().Format("2006-01-02_15-04-05")
	CSVInterimResultsPath  = filepath.Join(RESULTDIR, StartDtStr+"_result.csv")
	JSONInterimResultsPath = filepath.Join(RESULTDIR, StartDtStr+"_result.json")
	FinalResultsPathSorted = filepath.Join(RESULTDIR, StartDtStr+"_final.txt")
)

func (C Configuration) PrintInformation() {
	fmt.Printf(`-------------------------------------
Configuration :
User ID : %v%v%v
XHTTP Host: %v%v%v
XHTTP Path : %v%v%v
Address Port : %v%v%v
Upload Test : %v%v%v
Fronting Request Test : %v%v%v
Minimum Download Speed : %v%v%v
Maximum Download Time : %v%v%v
Minimum Upload Speed : %v%v%v
Maximum Upload Time : %v%v%v
Test URL : %v%v%v
Fronting Timeout : %v%v%v
Maximum Download Latency : %v%v%v
Maximum Upload Latency : %v%v%v
Number of Tries : %v%v%v
Xray-core : %v%v%v
Xray-loglevel : %v%v%v
Shuffling : %v%v%v
Writer : %v%v%v
Total Threads : %v%v%v
-------------------------------------
`,
		utils.Colors.OKBLUE, C.Config.UserId, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.WsHeaderHost, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.WsHeaderPath, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.AddressPort, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.DoUploadTest, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.DoFrontingTest, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Download.MinDlSpeed, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Download.MaxDlTime, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Upload.MinUlSpeed, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Upload.MaxUlTime, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.TestUrl, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.FrontingTimeout, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Download.MaxDlLatency, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Upload.MaxUlLatency, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.NTries, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Vpn, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.LogLevel, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Shuffling, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Config.Writer, utils.Colors.ENDC,
		utils.Colors.OKBLUE, C.Worker.Threads, utils.Colors.ENDC,
	)
}

// PrintInformationForTUI prints configuration info formatted for TUI logs

func (C Configuration) CreateTestConfig(configPath string) Configuration {

	if configPath == "" {
		log.Fatalf("Configuration file are not loaded please use the --config or -c flag to use the configuration file.")
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Printf("%vError occurred during opening the configuration file.\n%v",
			utils.Colors.WARNING, utils.Colors.ENDC)
		log.Fatal(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
		}
	}(jsonFile)

	var jsonFileContent map[string]interface{}
	byteValue, _ := io.ReadAll(jsonFile)

	content := json.Unmarshal(byteValue, &jsonFileContent)
	if content != nil {
		return Configuration{}
	}

	// Safely extract string values with nil checks
	if val, ok := jsonFileContent["userId"]; ok && val != nil {
		C.Config.UserId = val.(string)
	} else {
		fmt.Printf("Error: 'userId' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["wsHeaderHost"]; ok && val != nil {
		C.Config.WsHeaderHost = val.(string)
	} else {
		fmt.Printf("Error: 'wsHeaderHost' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["addressPort"]; ok && val != nil {
		C.Config.AddressPort = val.(string)
	} else {
		fmt.Printf("Error: 'addressPort' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["sni"]; ok && val != nil {
		C.Config.Sni = val.(string)
	} else {
		fmt.Printf("Error: 'sni' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["wsHeaderPath"]; ok && val != nil {
		C.Config.WsHeaderPath = val.(string)
	} else {
		fmt.Printf("Error: 'wsHeaderPath' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["localPort"]; ok && val != nil {
		C.Config.LocalPort = int(val.(float64))
	} else {
		fmt.Printf("Error: 'localPort' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["frontingTimeout"]; ok && val != nil {
		C.Config.FrontingTimeout = val.(float64)
	} else {
		fmt.Printf("Error: 'frontingTimeout' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["testUrl"]; ok && val != nil {
		C.Config.TestUrl = val.(string)
	} else {
		C.Config.TestUrl = "http://google.com/generate_204" // Default if not in JSON
	}

	if val, ok := jsonFileContent["nTries"]; ok && val != nil {
		C.Config.NTries = int(val.(float64))
	} else {
		fmt.Printf("Error: 'nTries' field is missing or null in config file\n")
		os.Exit(1)
	}

	if val, ok := jsonFileContent["writer"]; ok && val != nil {
		C.Config.Writer = val.(string)
	} else {
		fmt.Printf("Error: 'writer' field is missing or null in config file\n")
		os.Exit(1)
	}

	// Only print configuration if not using TUI
	// TUI will handle its own configuration display
	// C.PrintInformation()
	return C
}

func CreateInterimResultsFile(interimResultsPath string, nTries int, writer string) error {
	emptyFile, err := os.Create(interimResultsPath)
	if err != nil {
		return fmt.Errorf("failed to create interim results file: %w", err)
	}

	defer func(emptyFile *os.File) {
		err := emptyFile.Close()
		if err != nil {

		}
	}(emptyFile)

	if strings.ToLower(writer) == "csv" {

		titles := []string{
			"ip",
			"avg_download_speed", "avg_upload_speed",
			"avg_download_latency", "avg_upload_latency",
			"avg_download_jitter", "avg_upload_jitter",
		}

		for i := 1; i <= nTries; i++ {
			titles = append(titles, fmt.Sprintf("ip_%d", i))
		}

		for i := 1; i <= nTries; i++ {
			titles = append(titles, fmt.Sprintf("download_speed_%d", i))
		}

		for i := 1; i <= nTries; i++ {
			titles = append(titles, fmt.Sprintf("upload_speed_%d", i))
		}

		for i := 1; i <= nTries; i++ {
			titles = append(titles, fmt.Sprintf("download_latency_%d", i))
		}

		for i := 1; i <= nTries; i++ {
			titles = append(titles, fmt.Sprintf("upload_latency_%d", i))
		}

		if _, err := fmt.Fprintln(emptyFile, strings.Join(titles, ",")); err != nil {
			return fmt.Errorf("failed to write titles to interim results file: %w", err)
		}

	}
	return nil
}

// PrintInformationForTUI sets configuration info in TUI
func (C Configuration) PrintInformationForTUI(tuiController interface{}) {
	// Type assertion to avoid import cycle
	if controller, ok := tuiController.(interface{ SetConfig(string, string) }); ok {
		controller.SetConfig("host", C.Config.WsHeaderHost)
		controller.SetConfig("path", C.Config.WsHeaderPath)
		controller.SetConfig("port", C.Config.AddressPort)
		if len(C.Config.UserId) > 8 {
			controller.SetConfig("userid", C.Config.UserId[:8]+"...")
		} else {
			controller.SetConfig("userid", C.Config.UserId)
		}
		controller.SetConfig("upload_test", fmt.Sprintf("%v", C.Config.DoUploadTest))
		controller.SetConfig("fronting_test", fmt.Sprintf("%v", C.Config.DoFrontingTest))
		controller.SetConfig("test_url", C.Config.TestUrl)
		controller.SetConfig("download_speed", fmt.Sprintf("%.0f kBps", C.Worker.Download.MinDlSpeed))
		controller.SetConfig("upload_speed", fmt.Sprintf("%.0f kBps", C.Worker.Upload.MinUlSpeed))
		controller.SetConfig("threads", fmt.Sprintf("%d", C.Worker.Threads))
		controller.SetConfig("tries", fmt.Sprintf("%d", C.Config.NTries))
		controller.SetConfig("xray_core", fmt.Sprintf("%v", C.Worker.Vpn))
		controller.SetConfig("writer", C.Config.Writer)
	}
}
