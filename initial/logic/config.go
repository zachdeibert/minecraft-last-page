package logic

// Change to configure
const (
	minStringLen = 6
	maxStringLen = 20
	suffixLen    = 2  // determines number of GPU threads
	baseLen      = 3  // determines how often data is copied between the host memory and device memory
	outBuf       = 20 // maximum number of matches per batch
)

// Probably don't change
const (
	alphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()"
	outputFile = "../candidates.txt"
)
