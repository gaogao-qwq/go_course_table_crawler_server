package crawler

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strconv"
	"time"
)

type CourseTableCrawler struct {
	ctx      context.Context
	url      string
	account  string
	password string
}

func (c CourseTableCrawler) loginTasks() (err error) {
	var location string
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Navigate(c.url),
		chromedp.WaitVisible("body > div.foot"),
		chromedp.SendKeys("#username", c.account, chromedp.ByID),
		chromedp.SendKeys("#password", c.password, chromedp.ByID),
		chromedp.Sleep(time.Second),
		chromedp.Click("#loginForm > table.logintable > tbody > tr:nth-child(6) > td > input", chromedp.NodeVisible),
		chromedp.Sleep(time.Second),
		chromedp.Location(&location),
		chromedp.Click("#menu_panel > ul > li.expand > ul > div > li:nth-child(3) > a"),
	})
	if err != nil {
		return
	}
	if location != "http://jw.gzgs.edu.cn/eams/home.action" {
		return AuthorizationError{}
	}
	return
}

func (c CourseTableCrawler) getSemesterList() (semesterList []Semester, err error) {
	var nodes []*cdp.Node
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.Sleep(time.Second),
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
		err = chromedp.Run(c.ctx, chromedp.Tasks{
			chromedp.Click(".calendar-bar-td-blankBorder[index=\""+semester.Index+"\"]", chromedp.ByQuery),
			chromedp.Nodes("#semesterCalendar_termTb > tbody > tr:nth-child(1) > td", &semesterIdNode1, chromedp.ByQuery),
			chromedp.Nodes("#semesterCalendar_termTb > tbody > tr:nth-child(2) > td", &semesterIdNode2, chromedp.ByQuery),
		})
		semesterList[i].SemesterId1 = semesterIdNode1[0].Attributes[3]
		semesterList[i].SemesterId2 = semesterIdNode2[0].Attributes[3]
		if err != nil {
			return
		}
	}

	return
}

func (c CourseTableCrawler) selectSemester(semesterId string) (err error) {
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.SetValue("#semesterCalendar_target", semesterId, chromedp.ByID),
		chromedp.Click("#courseTableForm > div:nth-child(2) > input[type=submit]:nth-child(9)", chromedp.ByQuery),
		chromedp.Sleep(time.Second),
	})
	if err != nil {
		return
	}
	return
}

func (c CourseTableCrawler) setWeekNum(weekNum int) (err error) {
	err = chromedp.Run(c.ctx, chromedp.Tasks{
		chromedp.SetValue("#startWeek", strconv.Itoa(weekNum), chromedp.ByID),
		chromedp.Sleep(time.Second / 2),
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
