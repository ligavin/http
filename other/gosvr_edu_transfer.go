package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "regexp"
    "flag"
    "strings"
    "time"
    "os/exec"
)

type Size interface {
    Size() int64
}

var ARG_ADDR = flag.String("bind", ":8080", "http bind addr")
var ARG_PATH_PREFIX = flag.String("path", "^/data(/[\\w.0-9_-]+)+$", "file save path prefix, regexp")
var ARG_ALLOWED_IPS = flag.String("allow", "127.0.0.1,10.123.14.55,10.123.14.44,100.116.34.181,10.213.135.192,10.213.137.136", "only these IPs allowed")
var ARG_ALLOWED_CMDS = flag.String("allow_cmd", "cp", "only these cmd allowed")
var ARG_ALLOWED_CMD_DIR = flag.String("allow_cmd_dir", "/data/home/cmd", "only shell file in this dir are allowd")
var PATH_PREFIX_REGX *regexp.Regexp
var ALLOWED_IPS []string
var ALLOWED_CMDS []string

var EDU_ENV_INIT_SCRIPT_DIR string

func flowLog(remote string, path string, size int64, result string, msg string, err error) {
    fmt.Println(time.Now(), ", remote", remote, ", filepath", path, ", size", size, ", result", result, ", msg", msg, ", err", err)
}

func checkIP(w http.ResponseWriter, r *http.Request) (remote string, allowed bool) {
    remote = r.Header.Get("X-Real-IP")
    if remote == "" {
        remote = r.RemoteAddr
    }
    allowed = false
    for i := 0; i < len(ALLOWED_IPS); i++ {
        if ALLOWED_IPS[i] != "" && strings.Contains(remote, ALLOWED_IPS[i]) {
            allowed = true
            break
        }
    }
    if !allowed {
        msg := "remote is not allowed " + remote
        io.WriteString(w, msg)
        flowLog(remote, "-", 0, "rejected", msg, nil)
        return
    }
    return
}

func transfer(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain")
    remote, ok := checkIP(w, r)
    if !ok {
        return
    }

    file, _, err := r.FormFile("file")
    if err != nil {
        msg := "file is empty?"
        io.WriteString(w, msg)
        flowLog(remote, "-", 0, "error", msg, err)
        return
    }
    defer file.Close()

    var fileSize int64 = 0;
    fileStatHandle, ok := file.(Size)
    if !ok {
        fmt.Println("get file stat failed")
        fileSize = 0x7FFFFFFF
    } else {
        fileSize = fileStatHandle.Size()
    }

    path := r.FormValue("path")
    isMatched := PATH_PREFIX_REGX.MatchString(path)
    if !isMatched || strings.Contains(path, "..") {
        msg := "path error, path=" + path + ", reg exp=" + *ARG_PATH_PREFIX
        io.WriteString(w, msg)
        flowLog(remote, path, fileSize, "error", msg, nil)
        return
    }

    fW, err := os.Create(path)
    if err != nil {
        msg := "open file failed, path=" + path
        io.WriteString(w, msg)
        flowLog(remote, path, fileSize, "error", msg, err)
        return
    }
    defer fW.Close()
    _, err = io.Copy(fW, file)
    if err != nil {
        msg := "write file failed, path=" + path
        io.WriteString(w, msg)
        flowLog(remote, path, fileSize, "error", msg, err)
        return
    }
    io.WriteString(w, "done\n")
    flowLog(remote, path, fileSize, "done", "", nil)
}

func substr(s string, pos, length int) string {
    runes := []rune(s)
    l := pos + length
    if l > len(runes) {
        l = len(runes)
    }
    return string(runes[pos:l])
}

func getDirectory(dirctory string) string {
    return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

//收到具体的命令执行
func command(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain")
    remote, ok := checkIP(w, r)
    if !ok {
        return
    }

    cmd_str := r.FormValue("cmd")

    if cmd_str == "" {
        msg := "cmd is empty";
        io.WriteString(w, msg);
        flowLog(remote, "-", 0, "command not valid", msg, nil)
        return
    }

    allowed := false

    //执行操作
    if cmd_str != "" && strings.HasPrefix(cmd_str, "sh") {

        cmd_info := strings.Split(cmd_str, " ")

        if len(cmd_info) != 2 || cmd_info[1] == "" {

            msg := "cmd bad format,should be [sh cmd_file]"
            io.WriteString(w, msg)
            flowLog(remote, "-", 0, "rejected", msg, nil)
            return
        }

        cmd_directory := getDirectory(cmd_info[1])

        if cmd_directory != *ARG_ALLOWED_CMD_DIR {
            msg := "cmd not allowd,sh command shoulde be only in " + *ARG_ALLOWED_CMD_DIR
            io.WriteString(w, msg)
            flowLog(remote, "-", 0, "rejected", msg, nil)
            return
        }

        allowed = true

    } else {

        for _, cmd := range ALLOWED_CMDS {
            if cmd != "" && strings.HasPrefix(cmd_str, cmd) {
                allowed = true
            }
        }

    }

    if allowed == false {
        msg := "cmd is not allowd,only " + *ARG_ALLOWED_CMDS + " allowd"
        io.WriteString(w, msg)
        flowLog(remote, "-", 0, "command not valid", msg, nil)
        return
    }

    execCmdAndOutput2Http(w, r, cmd_str, remote)
}

// 执行EDU测试环境初始化、销毁
// action, app_type, app_name, app_port
func eduEnvInit(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain")
    remote, ok := checkIP(w, r)
    if !ok {
        return
    }

    var script string
    switch r.FormValue("action") {
    case "create":
        script = "create_"
    case "destroy":
        script = "destroy_"
    default:
        w.Write([]byte(remote + " Unknown ACTION" + r.FormValue("action")))
        return
    }
    switch r.FormValue("app_type") {
    case "JungleCGI":
        script += "cgi_testenv.sh"
    case "JungleServer":
        script += "svr_testenv.sh"
    default:
        w.Write([]byte(remote + " Unknown APP TYPE" + r.FormValue("app_type")))
        return
    }
    var paramAppName = r.FormValue("app_name")
    var paramAppPort = r.FormValue("app_port")
    var cmdStr = EDU_ENV_INIT_SCRIPT_DIR + "/" + script + " " + paramAppName + " " + paramAppPort
    execCmdAndOutput2Http(w, r, cmdStr, remote)
}

func execCmdAndOutput2Http(w http.ResponseWriter, r *http.Request, cmd_str, remote string) {
    fmt.Println("exec cmd", cmd_str)
    cmd := exec.Command("/bin/sh", "-c", cmd_str);

    stdout, err := cmd.StdoutPipe();

    if err != nil {
        msg := err.Error()
        io.WriteString(w, msg);
        flowLog(remote, "-", 0, cmd_str + " StdoutPipe failed", msg, nil)
        return
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        msg := err.Error()
        io.WriteString(w, msg);
        flowLog(remote, "-", 0, cmd_str + " StderrPipe failed", msg, nil)
        return
    }

    if err := cmd.Start(); err != nil {
        msg := err.Error()
        io.WriteString(w, msg)
        flowLog(remote, "-", 0, cmd_str + " Start failed", msg, nil)
        return
    }
    defer cmd.Process.Wait()

    bytesErr, err := ioutil.ReadAll(stderr)
    if err != nil {
        msg := err.Error()
        io.WriteString(w, msg)
        flowLog(remote, "-", 0, cmd_str + " ReadAll1 failed", msg, nil)
        return
    }

    if len(bytesErr) != 0 {
        msg := string(bytesErr)
        io.WriteString(w, msg)
        flowLog(remote, "-", 0, cmd_str + " ReadAll2 failed", msg, nil)
        return
    }

    bytes, err := ioutil.ReadAll(stdout)
    if err != nil {
        msg := err.Error()
        io.WriteString(w, msg)
        flowLog(remote, "-", 0, cmd_str + " ReadAll3 failed", msg, nil)
        return
    }

    io.WriteString(w, string(bytes))
}

func main() {

    flag.Parse()
    PATH_PREFIX_REGX = regexp.MustCompile(*ARG_PATH_PREFIX)
    ALLOWED_IPS = strings.Split(*ARG_ALLOWED_IPS, ",")
    ALLOWED_CMDS = strings.Split(*ARG_ALLOWED_CMDS, ",")
    var err error
    EDU_ENV_INIT_SCRIPT_DIR, err = os.Getwd()
    if err != nil {
        fmt.Println("server start failed ", err)
        return
    }

    fmt.Println("http bind to ", (*ARG_ADDR),
        ", file save path prefix ", (*ARG_PATH_PREFIX),
        ", allowed ip ", ALLOWED_IPS,
        ", EDU_ENV_INIT_SCRIPT_DIR=", EDU_ENV_INIT_SCRIPT_DIR);
    http.HandleFunc("/go/transfer", transfer)
    http.HandleFunc("/go/cmd", command)
    http.HandleFunc("/go/edu_env", eduEnvInit)
    err = http.ListenAndServe(*ARG_ADDR, nil)
    if err != nil {
        fmt.Println("server start failed", err)
        return
    }
    fmt.Println("server start success")
}
