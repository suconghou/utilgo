package utilgo

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//Bar progress
func Bar(loaded int, width int) string {
	if loaded > 100 {
		loaded = 100
	}
	loaded = int(float32(loaded) / 100 * float32(width))
	remain := width - loaded
	if remain < 0 {
		remain = 0
	}
	return fmt.Sprintf("%s %s", strings.Repeat("█", loaded), strings.Repeat(" ", remain))
}

//ProgressBar draw in cli
func ProgressBar(before string, after string) func(received int64, readed int64, total int64, duration float64, start int64, end int64) {
	return func(received int64, readed int64, total int64, duration float64, start int64, end int64) {
		loaded := float64(start+received) / float64(end) * 100
		speed := float64(received) / 1024 / duration
		remain := float64(total-received) / 1024 / speed
		fmt.Printf("\r%s%s%.1f%% %s/%s %.2fKB/s %.1f %.1f%s  ", before, Bar(int(loaded), 25), loaded, ByteFormat(uint64(start+readed)), ByteFormat(uint64(start+received)), speed, duration, remain, after)
	}
}

//BoolString quick
func BoolString(b bool, s, s1 string) string {
	if b {
		return s
	}
	return s1
}

//ByteFormat man read
func ByteFormat(bytes uint64) string {
	unit := [...]string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	if bytes >= 1024 {
		e := math.Floor(math.Log(float64(bytes)) / math.Log(float64(1024)))
		return fmt.Sprintf("%.2f%s", float64(bytes)/math.Pow(1024, math.Floor(e)), unit[int(e)])
	}
	return fmt.Sprintf("%d%s", bytes, unit[0])
}

//GetStorePath give the save path
func GetStorePath(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fileName := strings.Trim(path.Base(u.Path), " /.")
	if fileName == "" {
		fileName = strings.Trim(u.RawQuery, " /.")
	}
	if fileName == "" {
		fileName = "index"
	}
	return filepath.Join(dir, fileName), nil
}

//GetContinue create file or give file size and hanle
func GetContinue(saveas string) (*os.File, int64, error) {
	if stat, err := os.Stat(saveas); os.IsNotExist(err) {
		f, err := os.Create(saveas)
		return f, 0, err
	} else if !stat.IsDir() {
		file, err := os.OpenFile(saveas, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		return file, stat.Size(), err
	}
	return nil, 0, fmt.Errorf("%s is dir", saveas)
}
