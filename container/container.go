package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/leeli73/TinyContainer/cgroups"
	"golang.org/x/crypto/ssh/terminal"
)

type Container struct {
	Name          string
	FsDir         string
	Command       string
	EnableTTY     bool
	EnableNetwork bool
	MemoryLimit   string
	CpuSet        string
	CpuShare      string
}

func NewContainer(name string, fsDir string, command string, enableTTY bool, enableNetWork bool) Container {
	res := Container{
		Name:          name,
		FsDir:         fsDir,
		Command:       command,
		EnableTTY:     enableTTY,
		EnableNetwork: enableNetWork,
	}
	return res
}

func (c *Container) Child(parameter []string) {
	if c.EnableNetwork {
		fmt.Println(c.Name, "startup with network success.")
	} else {
		fmt.Println(c.Name, "startup success.")
	}
	syscall.Sethostname([]byte(c.Name))
	tmpMountPoint := "/tmp/" + c.Name + "/"
	syscall.Mount("/sys", filepath.Join(tmpMountPoint, "sys"), "sysfs", 0, "")
	syscall.Mount("udev", filepath.Join(tmpMountPoint, "dev"), "devtmpfs", 0, "")
	syscall.Mount("devpts", filepath.Join(tmpMountPoint, "dev/pts"), "devpts", 0, "")
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	syscall.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin:/")
	if c.FsDir == "/" {
		err := syscall.Chroot(c.FsDir)
		if err != nil {
			fmt.Printf("Chroot error:%v\n", err)
			os.Exit(1)
		}
	} else {
		err := pivotRoot(c.FsDir)
		if err != nil {
			fmt.Printf("pivotRoot error:%v\n", err)
			os.Exit(1)
		}
	}
	err := syscall.Chdir("/")
	if err != nil {
		fmt.Printf("Chdir error:%v\n", err)
		os.Exit(1)
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		fmt.Printf("Mount proc error:%v\n", err)
		os.Exit(1)
	}
	parameter[0] = strings.Replace(parameter[0], "'", "", -1)
	parameter = strings.Split(parameter[0], " ")
	err = syscall.Exec(parameter[0], parameter[1:], syscall.Environ())
	if err != nil {
		fmt.Printf("some error happend when Child:%v\n", err)
		os.Exit(1)
	}
}
func (c *Container) Run(parameter []string) {
	cgContral := cgroups.CGroups{
		Name:        c.Name,
		MemoryLimit: c.MemoryLimit,
		CpuSet:      c.CpuSet,
		CpuShare:    c.CpuShare,
	}
	res := cgContral.Create()
	defer cgContral.Remove()
	if res != "" {
		fmt.Printf("Create CGroups Suuess.\n%s", res)
	}
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, parameter[:]...)...)
	cmd.Dir = c.FsDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if c.EnableNetwork {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
		}
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		}
	}
	err := cmd.Run()
	if err != nil {
		fmt.Printf("some error happend when Run:%v\n", err)
		os.Exit(1)
	}
}
func (c *Container) RunWithPTY(parameter []string) {
	cgContral := cgroups.CGroups{
		Name:        c.Name,
		MemoryLimit: c.MemoryLimit,
		CpuSet:      c.CpuSet,
		CpuShare:    c.CpuShare,
	}
	res := cgContral.Create()
	defer cgContral.Remove()
	if res != "" {
		fmt.Printf("Create CGroups Suuess.\n%s", res)
	}
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, parameter[:]...)...)
	cmd.Dir = c.FsDir
	if c.EnableNetwork {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
		}
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		}
	}
	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Printf("some error happend when PTY TerminalL:%v\n", err)
	}
	defer ptmx.Close()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				fmt.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("error when make terminal row:%v\n", err)
		os.Exit(-1)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)
}
func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	if _, err := os.Stat(pivotDir); os.IsNotExist(err) {
		if err := os.Mkdir(pivotDir, 0777); err != nil {
			return err
		}
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	return os.Remove(pivotDir)
}
