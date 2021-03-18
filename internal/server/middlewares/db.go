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
	"github.com/lib/pq"

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

type Syncable interface {
	GetSyncID() uuid.UUID
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
			ctx := r.Context()
			if err != nil {
				if s, ok := o.(Syncable); ok {
					if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
						idStruct := struct {
							ID uuid.UUID `db:"id"`
						}{}
						err := sess.Select("id").From(collection).Where("syncid = ?", s.GetSyncID()).One(&idStruct)
						if err != nil {
							logrus.Errorf("sess.Select in MultipleInsertObjects %q - %s %+v", err, collection, o)
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						ctx = context.WithValue(r.Context(), InsertedIDContextKey{}, idStruct.ID)
						fn(w, r.WithContext(ctx), p)
						return
					}
				}
				logrus.Errorf("Insert in InsertObject %q - %s %+v", err, collection, o)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				ctx = context.WithValue(r.Context(), InsertedIDContextKey{}, uuid.FromStringOrNil(string(id.([]uint8))))
				fn(w, r.WithContext(ctx), p)
			}
		}
	}
}

// MultipleInsertedIDsContextKey - context key which stores the inserted object's ID
type MultipleInsertedIDsContextKey struct{}
type MultipleInsertErrorContextKey struct{}

type MultipleObjects interface {
	ToInterfaceArray() []interface{}
}

// InsertObject - Insert the payload object to DB
func InsertMultipleObjects(collection string) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			os := r.Context().Value(ObjectContextKey{}).(MultipleObjects).ToInterfaceArray()
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			col := sess.Collection(collection)
			ids := make([]uuid.UUID, 0, len(os))
			ctx := r.Context()
			for _, o := range os {
				id, err := col.Insert(o)
				if err != nil {
					if s, ok := o.(Syncable); ok {
						if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
							idStruct := struct {
								ID uuid.UUID `db:"id"`
							}{}
							err := sess.Select("id").From(collection).Where("syncid = ?", s.GetSyncID()).One(&idStruct)
							if err != nil {
								logrus.Errorf("sess.Select in MultipleInsertObjects %q - %s %+v", err, collection, o)
								ctx = context.WithValue(ctx, MultipleInsertErrorContextKey{}, err)
								break
							}
							ids = append(ids, idStruct.ID)
							continue
						}
					}
					logrus.Errorf("Insert in MultipleInsertObjects %q - %s %+v", err, collection, o)
					ctx = context.WithValue(ctx, MultipleInsertErrorContextKey{}, err)
					break
				}
				ids = append(ids, uuid.FromStringOrNil(string(id.([]uint8))))
			}
			ctx = context.WithValue(r.Context(), MultipleInsertedIDsContextKey{}, ids)
			fn(w, r.WithContext(ctx), p)
		}
	}
}
