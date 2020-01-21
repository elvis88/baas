package main

import (
	"time"

	"github.com/elvis88/baas/core/ws"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

func serverInfo() *ws.StatusServer {
	v, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := host.Info()

	ss := new(ws.StatusServer)
	ss.Uptime = n.Uptime
	ss.BootTime = n.BootTime
	ss.OS = n.OS
	ss.MemPercent = v.UsedPercent
	ss.CPUPercent = cc[0]
	ss.SwapPercent = s.UsedPercent
	ss.DiskPercent = d.UsedPercent
	return ss
}

func processInfo(pid int32) *ws.StatusProcess {
	ps := new(ws.StatusProcess)
	proc, err := process.NewProcess(pid)
	if err != nil {
		return ps
	}
	ps.PID = pid
	ps.BootTime, _ = proc.CreateTime()
	ps.MemPercent, _ = proc.MemoryPercent()
	ps.CPUPercent, _ = proc.CPUPercent()
	return ps
}
