package sysinfo

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type sysInfo struct {
	//Disk
	DiskAll         uint64
	DiskFree        uint64
	DiskUsed        uint64
	DiskUsedPercent float64
	//Mem
	MemAll         uint64
	MemFree        uint64
	MemUsed        uint64
	MemUsedPercent float64
	//Boot time
	Days    int64
	Hours   int64
	Minutes int64
	Seconds int64
	//CPU
	CpuUsedPercent float64
	CpuCores       int
	CpuInfo        string
	//OS
	OS   string
	Arch string
	//Gonebot
	SentTotal     int64
	ReceivedTotal int64
}

func GetSysInfo() (info sysInfo) {
	//Mem
	v, _ := mem.VirtualMemory()
	info.MemAll = v.Total
	info.MemFree = v.Free
	info.MemUsed = info.MemAll - info.MemFree
	info.MemUsedPercent = float64(info.MemUsed) / float64(info.MemAll) * 100.0
	unit := uint64(1024 * 1024) // MB
	info.MemAll /= unit
	info.MemUsed /= unit
	info.MemFree /= unit
	//CPU
	info.CpuCores = runtime.GOMAXPROCS(0)
	cc, _ := cpu.Percent(time.Millisecond*200, false) //CPU usage in 200ms
	info.CpuUsedPercent = cc[0]
	//OS
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
	return
}
