package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strconv"
	"time"
)

var (
	url      = "http://jw.gzgs.edu.cn/eams/login.action"
	account  = "202110610352"
	password = "060911"
)

type Semester struct {
	value       string
	index       string
	semesterId1 string
	semesterId2 string
}

type RawCourseInfo struct {
	id      string
	rowspan int
	title   string
}

func initCrawler() context.Context {
	ctx, _ := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 encors/0.0.6"),
			chromedp.Flag("headless", false),
		)...,
	)
	return ctx
}

func loginTasks(ctx *context.Context) error {
	err := chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible("body > div.foot"),
		chromedp.SendKeys("#username", account, chromedp.ByID),
		chromedp.SendKeys("#password", password, chromedp.ByID),
		chromedp.Sleep(time.Second),
		chromedp.Click("#loginForm > table.logintable > tbody > tr:nth-child(6) > td > input", chromedp.NodeVisible),
		chromedp.Sleep(time.Second),
		chromedp.Click("#menu_panel > ul > li.expand > ul > div > li:nth-child(3) > a"),
		chromedp.Sleep(time.Second),
	})
	if err != nil {
		return err
	}
	return nil
}

func getSemesterList(ctx *context.Context) ([]Semester, error) {
	var nodes []*cdp.Node
	err := chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.Click(".calendar-text-state-default", chromedp.ByQuery),
		chromedp.Nodes(".calendar-bar-td-blankBorder", &nodes, chromedp.ByQueryAll),
	})
	if err != nil {
		return nil, err
	}

	var semesterList []Semester
	for _, node := range nodes {
		if node.Children == nil {
			break
		}
		semesterList = append(semesterList, Semester{
			value:       node.Children[0].NodeValue,
			index:       node.Attributes[3],
			semesterId1: "",
			semesterId2: "",
		})
	}

	for i, semester := range semesterList {
		var semesterIdNode1 []*cdp.Node
		var semesterIdNode2 []*cdp.Node
		err = chromedp.Run(*ctx, chromedp.Tasks{
			chromedp.Click(".calendar-bar-td-blankBorder[index=\""+semester.index+"\"]", chromedp.ByQuery),
			chromedp.Nodes("#semesterCalendar_termTb > tbody > tr:nth-child(1) > td", &semesterIdNode1, chromedp.ByQuery),
			chromedp.Nodes("#semesterCalendar_termTb > tbody > tr:nth-child(2) > td", &semesterIdNode2, chromedp.ByQuery),
		})
		semesterList[i].semesterId1 = semesterIdNode1[0].Attributes[3]
		semesterList[i].semesterId2 = semesterIdNode2[0].Attributes[3]
		if err != nil {
			return nil, err
		}
	}

	return semesterList, nil
}

func selectSemester(ctx *context.Context, semesterId string) error {
	err := chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.SetValue("#semesterCalendar_target", semesterId, chromedp.ByID),
		chromedp.Click("#courseTableForm > div:nth-child(2) > input[type=submit]:nth-child(9)", chromedp.ByQuery),
		chromedp.Sleep(time.Second),
	})
	if err != nil {
		return err
	}
	return nil
}

func setWeekNum(ctx *context.Context, weekNum int) error {
	err := chromedp.Run(*ctx, chromedp.Tasks{
		chromedp.SetValue("#startWeek", strconv.Itoa(weekNum), chromedp.ByID),
		chromedp.Sleep(time.Second),
	})
	if err != nil {
		return err
	}
	return nil
}

func getRawCourseInfo(ctx *context.Context) ([]RawCourseInfo, error) {
	var (
		tableBodyNode     []*cdp.Node
		placeholderNode   []*cdp.Node
		rawCourseInfoList []RawCourseInfo
	)

	for weekNum := 1; weekNum <= 20; weekNum += 1 {
		err := setWeekNum(ctx, weekNum)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		for nthTR := int64(1); nthTR <= tableBodyNode[0].ChildNodeCount; nthTR += 1 {
			for nthTD := int64(1); nthTD <= tableBodyNode[0].Children[nthTR-1].ChildNodeCount; nthTD += 1 {
				err := chromedp.Run(*ctx, chromedp.Tasks{
					chromedp.Nodes(
						"#manualArrangeCourseTable > tbody > tr:nth-child("+strconv.Itoa(int(nthTR))+") > td:nth-child("+strconv.Itoa(int(nthTD))+")",
						&placeholderNode,
						chromedp.ByQuery,
					),
				})
				if err != nil {
					return nil, err
				}

				title, isExist := placeholderNode[0].Attribute("title")
				if !isExist {
					continue
				}
				id := placeholderNode[0].AttributeValue("id")
				rowspan, err := strconv.Atoi(placeholderNode[0].AttributeValue("rowspan"))
				if err != nil {
					return nil, err
				}

				rawCourseInfoList = append(rawCourseInfoList, RawCourseInfo{
					id:      id,
					rowspan: rowspan,
					title:   title,
				})
			}
		}
	}

	return rawCourseInfoList, nil
}

func Crawler(url string, account string, password string) {
	ctx := initCrawler()
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	ctx, _ = chromedp.NewContext(ctx)

	err := loginTasks(&ctx)
	if err != nil {
		panic(err)
	}

	semesterList, err := getSemesterList(&ctx)
	if err != nil {
		panic(err)
	}

	err = selectSemester(&ctx, semesterList[6].semesterId1)
	if err != nil {
		panic(err)
	}

	courseInfoRawList, err := getRawCourseInfo(&ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(courseInfoRawList)
}

func main() {
	Crawler(url, account, password)
}
