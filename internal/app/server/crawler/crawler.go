package crawler

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strconv"
	"time"
)

type AuthorizationError struct {
	Err string
}

func (e AuthorizationError) Error() string {
	return e.Err
}

type CourseTable struct {
	Row  int          `json:"row"`
	Col  int          `json:"col"`
	Week int          `json:"week"`
	Data []CourseInfo `json:"data"`
}

type CourseInfo struct {
	IsEmpty       bool   `json:"isEmpty"`
	CourseId      string `json:"courseId"`
	CourseName    string `json:"courseName"`
	LocationName  string `json:"locationName"`
	SectionBegin  int    `json:"sectionBegin"`
	SectionLength int    `json:"sectionLength"`
	WeekNum       int    `json:"weekNum"`
	DateNum       int    `json:"dateNum"`
}

type Semester struct {
	Value       string
	Index       string
	SemesterId1 string
	SemesterId2 string
}

type RawCourseInfo struct {
	Id      string
	Rowspan int
	Title   string
}

func initCrawler() (ctx context.Context) {
	ctx, _ = chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 encors/0.0.6"),
			chromedp.Flag("headless", true),
		)...,
	)
	return
}

func loginTasks(ctx *context.Context, url string, account string, password string) (err error) {
	var location string
	err = chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible("body > div.foot"),
		chromedp.SendKeys("#username", account, chromedp.ByID),
		chromedp.SendKeys("#password", password, chromedp.ByID),
		chromedp.Sleep(time.Second),
		chromedp.Click("#loginForm > table.logintable > tbody > tr:nth-child(6) > td > input", chromedp.NodeVisible),
		chromedp.Sleep(time.Second),
		chromedp.Location(&location),
	})
	if err != nil {
		return
	}
	if location != "http://jw.gzgs.edu.cn/eams/home.action" {
		return AuthorizationError{Err: "Wrong username or password"}
	}
	return
}

func getSemesterList(ctx *context.Context) (semesterList []Semester, err error) {
	var nodes []*cdp.Node
	err = chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.Click("#menu_panel > ul > li.expand > ul > div > li:nth-child(3) > a"),
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
		err = chromedp.Run(*ctx, chromedp.Tasks{
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

func selectSemester(ctx *context.Context, semesterId string) (err error) {
	err = chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.SetValue("#semesterCalendar_target", semesterId, chromedp.ByID),
		chromedp.Click("#courseTableForm > div:nth-child(2) > input[type=submit]:nth-child(9)", chromedp.ByQuery),
		chromedp.Sleep(time.Second),
	})
	if err != nil {
		return
	}
	return
}

func setWeekNum(ctx *context.Context, weekNum int) (err error) {
	err = chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.SetValue("#startWeek", strconv.Itoa(weekNum), chromedp.ByID),
		chromedp.Sleep(time.Second / 2),
	})
	if err != nil {
		return
	}
	return
}

func getCourseTableSize(ctx *context.Context) (row int, col int, week int, err error) {
	var placeholderNode []*cdp.Node

	err = chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.Nodes("#startWeek", &placeholderNode, chromedp.ByQuery),
	})
	if err != nil {
		return
	}
	week = int(placeholderNode[0].ChildNodeCount - 1)

	err = chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.Nodes("#manualArrangeCourseTable > tbody", &placeholderNode, chromedp.ByQuery),
	})
	if err != nil {
		return
	}
	row = int(placeholderNode[0].ChildNodeCount)

	err = chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.Nodes("#manualArrangeCourseTable > thead > tr", &placeholderNode, chromedp.ByQuery),
	})
	if err != nil {
		return
	}
	col = int(placeholderNode[0].ChildNodeCount - 1)

	return
}

func getCourseTable(ctx *context.Context) (courseTable CourseTable, err error) {
	var (
		tableBodyNode     []*cdp.Node
		placeholderNode   []*cdp.Node
		rawCourseInfoList []RawCourseInfo
	)

	row, col, weekNum, err := getCourseTableSize(ctx)
	if err != nil {
		return
	}
	courseTable.Row = row
	courseTable.Col = col
	courseTable.Week = weekNum

	for week := 1; week <= weekNum; week += 1 {
		err = setWeekNum(ctx, week)
		if err != nil {
			return
		}
		err = chromedp.Run(*ctx, chromedp.Tasks{
			chromedp.Nodes(
				"#manualArrangeCourseTable > tbody > tr:nth-child(1) > td:nth-child(1)",
				&placeholderNode,
				chromedp.ByQuery,
			),
			chromedp.Nodes("#manualArrangeCourseTable > tbody", &tableBodyNode, chromedp.ByQuery),
		})
		if err != nil {
			return
		}

		for nthTR := int64(1); nthTR <= tableBodyNode[0].ChildNodeCount; nthTR += 1 {
			for nthTD := int64(1); nthTD <= tableBodyNode[0].Children[nthTR-1].ChildNodeCount; nthTD += 1 {
				err = chromedp.Run(*ctx, chromedp.Tasks{
					chromedp.Nodes(
						"#manualArrangeCourseTable > tbody > tr:nth-child("+strconv.Itoa(int(nthTR))+") > td:nth-child("+strconv.Itoa(int(nthTD))+")",
						&placeholderNode,
						chromedp.ByQuery,
					),
				})
				if err != nil {
					return
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
	return
}

func Crawler(url string, account string, password string) (courseTable CourseTable, err error) {
	ctx := initCrawler()
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	ctx, _ = chromedp.NewContext(ctx)

	err = loginTasks(&ctx, url, account, password)
	if err != nil {
		return
	}

	semesterList, err := getSemesterList(&ctx)
	if err != nil {
		return
	}

	err = selectSemester(&ctx, semesterList[6].SemesterId2)
	if err != nil {
		return
	}

	courseTable, err = getCourseTable(&ctx)
	if err != nil {
		return
	}

	return
}
