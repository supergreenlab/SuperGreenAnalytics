/*
 * Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
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

package series

import (
	"encoding/json"

	"github.com/SuperGreenLab/Analytics/internal/data/db"
	"github.com/SuperGreenLab/Analytics/internal/server/middlewares"
	"github.com/SuperGreenLab/Analytics/internal/services/pubsub"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	client   influxdb2.Client
	writeAPI api.WriteAPI
	_        = pflag.String("influxdbtoken", "", "Influx DB token")
)

func init() {
	viper.SetDefault("InfluxDBToken", "")
}

func listenEventsCreated() {
	ch := pubsub.SubscribeOject("insert.events")
	for c := range ch {
		e := c.(middlewares.InsertMessage).Object.(*db.Event)
		id := c.(middlewares.InsertMessage).ID

		p := influxdb2.NewPointWithMeasurement(e.Type)
		params := map[string]interface{}{}
		if err := json.Unmarshal([]byte(e.Params), &params); err != nil {
			logrus.Errorf("json.Unmarshal in listenEventsCreated %q - id: %s e: %+v", err, id, e)
		}
		p = p.AddTag("sourceID", e.SourceID.UUID.String())
		p = p.AddTag("sessionID", e.SessionID.UUID.String())
		p = p.AddTag("visitorID", e.VisitorID.UUID.String())
		p = p.SetTime(e.Date)
		for k, i := range params {
			switch v := i.(type) {
			case int, float64:
				p = p.AddField(k, v)
			case string:
				p = p.AddTag(k, v)
			}
		}
		writeAPI.WritePoint(p)
	}
}

func Init() {
	client = influxdb2.NewClient("http://influxdb:8086", viper.GetString("InfluxDBToken"))
	writeAPI = client.WriteAPI("Hackerman", "analytics")

	go listenEventsCreated()
}
