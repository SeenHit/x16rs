package worker

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/mint/difficulty"
	"github.com/hacash/x16rs"
	"sync"
	"sync/atomic"
)

// Turn off force statistics
func (g *GpuMiner) CloseUploadHashrate() {
}

// Start mining
func (g *GpuMiner) GetSuperveneWide() int {
	return len(g.devices)
}

// Start mining
func (g *GpuMiner) DoMining(blockHeight uint64, reporthashrate bool, stopmark *byte, tarhashvalue []byte, blockheadmeta [][]byte) (bool, int, []byte, []byte) {

	deviceNum := len(g.devices)
	//fmt.Print(overallstep)

	var successed bool = false
	var successMark uint32 = 0
	var successStuffIdx int = 0
	var successNonce []byte = nil
	var successHash []byte = nil

	// Sync wait
	var syncWait = sync.WaitGroup{}
	syncWait.Add(deviceNum)

	// Equipment execution
	for i := 0; i < deviceNum; i++ {
		go func(did int) {
			defer syncWait.Done()
			//fmt.Println("mr.deviceworkers[i]", did, len(g.deviceworkers), g.deviceworkers)
			//devideCtx := g.deviceworkers[did]
			stuffbts := blockheadmeta[did]
			// implement
			x16rsrepeat := uint32(x16rs.HashRepeatForBlockHeight(blockHeight))
			var basenoncestart uint64 = 1
		RUNMINING:
			// Initialize execution environment
			//devideCtx := g.createWorkContext(did)
			devideCtx := g.deviceworkers[did]
			devideCtx.ReInit(stuffbts, tarhashvalue)
			//fmt.Println("DO RUNMINING...")
			//ttstart := time.Now()
			groupsize := g.devices[did].MaxWorkGroupSize()
			if g.groupSize > 0 {
				groupsize = int(g.groupSize)
			}
			globalwide := groupsize * g.groupNum
			overstep := globalwide * g.itemLoop // 单次挖矿 nonce 范围
			//fmt.Println(overstep, groupsize)
			success, nonce, endhash := g.doGroupWork(devideCtx, globalwide, groupsize, x16rsrepeat, uint32(basenoncestart))
			//devideCtx.Release() // 释放
			fmt.Print("_")
			//fmt.Println("END RUNMINING:", time.Now().Unix(), time.Now().Unix() - ttstart.Unix(), success, hex.EncodeToString(nonce), hex.EncodeToString(endhash) )
			if success && atomic.CompareAndSwapUint32(&successMark, 0, 1) {
				successed = true
				*stopmark = 1
				successStuffIdx = did
				successNonce = nonce
				successHash = endhash
				// Check whether it is really successful
				blk, _, _ := blocks.ParseExcludeTransactions(stuffbts, 0)
				blk.SetNonceByte(nonce)
				nblkhx := blk.HashFresh()
				if difficulty.CheckHashDifficultySatisfy(nblkhx, tarhashvalue) == false || bytes.Compare(nblkhx, endhash) != 0 {
					fmt.Println("挖矿失败！！！！！！！！！！！！！！！！")
					fmt.Println(nblkhx.ToHex(), hex.EncodeToString(endhash))
					fmt.Println(hex.EncodeToString(stuffbts))
				}

				return // Successful excavation, end
			}
			if *stopmark == 1 {
				//fmt.Println("ok.")
				return // One block is scarce, end
			}
			// Continue to collect funds
			basenoncestart += uint64(overstep)
			if basenoncestart > uint64(4294967295) {
				//if basenoncestart > uint64(529490) {
				return // Completion of this round of mining
			}
			//time.Sleep(time.Second * 5)
			goto RUNMINING
		}(i)
	}

	//fmt.Println("syncWait.Wait()")
	// wait for
	syncWait.Wait()

	//fmt.Println("syncWait.Wait() ok  返回")

	// return
	return successed, successStuffIdx, successNonce, successHash

}
