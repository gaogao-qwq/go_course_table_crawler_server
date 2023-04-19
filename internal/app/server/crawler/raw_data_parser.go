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
	"strconv"
	"strings"
)

// RawCourseInfo.Id example:
// TD12_0
// TD16_0

// RawCourseInfo.Title example:
// 从“愚昧”到“科学”：科学技术简史[508000056.01] (徐绍铜);(6,超星尔雅(网课))
// 毛泽东思想和中国特色社会主义理论体系概论[110010021.35] (刘于亮);(6,教I-D401(D-M))
// 大学体育Ⅳ[B340007.57] (陈晓燕);(4)
// 思想道德与法治[110010014.02] (穆建叶,曾广志);(9)
// 大学体育Ⅳ[B340007.57] (陈晓燕);(7);概率论与数理统计[111010006.03] (姚雪超);(7,教II-C104(D-S))

// Parser 函数接受一个 RawCourseInfo 切片，对其处理后返回 CourseInfo 二维切片
func Parser(rawCourseInfoList []RawCourseInfo) (courseInfo [][]CourseInfo) {
	for _, info := range rawCourseInfoList {
		// 获取 TDn_0 中的 n
		id, err := strconv.Atoi(info.Id[2:strings.IndexByte(info.Id, '_')])
		if err != nil {
			id = 0
		}

		// 对 id 处理后获取当前课时开始于周几的第几节课
		sectionBegin := id - (id/12)*12 + 1
		dateNum := id/12 + 1

		// Rowspan 即 HTML 的 Table 元素中单元格的纵向延长，可以直接转化为课时长度
		sectionLength := info.Rowspan

		// 分号数量决定了当前课时有几节课程存在（若存在多个课时时则视作课程冲突）
		// 设分号数量为 n，当前课时中课程数量为 c，有 c = floor(n / 2)
		n := strings.Count(info.Title, ";")
		c := n/2 + 1

		// 缓存 Title
		title := info.Title
		var courses []CourseInfo
		for i := 0; i < c; i += 1 {
			// 课程名称、课程代号、课程任课老师在分号左侧，截取为 courseName
			courseName := title[:strings.IndexByte(title, ';')]
			// courseId 只会且必会出现在 courseName 中被中括号括起来的位置
			courseId := courseName[strings.IndexByte(courseName, '[')+1 : strings.IndexByte(courseName, ']')]
			courseName = courseName[:strings.IndexByte(courseName, '[')]

			// 课程上课周、上课教室（可能不存在）在分号右侧，截取为 courseDetail
			courseDetail := title[strings.IndexByte(title, ';')+1:]
			// 若当前课程后存在其它课，则以分号作为末尾下标，若当前课程为当前课时最后的课程则以最末位作为末尾下标
			locationName := func() string {
				if i < c-1 {
					return courseDetail[:strings.IndexByte(courseDetail, ';')]
				}
				return courseDetail
			}()

			// 上课周数一定会存在且只存在于课程地点的首位
			weekNum := func() int {
				if strings.IndexByte(locationName, ',') == -1 {
					tmp, err := strconv.Atoi(locationName[1:strings.IndexByte(locationName, ')')])
					if err != nil {
						return 0
					}
					return tmp
				}
				tmp, err := strconv.Atoi(locationName[1:strings.IndexByte(locationName, ',')])
				if err != nil {
					return 0
				}
				return tmp
			}()
			courses = append(courses, CourseInfo{
				IsEmpty:       false,
				CourseId:      courseId,
				CourseName:    courseName,
				LocationName:  locationName,
				SectionBegin:  sectionBegin,
				SectionLength: sectionLength,
				WeekNum:       weekNum,
				DateNum:       dateNum,
			})
			if i < c-1 {
				for j := 0; j < 2; j++ {
					title = title[strings.IndexByte(title, ';')+1:]
				}
			}
		}

		courseInfo = append(courseInfo, courses)
	}

	return
}
