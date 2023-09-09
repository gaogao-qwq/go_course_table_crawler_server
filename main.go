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

package main

import (
	"course_table_server/internal/app/server/config"
	"course_table_server/internal/app/server/handler"
	"course_table_server/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	e := gin.Default()
	e.GET("/v1/semester-list", middleware.AuthMiddleware, handler.SemesterListHandler)
	e.GET("/v1/course-table", middleware.AuthMiddleware, handler.CourseTableHandler)
	e.GET("/v1/test", handler.TestHandler)
	log.Fatal(e.Run(config.Address + ":" + config.Port))
}
