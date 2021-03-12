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
	"context"
	"net/http"

	"github.com/SuperGreenLab/Analytics/internal/data/db"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

// SessContextKey - context key which stores the DB session object
type SessContextKey struct{}

// CreateDBSession - Creates a DB session and stores it in the context
func CreateDBSession(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), SessContextKey{}, db.Sess)
		fn(w, r.WithContext(ctx), p)
	}
}

// InsertedIDContextKey - context key which stores the inserted object's ID
type InsertedIDContextKey struct{}

// InsertObject - Insert the payload object to DB
func InsertObject(collection string) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(ObjectContextKey{})
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			col := sess.Collection(collection)
			id, err := col.Insert(o)
			if err != nil {
				logrus.Errorf("Insert in InsertObject %q - %s %+v", err, collection, o)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), InsertedIDContextKey{}, uuid.FromStringOrNil(string(id.([]uint8))))
			fn(w, r.WithContext(ctx), p)
		}
	}
}
