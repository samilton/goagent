package checks

import (
	"log"
	"strconv"
	"syscall"
	"time"
)

type Disk struct {
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

func NewDiskCheck(name string, partition string, interval time.Duration, threshold float64) (d *Disk) {
	d = &Disk{
		Name:      name,
		Partition: partition,
		Interval:  interval,
		Threshold: threshold,
	}

	return
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
	return
}

func (e Disk) Report(disk *Disk) Message {
	t := time.Now().Unix()
	log.Printf("Reporting %s status", e.Name)
	log.Printf("%T: %t", disk, disk)
	return Message{t, e.Name, e.Status, strconv.FormatFloat(e.PercentUsed, 'f', -1, 32)}

}

func (e Disk) Run() {
	for {
		log.Printf("Performing %s Check", e.Name)
		DiskUsage(e)
		log.Printf("%T: %t", e, e)
		time.Sleep(e.Interval * time.Second)
	}
}
