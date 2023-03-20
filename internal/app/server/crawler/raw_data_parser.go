package crawler

import (
	"strconv"
	"strings"
)

// 从“愚昧”到“科学”：科学技术简史[508000056.01] (徐绍铜);(6,超星尔雅(网课))
// 毛泽东思想和中国特色社会主义理论体系概论[110010021.35] (刘于亮);(6,教I-D401(D-M))
// 大学体育Ⅳ[B340007.57] (陈晓燕);(4)

func Parser(rawCourseInfoList []RawCourseInfo) []CourseInfo {
	var courseTable []CourseInfo
	for _, info := range rawCourseInfoList {
		id, err := strconv.Atoi(info.Id[2:strings.IndexByte(info.Id, '_')])
		if err != nil {
			id = 0
		}

		courseId := info.Title[strings.IndexByte(info.Title, '[')+1 : strings.IndexByte(info.Title, ']')]
		courseName := info.Title[:strings.IndexByte(info.Title, ';')]
		locationName := func() string {
			if strings.LastIndexByte(info.Title, ',') == -1 {
				return ""
			}
			return info.Title[strings.LastIndexByte(info.Title, ',')+1 : strings.LastIndexByte(info.Title, ')')]
		}()
		sectionBegin := id - (id/12)*12 + 1
		sectionLength := info.Rowspan
		weekNum := func() int {
			if strings.LastIndexByte(info.Title, ',') == -1 {
				tmp, err := strconv.Atoi(info.Title[strings.LastIndexByte(info.Title, ';')+2 : strings.LastIndexByte(info.Title, ')')])
				if err != nil {
					return -1
				}
				return tmp
			}
			tmp, err := strconv.Atoi(info.Title[strings.LastIndexByte(info.Title, ';')+2 : strings.LastIndexByte(info.Title, ',')])
			if err != nil {
				return -1
			}
			return tmp
		}()
		dateNum := id/12 + 1

		courseTable = append(courseTable, CourseInfo{
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
	return courseTable
}

//func main() {
//	rawCourseInfoList := []RawCourseInfo{{
//		id:      "TD70_0",
//		rowspan: 2,
//		title:   "从“愚昧”到“科学”：科学技术简史[508000056.01] (徐绍铜);(6,超星尔雅(网课))",
//	}, {
//		id:      "TD0_0",
//		rowspan: 2,
//		title:   "毛泽东思想和中国特色社会主义理论体系概论[110010021.35] (刘于亮);(4,教I-D401(D-M))",
//	}}
//	courseTable := Parser(rawCourseInfoList)
//	fmt.Println(courseTable)
//}
