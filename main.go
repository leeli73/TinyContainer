package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	tc "github.com/leeli73/TinyContainer/container"
)

var Version = "beta/0.0.1"
var (
	name          string
	fsDir         string
	memoryLimit   string
	cpuShare      string
	cpuSet        string
	enableTTY     bool
	enableNetwork bool
)

func main() {
	getValue(os.Args)
	if len(os.Args) <= 1 {
		fmt.Println("Can't find command.")
		return
	}
	switch os.Args[1] {
	case "run":
		parm := []string{}
		if name == "" {
			name = randomString(16)
			parm = append(parm, "name")
			parm = append(parm, name)
		}
		parm = append(parm, os.Args[2:]...)
		container := tc.Container{
			Name:          name,
			FsDir:         fsDir,
			Command:       os.Args[len(os.Args)-1],
			EnableTTY:     enableTTY,
			EnableNetwork: enableNetwork,
			MemoryLimit:   memoryLimit,
			CpuSet:        cpuSet,
			CpuShare:      cpuShare,
		}
		container.Run(parm)
	case "child":
		if os.Args[0] != "/proc/self/exe" {
			fmt.Println("This command must run by TinyContainer.")
			return
		}
		container := tc.Container{
			Name:          name,
			FsDir:         fsDir,
			Command:       os.Args[len(os.Args)-1],
			EnableTTY:     enableTTY,
			EnableNetwork: enableNetwork,
			MemoryLimit:   memoryLimit,
			CpuSet:        cpuSet,
			CpuShare:      cpuShare,
		}
		container.Child(os.Args[len(os.Args)-1:])
	default:
		fmt.Println("can't find the main command!")
	}
}
func getValue(parm []string) {
	for i := 1; i < len(parm); i++ {
		switch parm[i] {
		case "cpuset":
			cpuSet = parm[i+1]
			i = i + 1
		case "cpushare":
			cpuShare = parm[i+1]
			i = i + 1
		case "fs":
			fsDir = parm[i+1]
			i = i + 1
		case "memlimit":
			memoryLimit = parm[i+1]
			i = i + 1
		case "name":
			name = parm[i+1]
			i = i + 1
		case "tty":
			enableTTY = true
		case "net":
			enableNetwork = true
		case "-help":
			usage()
			os.Exit(0)
		case "-h":
			usage()
			os.Exit(0)
		}
	}
}
func usage() {
	fmt.Fprintf(os.Stderr, `TinyContainer version: `+Version+`
Usage Simple Example: tcontainer [run] [-name ContainerName] [-fs ContainerFSroot]

Options:
  run 
	Run Container
  cpuset string
	CPU core that container can use.Unlimited by default.
  cpushare string
	CPU Share that container can use.Unlimited by default.
  fs string
	The path of the Fire System used by the container.Default to use you fsroot. (default "/")
  memlimit string
	Max Memory that container can use.Unlimited by default.
  name string
	The name of container.Default to random string.
  net bool
	Whether to init network.Share your machine.Default to false.
  tty bool
	Whether to start tty.Default to true.
  -help or -h
	Show help information.
`)
}
func randomString(len int) string {
	var res string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		res += string(str[randomInt.Int64()])
	}
	return res
}
