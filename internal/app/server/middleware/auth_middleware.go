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
	"course_table_server/internal/app/server/crawler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err := crawler.Authorizer(username, password)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	c.Next()
}
