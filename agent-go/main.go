package main

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
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
    HostID string `json:"host_id"`
    Token  string `json:"token"`
    Server string `json:"server"`
}

func main() {
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
        }

        jsonData, _ := json.Marshal(data)
        http.Post(config.Server, "application/json", bytes.NewBuffer(jsonData))
        time.Sleep(60 * time.Second)
    }
}

func loadConfig() Config {
    file, _ := os.Open("config.json")
    defer file.Close()
    bytes, _ := ioutil.ReadAll(file)
    var config Config
    json.Unmarshal(bytes, &config)
    return config
}
