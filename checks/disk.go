package checks

import (
	"log"
	"strconv"
	"syscall"
	"time"
)

type Disk struct {
	Queue       chan Message
	Name        string
	Interval    time.Duration
	Partition   string
	Threshold   float64
	Total       uint64
	Used        uint64
	Free        uint64
	PercentUsed float64
	Status      string
}

func GetStatus(percentUsed float64) (status string) {
	switch {
	case percentUsed < .5:
		return StatusClear
	case percentUsed > .5:
		return StatusWarn
	default:
		return StatusUnknown
	}
}

func DiskUsage(d Disk) (disk Disk) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(d.Partition, &fs)
	if err != nil {
		return
	}
	d.Total = fs.Blocks * uint64(fs.Bsize)
	d.Free = fs.Bfree * uint64(fs.Bsize)
	d.Used = d.Total - d.Free
	d.PercentUsed = float64(d.Used) / float64(d.Total)
	d.Status = GetStatus(d.PercentUsed)
	log.Printf("All: %d\tFree: %d\tUsed: %d\tPercent: %f", d.Total, d.Free, d.Used, d.PercentUsed)
	return d
}

func (e Disk) Run() {
	for {
		log.Printf("Performing %s Check", e.Name)
		t := time.Now().Unix()
		d := DiskUsage(e)
		log.Printf("%T: %t", d, d)
		m := Message{t, e.Name, d.Status, strconv.FormatFloat(d.PercentUsed, 'f', -1, 32)}
		e.Queue <- m
		time.Sleep(e.Interval * time.Second)
	}
}
