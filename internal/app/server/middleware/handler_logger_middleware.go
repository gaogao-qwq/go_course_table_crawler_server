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

package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func HandlerLoggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	latency := time.Since(start).Milliseconds()
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("{%s} [%s]->%s %s %dms \n", dateTime, c.Request.Method, c.Request.URL.Path, c.ClientIP(), latency)
}
