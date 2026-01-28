package scanner

import (
	config "CFScanner/configuration"
	"CFScanner/tui"
	"sync"
)

// Global TUI controller and state
var tuiCtrl *tui.Controller
var tuiMu sync.Mutex

// SetTUIController sets the global TUI controller
func SetTUIController(ctrl *tui.Controller) {
	tuiMu.Lock()
	defer tuiMu.Unlock()
	tuiCtrl = ctrl
}

// UpdateTUI sends an update to the TUI
func UpdateTUI(workerID int, ip string, success bool) {
	tuiMu.Lock()
	ctrl := tuiCtrl
	tuiMu.Unlock()

	if ctrl != nil {
		ctrl.SendUpdate(workerID, ip, success)
	}
}

// UpdateWorker sends a worker status update to the TUI
func UpdateWorker(workerID int, ip string) {
	tuiMu.Lock()
	ctrl := tuiCtrl
	tuiMu.Unlock()

	if ctrl != nil {
		ctrl.SendWorkerUpdate(workerID, ip)
	}
}

// IsTUIPaused returns whether the TUI is paused
func IsTUIPaused() bool {
	tuiMu.Lock()
	ctrl := tuiCtrl
	tuiMu.Unlock()

	if ctrl != nil {
		return ctrl.IsPaused()
	}
	return false
}

// WaitTUIResume waits for the TUI to resume
func WaitTUIResume() {
	tuiMu.Lock()
	ctrl := tuiCtrl
	tuiMu.Unlock()

	if ctrl != nil {
		ctrl.WaitForResume()
	}
}

// GetTUIContext returns the TUI context for cancellation
func GetTUIContext() <-chan struct{} {
	tuiMu.Lock()
	ctrl := tuiCtrl
	tuiMu.Unlock()

	if ctrl != nil {
		return ctrl.GetStopChannel()
	}
	return nil
}

// IsTUIActive returns whether TUI is active
func IsTUIActive() bool {
	tuiMu.Lock()
	ctrl := tuiCtrl
	tuiMu.Unlock()

	return ctrl != nil
}

// StartWithTUI starts the scanner with TUI interface
func StartWithTUI(C config.Configuration, Worker config.Worker, ipList []string, threadsCount int, tuiController *tui.Controller) {
	// Set global TUI controller
	SetTUIController(tuiController)

	// Start scanning in background
	scanDone := make(chan struct{})
	go func() {
		defer close(scanDone)
		StartScanWithWorkerIDs(C, Worker, ipList, threadsCount)
	}()

	// Start TUI in main thread (blocks until quit)
	_ = tuiController.Start()

	// TUI has exited, scanning should also stop
	// The scanner will check for stop signal in its loop
}
