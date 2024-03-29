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
	"course_table_server/internal/app/server/config"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strconv"
	"time"
)

type CourseTableCrawler struct {
	ctx      context.Context
	account  string
	password string
}

func (c CourseTableCrawler) loginTasks() (err error) {
	var location string
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Navigate(config.LoginUrl),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.SendKeys("#username", c.account, chromedp.ByID),
		chromedp.SendKeys("#password", c.password, chromedp.ByID),
		chromedp.Sleep(time.Second),
		chromedp.Click("#loginForm > table.logintable > tbody > tr:nth-child(6) > td > input", chromedp.NodeVisible),
		chromedp.Sleep(time.Second),
		chromedp.Location(&location),
	})
	if err != nil {
		return
	}
	if location != config.HomeUrl {
		return AuthorizationError{}
	}
	return
}

func (c CourseTableCrawler) getSemesterList() (semesterList []Semester, err error) {
	var nodes []*cdp.Node
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Sleep(time.Second),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.Click("#menu_panel > ul > li.expand > ul > div > li:nth-child(3) > a"),
		chromedp.Click(".calendar-text-state-default", chromedp.ByQuery),
		chromedp.Nodes(".calendar-bar-td-blankBorder", &nodes, chromedp.ByQueryAll),
	})
	if err != nil {
		return
	}

	for _, node := range nodes {
		if node.Children == nil {
			break
		}
		semesterList = append(semesterList, Semester{
			Value:       node.Children[0].NodeValue,
			Index:       node.Attributes[3],
			SemesterId1: "",
			SemesterId2: "",
		})
	}

	for i, semester := range semesterList {
		var semesterIdNode1 []*cdp.Node
		var semesterIdNode2 []*cdp.Node
		_ = chromedp.Run(c.ctx, chromedp.Tasks{
			chromedp.Click(".calendar-bar-td-blankBorder[index=\""+semester.Index+"\"]", chromedp.ByQuery),
			chromedp.Nodes("#semesterCalendar_termTb > tbody > tr:nth-child(1) > td", &semesterIdNode1, chromedp.AtLeast(0)),
			chromedp.Nodes("#semesterCalendar_termTb > tbody > tr:nth-child(2) > td", &semesterIdNode2, chromedp.AtLeast(0)),
			chromedp.Sleep(time.Second / 2),
		})

		if len(semesterIdNode1) > 0 {
			semesterList[i].SemesterId1 = semesterIdNode1[0].Attributes[3]
		}
		if len(semesterIdNode2) > 0 {
			semesterList[i].SemesterId2 = semesterIdNode2[0].Attributes[3]
		}
	}

	return
}

func (c CourseTableCrawler) selectSemester(semesterId string) (err error) {
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Click("#menu_panel > ul > li.expand > ul > div > li:nth-child(3) > a"),
		chromedp.SetValue("#semesterCalendar_target", semesterId, chromedp.ByQuery),
		chromedp.Click("#courseTableForm > div:nth-child(2) > input[type=submit]:nth-child(9)", chromedp.ByQuery),
		chromedp.Sleep(2 * time.Second),
	})
	if err != nil {
		return
	}
	return
}

func (c CourseTableCrawler) setWeekNum(weekNum int) (err error) {
	var placeholderNode []*cdp.Node
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.SetValue("#startWeek", strconv.Itoa(weekNum), chromedp.ByID),
		chromedp.Nodes(
			"#manualArrangeCourseTable > tbody > tr:nth-child(1) > td:nth-child(1)",
			&placeholderNode,
			chromedp.ByQuery,
		),
		chromedp.Sleep(time.Second),
		chromedp.WaitVisible("#manualArrangeCourseTable > tbody > tr:nth-child(1) > td:nth-child(1)", chromedp.ByQuery),
	})
	if err != nil {
		return
	}
	return
}

func (c CourseTableCrawler) getCourseTableSize() (row int, col int, week int, err error) {
	var placeholderNode []*cdp.Node

	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Nodes("#startWeek", &placeholderNode, chromedp.ByQuery),
	})
	if err != nil {
		return
	}
	week = int(placeholderNode[0].ChildNodeCount - 1)

	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Nodes("#manualArrangeCourseTable > tbody", &placeholderNode, chromedp.ByQuery),
	})
	if err != nil {
		return
	}
	row = int(placeholderNode[0].ChildNodeCount)

	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Nodes("#manualArrangeCourseTable > thead > tr", &placeholderNode, chromedp.ByQuery),
	})
	if err != nil {
		return
	}
	col = int(placeholderNode[0].ChildNodeCount - 1)

	return
}

func (c CourseTableCrawler) getCourseTable() (courseTable CourseTable, err error) {
	var (
		tableBodyNode     []*cdp.Node
		placeholderNode   []*cdp.Node
		rawCourseInfoList []RawCourseInfo
	)

	row, col, weekNum, err := c.getCourseTableSize()
	if err != nil {
		return courseTable, err
	}
	courseTable.Row = row
	courseTable.Col = col
	courseTable.Week = weekNum

	for week := 1; week <= weekNum; week += 1 {
		err = c.setWeekNum(week)
		if err != nil {
			return courseTable, err
		}
		err = chromedp.Run(c.ctx, chromedp.Tasks{
			chromedp.Nodes(
				"#manualArrangeCourseTable > tbody > tr:nth-child(1) > td:nth-child(1)",
				&placeholderNode,
				chromedp.ByQuery,
			),
			chromedp.Nodes("#manualArrangeCourseTable > tbody", &tableBodyNode, chromedp.ByQuery),
		})
		if err != nil {
			return courseTable, err
		}

		for nthTR := int64(1); nthTR <= tableBodyNode[0].ChildNodeCount; nthTR += 1 {
			for nthTD := int64(1); nthTD <= tableBodyNode[0].Children[nthTR-1].ChildNodeCount; nthTD += 1 {
				err = chromedp.Run(c.ctx, chromedp.Tasks{
					chromedp.Nodes(
						"#manualArrangeCourseTable > tbody > tr:nth-child("+strconv.Itoa(int(nthTR))+") > td:nth-child("+strconv.Itoa(int(nthTD))+")",
						&placeholderNode,
						chromedp.ByQuery,
					),
				})
				if err != nil {
					return courseTable, err
				}

				title, isExist := placeholderNode[0].Attribute("title")
				if !isExist {
					continue
				}
				id := placeholderNode[0].AttributeValue("id")
				rowspan, err := strconv.Atoi(placeholderNode[0].AttributeValue("rowspan"))
				if err != nil {
					return courseTable, err
				}

				rawCourseInfoList = append(rawCourseInfoList, RawCourseInfo{
					Id:      id,
					Rowspan: rowspan,
					Title:   title,
				})
			}
		}

	}

	courseTable.Data = Parser(rawCourseInfoList)
	return courseTable, nil
}
