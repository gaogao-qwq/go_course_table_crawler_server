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

package handlers

import (
	"course_table_server/internal/app/server/crawler"
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle /login, ", "Method: ", r.Method)
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

	err := crawler.Authorizer(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
