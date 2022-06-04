package worker

import "github.com/xfong/go2opencl/cl"

// GPU mining planning
type GpuMiner struct {
	platform *cl.Platform
	context  *cl.Context
	program  *cl.Program
	devices  []*cl.Device // All equipment

	deviceworkers []*GpuMinerDeviceWorkerContext

	// config
	openclPath        string
	rebuild           bool   // Force recompile
	platName          string // Selected platforms
	groupNum          int    // Number of simultaneous execution groups
	groupSize         int    // Group size
	itemLoop          int    // Number of single execution cycles
	emptyFuncTest     bool   // Empty function compilation test
	useOneDeviceBuild bool   // Compile using a single device

}

// initialization
func NewGpuMiner(
	openclPath string,
	platName string,
	groupSize int, // Group width
	groupNum int, // Number of simultaneous execution groups: 1 ~ 64
	itemLoop int, // Suggestion 20 ~ 100
	useOneDeviceBuild bool, // Use a device to compile
	rebuild bool,
	emptyFuncTest bool,
) *GpuMiner {

	miner := &GpuMiner{
		openclPath:        openclPath,
		platName:          platName,
		rebuild:           rebuild,
		emptyFuncTest:     emptyFuncTest,
		useOneDeviceBuild: useOneDeviceBuild,
		groupSize:         groupSize,
		groupNum:          groupNum,
		itemLoop:          itemLoop,
	}

	// Created successfully
	return miner
}
