package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CourseInfo struct {
	IsEmpty      bool   `json:"isEmpty"`
	CourseId     string `json:"courseId"`
	CourseName   string `json:"courseName"`
	LocationName string `json:"locationName"`
	SectionNum   int    `json:"sectionNum"`
	WeekNum      int    `json:"weekNum"`
	DateNum      int    `json:"dateNum"`
}

type courseList struct {
	Semester string       `json:"semester"`
	List     []CourseInfo `json:"list"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	// 对无验证头的 GET 请求返回 StatusForbidden
	if r.Method == "GET" && r.Header["Authorization"] == nil {
		fmt.Println("无验证头的 GET 请求")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 对无表单头的 GET 请求返回 StatusBadRequest
	//if r.ParseForm() == nil {
	//	fmt.Println("无表单头的 GET 请求")
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	courseTable := []CourseInfo{
		{IsEmpty: false, CourseId: "123", CourseName: "wtf", LocationName: "1-C101", SectionNum: 1, WeekNum: 1, DateNum: 1},
	}
	err := json.NewEncoder(w).Encode(courseTable)
	if err != nil {
		return
	}
}
