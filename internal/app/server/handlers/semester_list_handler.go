package handlers

import (
	"course_table_server/internal/app/server/crawler"
	"encoding/json"
	"fmt"
	"net/http"
)

func SemesterListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle /semester-list, ", "Method: ", r.Method)
	w.Header().Set("Content-Type", "application/json")

	// 不响应 非 GET 请求
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(`"message": "ONLY GET requests are allowed"`))
		if err != nil {
			panic(err)
		}
		return
	}

	// 对无验证头的 GET 请求返回 StatusUnauthorized
	if r.Method == http.MethodGet && r.Header["Authorization"] == nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`"message": "GET requests without Authorization Header are not allowed"`))
		if err != nil {
			panic(err)
		}
		return
	}

	// 从请求头中获取 Basic 验证 token
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`"message": "No basic auth present"`))
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Println("username:", username, "\npassword", password)
	semesterList, err := crawler.GetSemesterList(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
		return
	}

	err = json.NewEncoder(w).Encode(semesterList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
		return
	}
}
