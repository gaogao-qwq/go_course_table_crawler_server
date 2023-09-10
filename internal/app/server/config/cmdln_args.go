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

package config

import "flag"

var (
	Address  string
	Port     string
	LoginUrl string
	HomeUrl  string
)

func init() {
	flag.StringVar(&Address, "address", "0.0.0.0", "server listen address")
	flag.StringVar(&Port, "port", "56789", "server listen port")
	flag.StringVar(&LoginUrl, "loginurl", "http://targeturl/login", "login page url")
	flag.StringVar(&HomeUrl, "homeurl", "http://targeturl/home", "home page url")
	flag.Parse()
}
