package worker

import (
	"fmt"
	"github.com/xfong/go2opencl/cl"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

func (mr *GpuMiner) buildOrLoadProgram() *cl.Program {

	var program *cl.Program

	binfilestuff := mr.platform.Name() // + "_" + mr.devices[0].Name()
	binfilename := strings.Replace(binfilestuff, " ", "_", -1)
	binfilepath := mr.openclPath + "/" + binfilename + ".objcache"
	binstat, staterr := os.Stat(binfilepath)
	if mr.rebuild || staterr != nil {
		fmt.Print("Create opencl program with source: " + mr.openclPath + ", Please wait...")
		buildok := false
		go func() { // Print
			for {
				time.Sleep(time.Second * 3)
				if buildok {
					break
				}
				fmt.Print(".")
			}
		}()
		emptyFuncTest := ""
		if mr.emptyFuncTest {
			emptyFuncTest = `_empty_test` // Quick compilation test of empty function
		}
		codeString := ` #include "x16rs_main` + emptyFuncTest + `.cl" `
		codeString += fmt.Sprintf("\n#define updateforbuild %d", rand.Uint64()) // Avoid compilation caching on some platforms
		program, _ = mr.context.CreateProgramWithSource([]string{codeString})
		bderr := program.BuildProgram(mr.devices, "-I "+mr.openclPath) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
		if bderr != nil {
			panic(bderr)
		}
		buildok = true // Build complete
		fmt.Println("\nBuild complete get binaries...")
		//fmt.Println("program.GetBinarySizes_2()")
		size := len(mr.devices)
		sizes, _ := program.GetBinarySizes_2(size)
		//fmt.Println(sizes)
		//fmt.Println("GetBinarySizes_2", sizes[0])
		//fmt.Println("program.GetBinaries_2()")
		bins, _ := program.GetBinaries_2(sizes)
		//fmt.Println("bins[0].size", len(bins[0]))
		f, e := os.OpenFile(binfilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if e != nil {
			panic(e)
		}
		//fmt.Println("f.Write(wbin) "+binfilepath, sizes[0])
		var berr error
		_, berr = f.Write(bins[0])
		if berr != nil {
			panic(berr)
		}
		berr = f.Close()
		if berr != nil {
			panic(berr)
		}

	} else {
		fmt.Printf("Load binary program file from \"%s\"\n", binfilepath)
		file, _ := os.OpenFile(binfilepath, os.O_RDONLY, 0777)
		bin := make([]byte, 0)
		//fmt.Println("file.Read(bin) size", binstat.Size())
		var berr error
		bin, berr = ioutil.ReadAll(file)
		if berr != nil {
			panic(berr)
		}
		if int64(len(bin)) != binstat.Size() {
			panic("int64(len(bin)) != binstat.Size()")
		}
		berr = file.Close()
		if berr != nil {
			panic(berr)
		}
		//fmt.Println(bin)
		// Only the same device on the same platform is supported
		bins := make([][]byte, len(mr.devices))
		sizes := make([]int, len(mr.devices))
		for k, _ := range mr.devices {
			bins[k] = bin
			sizes[k] = int(binstat.Size())
		}
		fmt.Println("Create program with binary...")
		program, berr = mr.context.CreateProgramWithBinary_2(mr.devices, sizes, bins)
		if berr != nil {
			panic(berr)
		}
		err := program.BuildProgram(mr.devices, "")
		if err != nil {
			panic(berr)
		}
		//fmt.Println("context.CreateProgramWithBinary")
	}
	fmt.Println("GPU miner program create complete successfully.")

	// return
	return program
}
