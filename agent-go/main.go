package main

import (
    "bytes"
    "encoding/json"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/shirou/gopsutil/v4/cpu"
    "github.com/shirou/gopsutil/v4/disk"
    "github.com/shirou/gopsutil/v4/mem"
    "github.com/shirou/gopsutil/v4/net"
    "github.com/shirou/gopsutil/v4/host"
)

type Config struct {
    HostID   string `json:"host_id"`
    Token    string `json:"token"`
    Server   string `json:"server"`
    Interval int    `json:"interval"` // 上报间隔（秒）
}

var logger *log.Logger

func initLogger() {
    logFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    multi := io.MultiWriter(os.Stdout, logFile)
    logger = log.New(multi, "", log.LstdFlags)
}

func main() {
    initLogger()
    config := loadConfig()
    for {
        cpuPercents, _ := cpu.Percent(0, false)
        memInfo, _ := mem.VirtualMemory()
        diskInfo, _ := disk.Usage("/")
        netIO, _ := net.IOCounters(false)
        uptime, _ := host.Uptime()

        data := map[string]interface{}{
            "host_id": config.HostID,
            "token":   config.Token,
            "cpu":     cpuPercents[0],
            "mem":     memInfo.UsedPercent,
            "disk":    diskInfo.UsedPercent,
            "net_in":  netIO[0].BytesRecv,
            "net_out": netIO[0].BytesSent,
            "uptime":  uptime,
            "timestamp": time.Now().Unix(),
        }

        jsonData, _ := json.Marshal(data)
        resp, err := http.Post(config.Server, "application/json", bytes.NewBuffer(jsonData))

        if err != nil {
            logger.Printf("[ERROR] Failed to POST to %s: %v\n", config.Server, err)
        } else {
            defer resp.Body.Close()
            logger.Printf("[INFO] Sent data to %s - Status: %s\n", config.Server, resp.Status)
        }

        time.Sleep(time.Duration(config.Interval) * time.Second)
    }
}

func loadConfig() Config {
    file, _ := os.Open("config.json")
    defer file.Close()
    bytes, _ := ioutil.ReadAll(file)
    var config Config
    json.Unmarshal(bytes, &config)
    if config.Interval == 0 {
        config.Interval = 60
    }
    return config
}
