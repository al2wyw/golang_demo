package main

import (
	"fmt"
	"testing"
	"time"
)

func AddYears(t time.Time, years int) time.Time {
	return AddDate(t, years, 0, 0)
}

func AddMonths(t time.Time, months int) time.Time {
	return AddDate(t, 0, months, 0)
}

// AddDate 解决 go time包 AddDate() 添加年份/月份溢出到下一个月的问题.
// 例如:
//
//	2024-02-29 AddDate(1,0,0) 期望结果:2025-02-28 实际结果:2025-03-01
//	2024-08-31 AddDate(0,1,1) 期望结果:2024-09-30 实际结果:2024-10-01
func AddDate(t time.Time, year, month, day int) time.Time {
	//先跳到目标月的1号
	targetDate := t.AddDate(year, month, -t.Day()+1)
	//获取目标月的临界值
	targetDay := targetDate.AddDate(0, 1, -1).Day()
	//对比临界值与源日期值，取最小的值
	if targetDay > t.Day() {
		targetDay = t.Day()
	}
	//最后用目标月的1号加上目标值和入参的天数
	targetDate = targetDate.AddDate(0, 0, targetDay-1+day)
	return targetDate
}

func TestAddDate(t *testing.T) {
	caseRun("2025-01-31")
	caseRun("2024-01-31")
	caseRun("2024-08-31")
	caseRun("2024-02-29")

	caseRunYear("2024-02-29")
}

func caseRun(date string) {
	// 获取当前时间
	now, _ := time.ParseInLocation("2006-01-02", date, time.Local)

	// 获取下下个月的第一天
	next := AddMonths(now, 1)

	fmt.Printf("%s->:%s\n", date, next.Format("2006-01-02"))
}

func caseRunYear(date string) {
	// 获取当前时间
	now, _ := time.ParseInLocation("2006-01-02", date, time.Local)

	// 获取下下个月的第一天
	next := AddYears(now, 1)

	fmt.Printf("%s->:%s\n", date, next.Format("2006-01-02"))
}
