// Code generated; DO NOT EDIT.

package < .Package >

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// < .UpperHump >WithID is < .UpperHump > with ID
type < .UpperHump >WithID struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"< .CrudLower >_id"`
	< .UpperHump > `bson:",inline"`
}

// < .UpperHump >Service #path:"/< .CrudLower >/"#
type < .UpperHump >Service struct {
	db *mgo.Collection
}

// New< .UpperHump >Service Create a new < .UpperHump >Service
func New< .UpperHump >Service(db *mgo.Collection) (*< .UpperHump >Service, error) {
	return &< .UpperHump >Service{db}, nil
}

// Create a < .UpperHump > #route:"POST /"#
func (b *< .UpperHump >Service) Create(< .CrudLower > *< .UpperHump >) (err error) {
	return b.db.Insert(< .CrudLower >)
}

// Update the < .UpperHump > #route:"PUT /{< .CrudLower >_id}"#
func (s *< .UpperHump >Service) Update(< .CrudLower >_id bson.ObjectId, < .CrudLower > *< .UpperHump >) (err error) {
	return s.db.UpdateId(< .CrudLower >_id, < .CrudLower >)
}

// Delete the < .UpperHump > #route:"DELETE /{< .CrudLower >_id}"#
func (s *< .UpperHump >Service) Delete(< .CrudLower >_id bson.ObjectId) (err error) {
	return s.db.RemoveId(< .CrudLower >_id)
}

// Get the < .UpperHump > #route:"GET /{< .CrudLower >_id}"#
func (s *< .UpperHump >Service) Get(< .CrudLower >_id bson.ObjectId) (< .CrudLower > *< .UpperHump >, err error) {
	q := s.db.FindId(< .CrudLower >_id)
	err = q.One(&< .CrudLower >)
	if err != nil {
		return nil, err
	}
	return < .CrudLower >, nil
}

// List of the < .UpperHump > #route:"GET /"#
func (s *< .UpperHump >Service) List(offset, limit int) (< .CrudLower >s []*< .UpperHump >WithID, err error) {
	q := s.db.Find(nil).Skip(offset).Limit(limit)
	err = q.All(&< .CrudLower >s)
	if err != nil {
		return nil, err
	}
	return < .CrudLower >s, nil
}
