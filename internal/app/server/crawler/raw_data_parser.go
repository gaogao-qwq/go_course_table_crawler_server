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

// Parser 函数接受一个 RawCourseInfo 切片，对其处理后返回 CourseInfo 切片
func Parser(rawCourseInfoList []RawCourseInfo) (courseInfo []CourseInfo) {

	for _, info := range rawCourseInfoList {
		// 获取 TDn_0 中的 n
		id, err := strconv.Atoi(info.Id[2:strings.IndexByte(info.Id, '_')])
		if err != nil {
			id = 0
		}

		// 对 id 处理后获取当前课时开始于周几的第几节课
		sectionBegin := id - (id/12)*12 + 1
		dateNum := id/12 + 1

		// RowSpan 即 HTML 的 Table 元素中单元格的纵向延长，可以直接转化为课时长度
		sectionLength := info.Rowspan

		// 课程名称、课程代号、课程任课老师在分号左侧，截取为 courseName
		courseName := info.Title[:strings.IndexByte(info.Title, ';')]
		// courseId 只会且必会出现在 courseName 中被中括号括起来的位置
		courseId := courseName[strings.IndexByte(info.Title, '[')+1 : strings.IndexByte(info.Title, ']')]

		// 课程上课周、上课教室（可能不存在）在分号右侧，截取为 courseDetail
		courseDetail := info.Title[strings.IndexByte(info.Title, ';')+1:]
		// 上课教室可能不存在，当不存在时，作为分隔符的 ',' 也不会存在
		locationName := func() string {
			if strings.LastIndexByte(courseDetail, ',') == -1 {
				return ""
			}
			return info.Title[strings.LastIndexByte(courseDetail, ',')+1 : strings.LastIndexByte(courseDetail, ')')]
		}()

		// 上课周数一定会存在且只存在于首位
		weekNum := func() int {
			if strings.LastIndexByte(courseDetail, ',') == -1 {
				tmp, err := strconv.Atoi(courseDetail[strings.LastIndexByte(courseDetail, ';')+2 : strings.LastIndexByte(courseDetail, ')')])
				if err != nil {
					return -1
				}
				return tmp
			}
			tmp, err := strconv.Atoi(courseDetail[strings.LastIndexByte(courseDetail, ';')+2 : strings.LastIndexByte(courseDetail, ',')])
			if err != nil {
				return -1
			}
			return tmp
		}()

		courseInfo = append(courseInfo, CourseInfo{
			IsEmpty:       false,
			CourseId:      courseId,
			CourseName:    courseName,
			LocationName:  locationName,
			SectionBegin:  sectionBegin,
			SectionLength: sectionLength,
			WeekNum:       weekNum,
			DateNum:       dateNum,
		})
	}

	return
}
