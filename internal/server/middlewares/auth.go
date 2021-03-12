/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package middlewares

import (
	"github.com/rileyr/middleware/wares"

	"github.com/spf13/pflag"

	"github.com/rileyr/middleware"
	"github.com/spf13/viper"
)

var (
	_ = pflag.String("logrequests", "true", "Set to false in production") // TODO move this somewhere else
)

func init() {
	viper.SetDefault("LogRequests", "true")
}

// AnonStack - allows anonymous connection
func AnonStack() middleware.Stack {
	anon := middleware.NewStack()
	if viper.GetString("LogRequests") == "true" {
		anon.Use(wares.Logging)
	}
	anon.Use(CreateDBSession)
	return anon
}
