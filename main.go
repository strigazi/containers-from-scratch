package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
//	"golang.org/x/sys/unix"
)

// go run main.go run <cmd> <args>
// docker         run <cmd> <args>
func main() {
	fmt.Printf("Running as pid %v \n", os.Getpid())
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// https://golang.org/pkg/syscall/#SysProcAttr
	// https://godoc.org/golang.org/x/sys/unix
	// https://github.com/golang/sys/blob/master/unix/zerrors_linux_amd64.go#L13
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
//		AmbientCaps: []uintptr{unix.CAP_SYS_TIME,},
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cg()

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	// See Readme.md on how to download the root fs.
	must(syscall.Chroot("/root/ubuntufs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	// must(syscall.Mount("thing", "mytemp", "tmpfs", 0, ""))

	must(cmd.Run())

	must(syscall.Unmount("proc", 0))
	// must(syscall.Unmount("thing", 0))
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	// create dir for our demo
	// remove cgroup with rmdir https://lists.linuxcontainers.org/pipermail/lxc-users/2011-June/002106.html
	os.Mkdir(filepath.Join(pids, "demo"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "demo/pids.max"), []byte("10"), 0700))
	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "demo/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "demo/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
