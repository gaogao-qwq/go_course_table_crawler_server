package handlers

import (
	"course_table_server/internal/app/server/crawler"
	"encoding/json"
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	// 不响应 非 GET 请求
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(`message: "ONLY GET requests are allowed"`))
		if err != nil {
			panic(err)
		}
		return
	}

	// 对无验证头的 GET 请求返回 StatusUnauthorized
	if r.Method == http.MethodGet && r.Header["Authorization"] == nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`message: "GET requests without Authorization Header are not allowed"`))
		if err != nil {
			panic(err)
		}
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`message: "No basic auth present"`))
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Println("username:", username, "\npassword", password)
	courseTable, err := crawler.Crawler("http://jw.gzgs.edu.cn/eams/login.action", username, password)
	if err != nil && err.Error() == "wrong username or password" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`message: "Wrong username or password"`))
		if err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(courseTable)
	if err != nil {
		return
	}
}
