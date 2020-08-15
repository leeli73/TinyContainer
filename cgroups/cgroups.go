package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var groupPath = "/sys/fs/cgroup/"

type CGroups struct {
	Name        string
	MemoryLimit string
	CpuSet      string
	CpuShare    string
}

func (t *CGroups) Create() string {
	res := ""
	groupPath = path.Join(groupPath, t.Name)
	if _, err := os.Stat(groupPath); os.IsNotExist(err) {
		if err := os.Mkdir(groupPath, 0755); err != nil {
			fmt.Printf("Create cgroups error:%v\n", err)
			fmt.Println(`You can execute "mount -o remount,rw '/sys/fs/cgroup'" by root or sudo to fix this problem.`)
			os.Exit(1)
		}
	}
	if t.MemoryLimit != "" {
		t.writeMemLimit()
		res = res + "Max Memory Limit:" + t.MemoryLimit + "\n"
	}
	if t.CpuSet != "" {
		t.writeCPUSet()
		res = res + "Max CPU Limit:" + t.CpuSet + "\n"
	}
	if t.CpuShare != "" {
		t.writeCPUShare()
		res = res + "Max CPU Share Limit:" + t.CpuShare + "\n"
	}
	return res
}
func (t *CGroups) Remove() {
	if t.Name != "" {
		os.Remove(groupPath)
	}
}
func (t *CGroups) writeMemLimit() {
	err := ioutil.WriteFile(path.Join(groupPath, "memory.limit_in_bytes"), []byte(t.MemoryLimit), 0755)
	if err != nil {
		fmt.Printf("Write to cgroups memory error:%v\n", err)
		os.Exit(1)
	}
}
func (t *CGroups) writeCPUSet() {
	err := ioutil.WriteFile(path.Join(groupPath, "cpuset.cpus"), []byte(t.CpuSet), 0755)
	if err != nil {
		fmt.Printf("Write to cgroups cpuset error:%v\n", err)
		os.Exit(1)
	}
}
func (t *CGroups) writeCPUShare() {
	err := ioutil.WriteFile(path.Join(groupPath, "cpu.shares"), []byte(t.CpuShare), 0755)
	if err != nil {
		fmt.Printf("Write to cgroups cpu share error:%v\n", err)
		os.Exit(1)
	}
}
func (t *CGroups) writeTask() {
	err := ioutil.WriteFile(path.Join(groupPath, "tasks"), []byte(t.MemoryLimit), 0755)
	if err != nil {
		fmt.Printf("Write to cgroups task error:%v\n", err)
		os.Exit(1)
	}
}
