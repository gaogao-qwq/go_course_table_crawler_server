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

type AuthorizationError struct {
	Err string
}

func (e AuthorizationError) Error() string {
	return e.Err
}

type CourseTable struct {
	Row  int            `json:"row"`
	Col  int            `json:"col"`
	Week int            `json:"week"`
	Data [][]CourseInfo `json:"data"`
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

type Crawler interface {
	loginTasks() error
	getSemesterList() ([]Semester, error)
	selectSemester(string) error
	setWeekNum(int) error
	getCourseTableSize() (int, int, int, error)
	getCourseTable() (CourseTable, error)
}
