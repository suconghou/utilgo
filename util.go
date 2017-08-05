package utilgo

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
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
func ProgressBar(before string, after string, hook func(loaded float64, speed float64, remain float64), writer io.Writer) func(received int64, readed int64, total int64, duration float64, start int64, end int64) {
	return func(received int64, readed int64, total int64, duration float64, start int64, end int64) {
		loaded := float64(start+received) / float64(end) * 100
		speed := float64(received) / 1024 / duration
		remain := float64(total-received) / 1024 / speed
		if hook != nil {
			hook(float64(start+readed)/float64(end)*100, speed, remain)
		}
		if writer == nil {
			writer = os.Stdout
		}
		fmt.Fprintf(writer, "\r\033[2K\r%s%s%.1f%% %s/%s/%s %.2fKB/s %.1f %.1f%s", before, Bar(int(loaded), 25), loaded, ByteFormat(uint64(start+readed)), ByteFormat(uint64(start+received)), ByteFormat(uint64(total)), speed, duration, remain, after)
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

// StringPadding padding str to given len
func StringPadding(str string, le int) string {
	l := le - len(str)
	if l > 0 {
		for i := 0; i < l; i++ {
			str = str + " "
		}
	}
	return str
}

// DateFormat form given date
func DateFormat(times int64) string {
	return time.Unix(times, 0).Format("2006/01/02 15:04:05")
}

//GetStorePath give the save path from the url or the file path
func GetStorePath(urlStr string) (string, error) {
	var fileName = urlStr
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if IsURL(urlStr) {
		u, err := url.Parse(urlStr)
		if err != nil {
			return "", err
		}
		fileName = strings.Trim(path.Base(u.Path), " /.")
		if fileName == "" {
			fileName = strings.Trim(u.RawQuery, " /.")
		}
	} else {
		fileName = filepath.Base(fileName)
	}
	if fileName == "" {
		fileName = "index"
	}
	return filepath.Join(dir, fileName), nil
}

// GetOpenFile just open file for read
func GetOpenFile(file string) (*os.File, error) {
	var fullpath = file
	if !filepath.IsAbs(fullpath) {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		fullpath = filepath.Join(dir, file)
	}
	return os.Open(fullpath)
}

//GetContinue create file or give file size and hanle
func GetContinue(fullpath string) (*os.File, int64, error) {
	if stat, err := os.Stat(fullpath); os.IsNotExist(err) {
		f, err := os.Create(fullpath)
		return f, 0, err
	} else if !stat.IsDir() {
		file, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		return file, stat.Size(), err
	}
	return nil, 0, fmt.Errorf("%s is dir", fullpath)
}

//InArray in golang
func InArray(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

// IsURL if the given string is an url
func IsURL(url string) bool {
	return regexp.MustCompile(`^[a-zA-z]+://[^\s]*$`).MatchString(url)
}

// HasFlag return if has given param
func HasFlag(flag string) bool {
	for _, item := range os.Args {
		if item == flag {
			return true
		}
	}
	return false
}

// GetParam return key value
func GetParam(key string) (string, error) {
	var catched = false
	for _, item := range os.Args {
		if catched {
			return item, nil
		}
		if item == key {
			catched = true
		}
	}
	return "", fmt.Errorf("%s value not found", key)
}

// CallPlayer player play media file
func CallPlayer(file string) {
	var player string
	if runtime.GOOS == "windows" {
		player = "PotPlayerMini.exe"
	} else {
		player = "mpv"
	}
	exec.Command(player, file).Start()
}

// PathMustHave return absolute path which must be a dir
func PathMustHave(p string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return p, err
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(pwd, p)
	}
	stat, err := os.Stat(p)
	if err != nil {
		return p, err
	}
	if !stat.Mode().IsDir() {
		return p, fmt.Errorf("%s is not directory", p)
	}
	return p, err
}

// GetFileHash return md5sum sha1sum ...
func GetFileHash(file *os.File, t string) ([]byte, error) {
	switch t {
	case "md5":
		h := md5.New()
		_, err := io.Copy(h, file)
		if err != nil {
			return nil, err
		}
		return h.Sum(nil), nil
	case "sha1":
		h := sha1.New()
		_, err := io.Copy(h, file)
		if err != nil {
			return nil, err
		}
		return h.Sum(nil), nil
	case "sha256":
		h := sha256.New()
		_, err := io.Copy(h, file)
		if err != nil {
			return nil, err
		}
		return h.Sum(nil), nil
	default:
		h := crc32.NewIEEE()
		_, err := io.Copy(h, file)
		if err != nil {
			return nil, err
		}
		return h.Sum(nil), nil
	}
}

// JSONPut resp json
func JSONPut(w http.ResponseWriter, bs []byte, cors bool, cacheTime int) (int, error) {
	h := w.Header()
	h.Set("Content-Type", "text/json; charset=utf-8")
	if cors {
		CrossShare(h)
	}
	if cacheTime > 1 {
		UseHTTPCache(h, cacheTime)
	}
	return w.Write(bs)
}

// UseHTTPCache set header
func UseHTTPCache(h http.Header, cacheTime int) {
	h.Set("Expires", time.Now().Add(time.Second*time.Duration(cacheTime)).Format(http.TimeFormat))
	h.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", cacheTime))
}

// CrossShare set header
func CrossShare(h http.Header) {
	h.Set("Access-Control-Max-Age", "86400")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS")
	h.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Content-Length, Accept, Accept-Encoding, Cache-Control, Expires")
}
