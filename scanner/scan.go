package scanner

import (
	config "CFScanner/configuration"
	"CFScanner/logger"
	"CFScanner/speedtest"
	"CFScanner/utils"
	"CFScanner/vpn"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PrintLog prints a log message only if TUI is not active
func PrintLog(ld logger.ScannerManage) {
	if !IsTUIActive() {
		ld.Print()
	}
}

var results [][]string

var (
	downloadSpeed   float64
	downloadLatency float64
	uploadSpeed     float64
	uploadLatency   float64
)

type ScanResult struct {
	IP       string
	Download struct {
		Speed   []float64
		Latency []int
	}
	Upload struct {
		Speed   []float64
		Latency []int
	}
}

// Running Possible worker state.
var (
	Running bool
	MaxProc = runtime.NumCPU() * 2 // Max CPU + Thread * 2
)

// const WorkerCount = 48

func scanner(ip string, Config config.Configuration, Worker config.Worker) *ScanResult {

	result := &ScanResult{
		IP: ip,
	}

	var Upload = &Worker.Upload
	var Download = &Worker.Download

	var proxies map[string]string = nil
	var process vpn.ScanWorker

	if Worker.Vpn {
		// create config for desired ip
		xrayConfigPath := vpn.XRayConfig(ip, &Config)
		listen, port, _ := vpn.XRayReceiver(xrayConfigPath)

		// bind proxy
		proxies = vpn.ProxyBind(listen, port)

		// wait for port
		utils.WaitForPort(listen, port, time.Duration(5))

		var err error
		process = vpn.XRayInstance(xrayConfigPath)

		if err != nil {
			ld := logger.ScannerManage{
				IP:      "",
				Status:  logger.ErrorStatus,
				Message: "Could not start vpn service",
				Cause:   err.Error(),
			}
			PrintLog(ld)
			os.Exit(1)
			return nil
		}

		defer func() {
			// terminate process
			err = process.Instance.Close()
			if err != nil {
				ld := logger.ScannerManage{
					IP:      "",
					Status:  logger.ErrorStatus,
					Message: "Failed to stop xray-core instance",
					Cause:   err.Error(),
				}
				PrintLog(ld)
			}

			// Clean up config file
			os.Remove(xrayConfigPath)
		}()
	}

	for tryIdx := 0; tryIdx < Config.Config.NTries; tryIdx++ {
		// Step 1: Connectivity test - check if connection works at all
		connectivityOK, connErr := speedtest.ConnectivityTest(proxies, 5)
		if !connectivityOK || connErr != nil {
			ld := logger.ScannerManage{
				IP:      ip,
				Status:  logger.FailStatus,
				Message: "Connectivity check failed",
				Cause:   fmt.Sprintf("%v", connErr),
			}
			PrintLog(ld)
			return nil
		}

		ldConn := logger.ScannerManage{
			IP:      ip,
			Status:  logger.InfoStatus,
			Message: "Connectivity OK, proceeding to speed test",
		}
		PrintLog(ldConn)

		// Fronting test
		if Config.Config.DoFrontingTest {
			fronting := speedtest.FrontingTest(ip, proxies, time.Duration(Config.Config.FrontingTimeout))

			if !fronting {
				return nil
			}
		}

		// Check download speed
		if m, done := downloader(ip, Download, proxies, result); done {
			return m
		}

		// upload speed test
		if Config.Config.DoUploadTest {
			if m2, done2 := uploader(ip, Upload, proxies, result); done2 {
				return m2
			}
		}

		dlTimeLatency := math.Round(downloadLatency * 1000)
		upTimeLatency := math.Round(uploadLatency * 1000)

		ld := logger.ScannerManage{
			IP:     ip,
			Status: logger.InfoStatus,
			Message: fmt.Sprintf("Download: %7.4fmbps , Upload: %7.4fmbps , UP_Latency: %vms , DL_Latency: %vms",
				downloadSpeed, uploadSpeed, upTimeLatency, dlTimeLatency),
		}
		PrintLog(ld)
	}

	return result
}

func uploader(ip string, Upload *config.Upload, proxies map[string]string, result *ScanResult) (*ScanResult, bool) {
	var err error
	nBytes := Upload.MinUlSpeed * 1000 * Upload.MaxUlTime
	uploadSpeed, uploadLatency, err = speedtest.UploadSpeedTest(int(nBytes), proxies,
		time.Duration(Upload.MaxUlLatency))

	if err != nil {
		ld := logger.ScannerManage{
			IP:      ip,
			Status:  logger.FailStatus,
			Message: logger.UploadError,
			Cause:   err.Error(),
		}
		PrintLog(ld)
		return nil, true
	}

	if uploadLatency >= Upload.MaxUlLatency {
		ld := logger.ScannerManage{
			IP:      ip,
			Status:  logger.FailStatus,
			Message: logger.UploadLatency,
		}
		PrintLog(ld)
		return nil, true
	}

	uploadSpeedKbps := uploadSpeed / 8 * 1000

	if uploadSpeedKbps <= Upload.MinUlSpeed {
		ld := logger.ScannerManage{
			IP:     ip,
			Status: logger.FailStatus,
			Message: fmt.Sprintf("Upload too slow %f kBps < %f kBps",
				uploadSpeedKbps, Upload.MinUlSpeed),
		}
		PrintLog(ld)
		return nil, true
	}

	result.Upload.Speed = append(result.Upload.Speed, uploadSpeed)
	result.Upload.Latency = append(result.Upload.Latency, int(math.Round(uploadLatency*1000)))

	return nil, false
}

func downloader(ip string, Download *config.Download, proxies map[string]string, result *ScanResult) (*ScanResult, bool) {
	nBytes := Download.MinDlSpeed * 1000 * Download.MaxDlTime
	var err error

	downloadSpeed, downloadLatency, err = speedtest.DownloadSpeedTest(int(nBytes), proxies,
		time.Duration(Download.MaxDlLatency))

	if err != nil {
		ld := logger.ScannerManage{
			IP:      ip,
			Status:  logger.FailStatus,
			Message: logger.DownloadError,
			Cause:   err.Error(),
		}
		PrintLog(ld)
		return nil, true

	}

	if downloadLatency >= Download.MaxDlLatency {
		ld := logger.ScannerManage{
			IP:     ip,
			Status: logger.FailStatus,
			Message: fmt.Sprintf("High Download latency %.4f s > %.4f s",
				downloadLatency, Download.MaxDlLatency),
		}
		PrintLog(ld)

		return nil, true
	}
	downloadSpeedKBps := downloadSpeed / 8 * 1000

	if downloadSpeedKBps <= Download.MinDlSpeed {
		ld := logger.ScannerManage{
			IP:     ip,
			Status: logger.FailStatus,
			Message: fmt.Sprintf("Download too slow %.4f kBps < %.4f kBps",
				downloadSpeedKBps, Download.MinDlSpeed),
		}
		PrintLog(ld)
		return nil, true

	}
	result.Download.Speed = append(result.Download.Speed, downloadSpeed)
	result.Download.Latency = append(result.Download.Latency, int(math.Round(downloadLatency*1000)))

	return result, false
}

func scan(Config *config.Configuration, worker *config.Worker, ip string, workerID int) {
	// Check if TUI is paused
	if IsTUIPaused() {
		WaitTUIResume()
	}

	// Send worker update
	UpdateWorker(workerID, ip)

	res := scanner(ip, *Config, *worker)

	// Notify TUI of scan result
	UpdateTUI(workerID, ip, res != nil)

	if res == nil {
		return
	}

	// make downLatencyInt to float64
	downLatencyInt := res.Download.Latency
	downLatency := make([]float64, len(downLatencyInt))
	for i, v := range downLatencyInt {
		downLatency[i] = float64(v)
	}
	downMeanJitter := utils.MeanJitter(downLatency)

	// make uploadLatencyInt to float64
	uploadLatencyInt := res.Upload.Latency
	uploadLatency := make([]float64, len(uploadLatencyInt))
	for i, v := range uploadLatencyInt {
		uploadLatency[i] = float64(v)
	}
	upMeanJitter := -1.0

	if Config.Config.DoUploadTest {
		upMeanJitter = utils.MeanJitter(uploadLatency)
	}

	downSpeed := res.Download.Speed
	meanDownSpeed := utils.Mean(downSpeed)
	meanUploadSpeed := -1.0

	uploadSpeed := res.Upload.Speed
	if Config.Config.DoUploadTest {
		meanUploadSpeed = utils.Mean(uploadSpeed)
	}

	meanDownLatency := utils.Mean(downLatency)
	meanUploadLatency := -1.0
	if Config.Config.DoUploadTest {
		meanUploadLatency = utils.Mean(uploadLatency)
	}

	// change download latency to string type for using it with saveResults func
	var latencyDownloadString string
	for _, f := range downLatencyInt {
		latencyDownloadString = fmt.Sprintf("%d", f)
	}

	results = append(results, []string{latencyDownloadString, ip})

	var Writer Writer
	switch Config.Config.Writer {
	case "csv":
		Writer = CSV{
			res:                 res,
			IP:                  ip,
			DownloadMeanJitter:  downMeanJitter,
			UploadMeanJitter:    upMeanJitter,
			MeanDownloadSpeed:   meanDownSpeed,
			MeanDownloadLatency: meanDownLatency,
			MeanUploadSpeed:     meanUploadSpeed,
			MeanUploadLatency:   meanUploadLatency,
		}
	case "json":
		Writer = JSON{
			res:                 res,
			IP:                  ip,
			DownloadMeanJitter:  downMeanJitter,
			UploadMeanJitter:    upMeanJitter,
			MeanDownloadSpeed:   meanDownSpeed,
			MeanDownloadLatency: meanDownLatency,
			MeanUploadSpeed:     meanUploadSpeed,
			MeanUploadLatency:   meanUploadLatency,
		}
	default:
		cause := fmt.Errorf("invalid writer type: %s", Config.Config.Writer)
		ld := logger.ScannerManage{
			IP:      "",
			Status:  logger.ErrorStatus,
			Message: nil,
			Cause:   cause.Error(),
		}
		ld.Print()
		os.Exit(1)

	}

	Writer.Output()
	Writer.Write()

	// Save results & sort based on download latency
	err := saveResults(results, config.FinalResultsPathSorted, true)
	if err != nil {
		fmt.Println(err)
		return
	}

}

// StartScanWithWorkerIDs starts scanning with worker ID tracking for TUI
func StartScanWithWorkerIDs(C config.Configuration, Worker config.Worker, ipList []string, threadsCount int) {
	var wg sync.WaitGroup

	// limit the thread execution if it was higher than current cpu num * 2
	if threadsCount > MaxProc {
		fmt.Println("Max Thread limit setting thread to :", MaxProc)
		threadsCount = MaxProc
	}

	// Create batches
	n := len(ipList)
	batchSize := len(ipList) / threadsCount
	batches := make([][]string, threadsCount)

	for i := range batches {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == threadsCount-1 {
			end = n
		}
		batches[i] = ipList[start:end]
	}

	// Start workers
	Running = true
	for i := 0; i < threadsCount; i++ {
		wg.Add(1)
		go func(workerID int, batch []string) {
			defer wg.Done()
			for _, ip := range batch {
				select {
				case <-GetTUIContext():
					return // Stop if TUI signaled to quit
				default:
					scan(&C, &Worker, ip, workerID)
				}
			}
		}(i, batches[i])
	}

	wg.Wait()
}

// StartScan starts the scanning process with TUI support
func StartScan(C config.Configuration, Worker config.Worker, ipList []string, threadsCount int) {
	var wg sync.WaitGroup

	// limit the thread execution if it was higher than current cpu num * 2
	if threadsCount > MaxProc {
		fmt.Println("Max Thread limit setting thread to :", MaxProc)
		threadsCount = MaxProc
	}

	// Create batches
	n := len(ipList)
	batchSize := len(ipList) / threadsCount
	batches := make([][]string, threadsCount)

	for i := range batches {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == threadsCount-1 {
			end = n
		}
		batches[i] = ipList[start:end]
	}

	// Start workers
	Running = true
	stopCh := GetTUIContext()

	for i := 0; i < threadsCount; i++ {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			for _, ip := range batch {
				// Check if we should stop
				select {
				case <-stopCh:
					Running = false
					return
				default:
				}

				if !Running {
					return
				}

				// Check pause state
				for IsTUIPaused() {
					WaitTUIResume()
				}

				scan(&C, &Worker, ip, 0) // Use worker ID 0 for backward compatibility
			}
		}(batches[i])
	}

	wg.Wait()
}

// Start is the legacy function kept for backward compatibility
func Start(C config.Configuration, Worker config.Worker, ipList []string, threadsCount int) {
	var (
		wg         sync.WaitGroup
		pauseChan  = make(chan struct{})
		resumeChan = make(chan struct{})
		quitChan   = make(chan struct{})
	)

	// limit the thread execution if it was higher than current cpu num * 2
	if threadsCount > MaxProc {
		fmt.Println("Max Thread limit setting thread to :", MaxProc)
		threadsCount = MaxProc
	}

	// Create batches
	n := len(ipList)
	batchSize := len(ipList) / threadsCount
	batches := make([][]string, threadsCount)

	for i := range batches {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == threadsCount-1 {
			end = n
		}
		batches[i] = ipList[start:end]
	}

	// Start workers
	Running = true
	for i := 0; i < threadsCount; i++ {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			for _, ip := range batch {
				select {
				case <-pauseChan:
					// wait for resume signal
					<-resumeChan
				case <-quitChan:
					// quit the function
					return
				default:
					scan(&C, &Worker, ip, 0) // Use worker ID 0 for legacy function
				}
			}
		}(batches[i])
	}

	wg.Wait()

}

func saveResults(values [][]string, savePath string, sort bool) error {
	// clean the values and make sure the first element is integer
	for i := 0; i < len(values); i++ {
		ms, err := strconv.Atoi(strings.TrimSuffix(values[i][0], " ms"))
		if err != nil {
			return err
		}
		values[i][0] = strconv.Itoa(ms)
	}

	if sort {
		// sort the values based on response time using bubble sort
		for i := 0; i < len(values); i++ {
			for j := 0; j < len(values)-1; j++ {
				ms1, _ := strconv.Atoi(values[j][0])
				ms2, _ := strconv.Atoi(values[j+1][0])
				if ms1 > ms2 {
					values[j], values[j+1] = values[j+1], values[j]
				}
			}
		}
	}

	// write the values to file
	var lines []string
	for _, res := range values {
		lines = append(lines, strings.Join(res, " "))
	}
	data := []byte(strings.Join(lines, "\n") + "\n")
	err := os.WriteFile(savePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
