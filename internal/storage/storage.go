package storage

import (
	"errors"
	"reflect"
	"reminder-manager/utils"
	"sort"
	"time"
)

var data []*Reminder

func init() {
	data = make([]*Reminder, 0)
}

func Add(rem *Reminder) {
	index := sort.Search(len(data), func(i int) bool {
		return rem.Date.Before(data[i].Date)
	})
	data = append(data, rem)
	swap := reflect.Swapper(data)
	for i := index; i < len(data); i++ {
		swap(i, len(data)-1)
	}
}

func RemindersForDays(count int) []*Reminder {
	if count < 1 {
		return nil
	}
	l := sort.Search(len(data), func(i int) bool {
		today := utils.UpToDay(time.Now()).Add(-time.Millisecond)
		return data[i].Date.After(today)
	})
	r := sort.Search(len(data), func(i int) bool {
		border := utils.UpToDay(time.Now()).Add(24*time.Hour*time.Duration(count) - time.Millisecond)
		return data[i].Date.After(border)
	})
	if l == len(data) || l == r {
		return nil
	}
	res := make([]*Reminder, r-l)
	copy(res, data[l:r])
	return res
}

func AsStrings(rem []*Reminder) []string {
	if rem == nil {
		return nil
	}
	str := make([]string, 0, len(rem))
	for _, cur := range rem {
		str = append(str, cur.ToString())
	}
	return str
}

func RemoveOutdated() int {
	outdated := OutdatedCount()
	data = data[outdated:]
	return outdated
}

func OutdatedCount() (cnt int) {
	for i := 0; i < len(data) && data[i].Date.Before(utils.UpToDay(time.Now())); i++ {
		cnt++
	}
	return cnt
}

func RemoveById(id uint64) bool {
	index, err := indexById(id)
	if err == nil {
		copy(data[index:], data[index+1:])
		data = data[:len(data)-1]
		return true
	}
	return false
}

func Edit(id uint64, newText string) bool {
	index, err := indexById(id)
	if err == nil {
		data[index].What = newText
		return true
	}
	return false
}

func indexById(id uint64) (int, error) {
	for i, cur := range data {
		if cur.Id == id {
			return i, nil
		}
	}
	return -1, errors.New("no such id")
}

func Data() []*Reminder {
	res := make([]*Reminder, len(data))
	copy(res, data)
	return res
}
