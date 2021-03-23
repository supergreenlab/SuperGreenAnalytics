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

package db

import (
	"time"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID uuid.NullUUID `db:"id,omitempty" json:"id"`

	SyncID  uuid.NullUUID `db:"syncid" json:"syncid"`
	Skipped bool          `db:"-" json:"-"`

	SourceID  uuid.NullUUID `db:"sourceid" json:"sourceID"`
	SessionID uuid.NullUUID `db:"sessionid" json:"sessionID"`
	VisitorID uuid.NullUUID `db:"visitorid" json:"visitorID"`

	Type   string    `db:"etype" json:"type"`
	Params string    `db:"params" json:"params"`
	Date   time.Time `db:"createdat" json:"date"`

	CreatedAt time.Time `db:"cat,omitempty" json:"cat"`
	UpdatedAt time.Time `db:"uat,omitempty" json:"uat"`
}

func (e Event) GetSyncID() uuid.UUID {
	return e.SyncID.UUID
}

func (e *Event) SetIsSkipped() {
	e.Skipped = true
}

func (e Event) IsSkipped() bool {
	return e.Skipped
}

type Events struct {
	Events         []Event `json:"events"`
	interfaceCache []interface{}
	didCache       bool
}

func (es *Events) ToInterfaceArray() []interface{} {
	if es.didCache {
		return es.interfaceCache
	}
	es.interfaceCache = make([]interface{}, 0, len(es.Events))
	for _, e := range es.Events {
		es.interfaceCache = append(es.interfaceCache, &e)
	}
	es.didCache = true
	return es.interfaceCache
}
