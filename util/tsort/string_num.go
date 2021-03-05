package tsort

import (
	"regexp"
	"sort"
	"strconv"
)

type PersonWrapper struct {
	data []string
}

func (m PersonWrapper) Len() int {
	return len(m.data)
}
func (m PersonWrapper) Less(i, j int) bool {
	return numberFormString(m.data[i]) < numberFormString(m.data[j])
}
func (m PersonWrapper) Swap(i, j int) {
	m.data[i], m.data[j] = m.data[j], m.data[i]
}

func SortPerson(data []string) []string {
	sort.Sort(PersonWrapper{data})
	return data
}

func numberFormString(str string) uint64 {
	var str_num string
	for _, v := range regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(str, -1) {
		str_num += v[0]
	}
	i, _ := strconv.ParseUint(str_num, 10, 64)
	return i
}
