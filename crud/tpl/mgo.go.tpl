// Code generated; Do not regenerate the overwrite after editing.

package < .Package >

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// < .UpperHump >WithID is < .UpperHump > with ID
type < .UpperHump >WithID struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"< .LowerSnake >_id"`
	< .Original > `bson:",inline"`
}

// < .UpperHump >Service #path:"/< .LowerSnake >/"#
type < .UpperHump >Service struct {
	db *mgo.Collection
}

// New< .UpperHump >Service Create a new < .UpperHump >Service
func New< .UpperHump >Service(db *mgo.Collection) (*< .UpperHump >Service, error) {
	return &< .UpperHump >Service{db}, nil
}

// Create a < .UpperHump > #route:"POST /"#
func (s *< .UpperHump >Service) Create(< .LowerHump > *< .Original >) (err error) {
	return  s.db.Insert(< .LowerHump >)
}

// Update the < .UpperHump > #route:"PUT /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Update(< .LowerHump >ID bson.ObjectId /* #name:"< .LowerSnake >_id"# */, < .LowerHump > *< .Original >) (err error) {
	return s.db.UpdateId(< .LowerHump >ID, bson.D{{"$set", < .LowerHump >}})
}

// Delete the < .UpperHump > #route:"DELETE /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Delete(< .LowerHump >ID bson.ObjectId /* #name:"< .LowerSnake >_id"# */) (err error) {
	return s.db.RemoveId(< .LowerHump >ID)
}

// Get the < .UpperHump > #route:"GET /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Get(< .LowerHump >ID bson.ObjectId /* #name:"< .LowerSnake >_id"# */) (< .LowerHump > *< .UpperHump >WithID, err error) {
	q := s.db.FindId(< .LowerHump >ID)
	err = q.One(&< .LowerHump >)
	if err != nil {
		return nil, err
	}
	return < .LowerHump >, nil
}

// List of the < .UpperHump > #route:"GET /"#
func (s *< .UpperHump >Service) List(offset, limit int) (< .LowerHump >s []*< .UpperHump >WithID, err error) {
	q := s.db.Find(nil).Skip(offset).Limit(limit)
	err = q.All(&< .LowerHump >s)
	if err != nil {
		return nil, err
	}
	return < .LowerHump >s, nil
}
