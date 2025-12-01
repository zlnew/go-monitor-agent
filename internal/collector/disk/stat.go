package disk

import "syscall"

func (c *Collector) readStat() error {
	c.totalBytes = 0
	c.freeBytes = 0

	var stat syscall.Statfs_t

	if err := syscall.Statfs("/", &stat); err != nil {
		return err
	}

	bsize := uint64(stat.Bsize)
	c.totalBytes = float64(stat.Blocks * bsize)
	c.freeBytes = float64(stat.Bavail * bsize)

	return nil
}
