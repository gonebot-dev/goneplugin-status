package sysinfo

import (
	"fmt"
	"regexp"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"github.com/gonebot-dev/gonebot/api"
)

type diskInfo struct {
	Name        string  `json:"name"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type sysInfo struct {
	//Disk
	Disks []diskInfo `json:"disks"`
	//Mem
	MemAll         uint64  `json:"memAll"`
	MemUsed        uint64  `json:"memUsed"`
	MemUsedPercent float64 `json:"memUsedPercent"`
	//Boot time
	Days    int64 `json:"days"`
	Hours   int64 `json:"hours"`
	Minutes int64 `json:"minutes"`
	Seconds int64 `json:"seconds"`
	//CPU
	CpuUsedPercent float64 `json:"cpuUsedPercent"`
	CpuCores       int     `json:"cpuCores"`
	CpuInfo        string  `json:"cpuInfo"`
	CpuLoad1       float64 `json:"cpuLoad1"`
	CpuLoad5       float64 `json:"cpuLoad5"`
	CpuLoad15      float64 `json:"cpuLoad15"`
	//OS
	OS   string `json:"os"`
	Arch string `json:"arch"`
	//Gonebot
	SentTotal     int    `json:"sentTotal"`
	ReceivedTotal int    `json:"receivedTotal"`
	Backend       string `json:"backend"`
}

func GetSysInfo() (info sysInfo) {
	unit := uint64(1024 * 1024) // MB

	// Disks
	infos, _ := disk.Partitions(false)
	for _, inf := range infos {
		diskStat, err := disk.Usage(inf.Mountpoint)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		info.Disks = append(info.Disks, diskInfo{
			Name:        inf.Mountpoint,
			Total:       diskStat.Total / unit,
			Used:        diskStat.Used / unit,
			UsedPercent: diskStat.UsedPercent,
		})
	}

	// Mem
	v, _ := mem.VirtualMemory()
	info.MemAll = v.Total
	info.MemUsed = info.MemAll - v.Free
	info.MemUsedPercent = float64(info.MemUsed) / float64(info.MemAll) * 100.0
	info.MemAll /= unit
	info.MemUsed /= unit

	// CPU
	info.CpuCores, _ = cpu.Counts(true)
	cc, _ := cpu.Percent(time.Millisecond*200, false) //CPU usage in 200ms
	info.CpuUsedPercent = cc[0]
	dat, _ := cpu.Info()
	reg := regexp.MustCompile(`( @ ).*Hz`)
	info.CpuInfo = reg.ReplaceAllString(dat[0].ModelName, "")
	stat, _ := load.Avg()
	info.CpuLoad1 = stat.Load1
	info.CpuLoad5 = stat.Load5
	info.CpuLoad15 = stat.Load15

	// OS
	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH

	// 获取开机时间
	boottime, _ := host.BootTime()
	ntime := time.Now().Unix()
	btime := time.Unix(int64(boottime), 0).Unix()
	deltatime := ntime - btime

	info.Seconds = int64(deltatime)
	info.Minutes = info.Seconds / 60
	info.Seconds -= info.Minutes * 60
	info.Hours = info.Minutes / 60
	info.Minutes -= info.Hours * 60
	info.Days = info.Hours / 24
	info.Hours -= info.Days * 24

	info.SentTotal = api.GetResultCount()
	info.ReceivedTotal = api.GetIncomingCount()
	info.Backend = api.GetBackend()

	return
}
