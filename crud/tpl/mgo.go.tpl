// Code generated; DO NOT EDIT.

package < .Package >

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// < .CrudUpper >WithID is < .CrudUpper > with ID
type < .CrudUpper >WithID struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"< .CrudLower >_id"`
	< .CrudUpper > `bson:",inline"`
}

// < .CrudUpper >Service #path:"/< .CrudLower >/"#
type < .CrudUpper >Service struct {
	db *mgo.Collection
}

// New< .CrudUpper >Service Create a new < .CrudUpper >Service
func New< .CrudUpper >Service(db *mgo.Collection) (*< .CrudUpper >Service, error) {
	return &< .CrudUpper >Service{db}, nil
}

// Create a < .CrudUpper > #route:"POST /"#
func (b *< .CrudUpper >Service) Create(< .CrudLower > *< .CrudUpper >) (err error) {
	return b.db.Insert(< .CrudLower >)
}

// Update the < .CrudUpper > #route:"PUT /{< .CrudLower >_id}"#
func (s *< .CrudUpper >Service) Update(< .CrudLower >_id string, < .CrudLower > *< .CrudUpper >) (err error) {
	return s.db.UpdateId(< .CrudLower >_id, < .CrudLower >)
}

// Delete the < .CrudUpper > #route:"DELETE /{< .CrudLower >_id}"#
func (s *< .CrudUpper >Service) Delete(< .CrudLower >_id string) (err error) {
	return s.db.RemoveId(< .CrudLower >_id)
}

// Get the < .CrudUpper > #route:"GET /{< .CrudLower >_id}"#
func (s *< .CrudUpper >Service) Get(< .CrudLower >_id string) (< .CrudLower > *< .CrudUpper >, err error) {
	q := s.db.FindId(< .CrudLower >_id)
	err = q.One(&< .CrudLower >)
	if err != nil {
		return nil, err
	}
	return < .CrudLower >, nil
}

// List of the < .CrudUpper > #route:"GET /"#
func (s *< .CrudUpper >Service) List(offset, limit int) (< .CrudLower >s []*< .CrudUpper >WithID, err error) {
	q := s.db.Find(nil).Skip(offset).Limit(limit)
	err = q.All(&< .CrudLower >s)
	if err != nil {
		return nil, err
	}
	return < .CrudLower >s, nil
}
