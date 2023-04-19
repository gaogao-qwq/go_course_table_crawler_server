// An web crawler and API server implementation for the course table app below:
// https://github.com/gaogao-qwq/flutter_course_table_demo
// Copyright (C) 2023 Zhihao Zhou
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

func NewCourseTableCrawler(account string, password string) CourseTableCrawler {
	ctx, _ := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 encors/0.0.6"),
			chromedp.Flag("headless", true),
		)...,
	)
	return CourseTableCrawler{
		ctx:      ctx,
		account:  account,
		password: password,
	}
}

func Authorizer(account string, password string) (err error) {
	crawler := NewCourseTableCrawler(account, password)
	var cancel context.CancelFunc
	crawler.ctx, cancel = context.WithTimeout(crawler.ctx, 10*time.Second)
	defer cancel()
	crawler.ctx, _ = chromedp.NewContext(crawler.ctx)

	err = crawler.loginTasks()
	if err != nil {
		return
	}
	return
}

func GetSemesterList(account string, password string) (semesterList []Semester, err error) {
	crawler := NewCourseTableCrawler(account, password)
	var cancel context.CancelFunc
	crawler.ctx, cancel = context.WithTimeout(crawler.ctx, 20*time.Second)
	defer cancel()

	crawler.ctx, _ = chromedp.NewContext(crawler.ctx)
	err = crawler.loginTasks()
	if err != nil {
		return
	}

	semesterList, err = crawler.getSemesterList()
	if err != nil {
		return
	}
	return
}

func GetCourseTable(account string, password string, semesterId string) (courseTable CourseTable, err error) {
	crawler := NewCourseTableCrawler(account, password)
	var cancel context.CancelFunc
	crawler.ctx, cancel = context.WithTimeout(crawler.ctx, 60*time.Second)
	defer cancel()

	crawler.ctx, _ = chromedp.NewContext(crawler.ctx)
	err = crawler.loginTasks()
	if err != nil {
		return courseTable, err
	}

	err = crawler.selectSemester(semesterId)
	if err != nil {
		return courseTable, err
	}

	courseTable, err = crawler.getCourseTable()
	if err != nil {
		return courseTable, err
	}
	return
}
