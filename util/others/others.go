package others

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"utilware/decimal"
)

type (
	CutBackType struct {
		AllPage uint64 `json:"all_page"`
		DataLen uint64 `json:"data_len"`
		NowPage uint64 `json:"now_page"`
	}
)

func NewTimeticker(micro_second uint64, ticker_func func()) *time.Ticker {
	ticker := time.NewTicker(time.Duration(micro_second) * time.Microsecond)
	go func() {
		for range ticker.C {
			ticker_func()
		}
	}()
	return ticker
}

func StringToTimezone(str string) (*time.Location, error) {
	if data := strings.Split(str, ":"); len(data) == 2 {
		if n, err := strconv.Atoi(data[1]); err == nil {
			return time.FixedZone(data[0], int((time.Duration(n) * time.Hour).Seconds())), nil
		}
	}
	return nil, errors.New("split error")
}

func IntToTimezone(zone string, i int) *time.Location {
	return time.FixedZone(zone, int((time.Duration(i) * time.Hour).Seconds()))
}

func GetDirsAndFiles(path string) (dirs []string, files []string, e error) {
	if path, e = filepath.Abs(path); e == nil {
		fs, e := ioutil.ReadDir(path)
		if e != nil {
			return nil, nil, e
		}
		for _, file := range fs {
			if file.IsDir() {
				dirs = append(dirs, file.Name())
			} else {
				files = append(files, file.Name())
			}
		}
	}
	return
}

func FloatDecimal(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func DelRepeatStrings(strings []string, vals []string) []string {
	for i := 0; i < len(strings); i++ {
		for o := 0; o < len(vals); o++ {
			if strings[i] == vals[o] {
				strings = append(strings[:i], strings[i+1:]...)
			}
		}
	}
	return strings
}

func CutInterfaces(datas []interface{}, one_page_num, page_num uint64) (*CutBackType, []interface{}) {
	if one_page_num <= 0 {
		one_page_num = 10
	}

	if page_num <= 0 {
		page_num = 1
	}

	start := one_page_num*page_num - one_page_num
	stop := one_page_num * page_num
	info := &CutBackType{}

	info.DataLen = uint64(len(datas))

	if start >= info.DataLen || start < 0 {
		start = 0
	}
	if stop > info.DataLen || stop <= 0 {
		stop = info.DataLen
	}
	if (info.DataLen - stop) < one_page_num {
		stop = info.DataLen
	}

	info.NowPage = stop / one_page_num
	info.AllPage = info.DataLen / one_page_num

	if info.AllPage <= 0 {
		info.AllPage = 1
	}
	if info.NowPage <= 0 {
		info.NowPage = 1
	}

	return info, datas[start:stop]
}

func CheckOnlyAz09(s string) bool {
	exp, _ := regexp.Compile(`[a-zA-Z0-9_]`)
	return len(exp.ReplaceAllString(s, "")) == 0
}

func ReadFile(path string) (v []byte, e error) {
	if path, e = filepath.Abs(path); e == nil {
		if v, e = ioutil.ReadFile(path); e == nil {
			return v, nil
		}
	}
	return v, e
}

func WriteFile(path string, data []byte) (e error) {
	if path, e = filepath.Abs(path); e == nil {
		if e = ioutil.WriteFile(path, data, os.ModePerm); e == nil {
			return nil
		}
	}
	return e
}

func BytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func Uint64sHasUint64(vals []uint64, val uint64) (data []int) {
	for k, v := range vals {
		if v == val {
			data = append(data, k)
		}
	}
	return
}

func Uint64sDelUint64(vals []uint64, val uint64) []uint64 {
	for k, v := range vals {
		if v == val {
			vals = append(vals[:k], vals[k+1:]...)
		}
	}
	return vals
}

func InterfacesDeduplication(datas []interface{}) (vals []interface{}) {
	for i := 0; i < len(datas); i++ {
		repeat := false
		for j := i + 1; j < len(datas); j++ {
			if datas[i] == datas[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			vals = append(vals, datas[i])
		}
	}
	return vals
}

func FileExtName(fp string) string {
	return path.Ext(path.Base(fp))
}

func FullCPU(b bool) {
	if b {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(1)
	}
}

func AbsPath(path string) string {
	if v, err := filepath.Abs(path); err == nil {
		return v
	}
	return path
}

func MoveFile(source, to string) error {
	source, to = AbsPath(source), AbsPath(to)
	if fl, e := os.Stat(to); e == nil {
		if fl.IsDir() == true {
			os.RemoveAll(to)
		} else {
			os.Remove(to)
		}
	}
	return os.Rename(source, to)
}

func IsExist(path string) bool {
	_, e := os.Stat(AbsPath(path))
	return e == nil
}

func Mkdir(path string) error {
	path, e := filepath.Abs(path)
	if e != nil {
		return e
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	return os.MkdirAll(path, os.ModePerm)
}

func ReadConfig(path string, read_only bool, data interface{}) (interface{}, error) {
	path, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}

	if e = json.Unmarshal(b, &data); e != nil {
		return nil, e
	}

	if !read_only {
		WriteConfig(path, data)
	}
	return data, nil
}

func WriteConfig(path string, data interface{}) error {
	path, e := filepath.Abs(path)
	if e != nil {
		return e
	}

	byte_json, e := json.Marshal(data)
	if e != nil {
		return e
	}

	var buf bytes.Buffer
	e = json.Indent(&buf, byte_json, "", "    ")
	if e != nil {
		return e
	}

	return ioutil.WriteFile(path, buf.Bytes(), os.ModePerm)
}

func InitConfigDefault(path string, read_only bool, data, default_data interface{}) error {
	if _, err := os.Stat(AbsPath(path)); err != nil {
		if e := WriteConfig(path, default_data); e != nil {
			return e
		}
	}
	_, e := ReadConfig(path, read_only, data)
	return e
}

func DivUint64ToFloat64(x, d uint64) float64 {
	f, _ := decimal.NewFromUint64(x).Div(decimal.NewFromUint64(d)).Float64()
	return f
}

func DivUint64ToString(x, d uint64) string {
	return decimal.NewFromUint64(x).Div(decimal.NewFromUint64(d)).String()
}

func DivFloat64ToString(x, d float64) string {
	return decimal.NewFromFloat(x).Div(decimal.NewFromFloat(d)).String()
}

func MulUint64ToUint64(f, d uint64) uint64 {
	return decimal.NewFromUint64(f).Mul(decimal.NewFromUint64(d)).Uint64()
}

func MulUint64ToString(f, d uint64) string {
	return decimal.NewFromUint64(f).Mul(decimal.NewFromUint64(d)).String()
}

func MulFloat64ToUint64(f, d float64) uint64 {
	return decimal.NewFromFloat(f).Mul(decimal.NewFromFloat(d)).Uint64()
}

func MulStringToUint64(f, d string) uint64 {
	k, e := decimal.NewFromString(f)
	if e != nil {
		return 0
	}
	v, _ := decimal.NewFromString(d)
	if e != nil {
		return 0
	}
	return k.Mul(v).Uint64()
}
