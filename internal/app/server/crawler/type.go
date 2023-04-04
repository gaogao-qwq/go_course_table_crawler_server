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
