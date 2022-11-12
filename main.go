package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "syscall"
)

func main(){
    switch os.Args[1] {
    case "run":
        run()
    case "child":
        child()
    default:
        panic("Use run or child to run the program `go run main.go run /bin/bash`")
    }
}

func run(){
    fmt.Printf("Running %v \n", os.Args[2:])
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...),...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysprocAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
        Unshareflags : syscall.CLONE_NEWNS,
    }
    must(cmd.Run())
}

func child(){
    fmt.Printf("Running %v \n", os.Args[3:]...)
    cg()
    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    must(syscall.Sethostname([]byte("container")))
    must(syscall.Chroot("/home/ubuntu/ubuntufs"))
    must(os.Chdir("/"))
    must(syscall.Mount("proc", "proc", "proc", 0, ""))
    must(syscall.Mount("thing", "mytemp", "tmpds", 0, ""))
    must(cmd.Run())
    must(syscall.Ubmount("proc", 0))
    must(syscall.Ubmount("thing", 0))
}

func cg(){
    cgroups := "/sys/fs/cgroup/"
    pids := filepath.Join(cgroups, "pids")
    os.Mkdir(filepath.Join(pids, "dore"), 0755)
    must(ioutil.WriteFile(filepath.Join(pids, "dore/pids.max"), []byte("20"), 0700))
    // Removes the new cgroup in place after the container exits
    must(ioutil.WriteFile(filepath.Join(pids, "dore/notify_on_release"), []byte("1"), 0700))
    must(ioutil.WriteFile(filepath.Join(pids, "dore/cgroup.rocs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error){
    if err != nil {
        paic(err)
    }
}
