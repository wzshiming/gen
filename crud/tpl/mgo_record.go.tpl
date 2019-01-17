// Code generated; Do not regenerate the overwrite after editing.

package < .Package >

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// < .UpperHump >WithID is < .UpperHump > with ID
type < .UpperHump >WithID struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"< .LowerSnake >_id"`
	< .UpperHump >   `bson:",inline"`
	Update time.Time `bson:"update,omitempty" json:"update"`
}

type < .UpperHump >Record struct {
	ID                bson.ObjectId   `bson:"_id,omitempty" json:"< .LowerSnake >_record_id"`
	< .UpperHump >ID  bson.ObjectId   `bson:"< .LowerSnake >_id,omitempty" json:"< .LowerSnake >_id"`
	Recent            *< .UpperHump > `bson:"recent,omitempty" json:"recent"`
	Current           *< .UpperHump > `bson:"current,omitempty" json:"current"`
	RecentTime        time.Time       `bson:"recent_time,omitempty" json:"recent_time"`
	CurrentTime       time.Time       `bson:"current_time,omitempty" json:"current_time"`
	Times             int             `bson:"times,omitempty" json:"times"`
}

// < .UpperHump >Service #path:"/< .LowerSnake >/"#
type < .UpperHump >Service struct {
	db       *mgo.Collection
	dbRecord *mgo.Collection
}

// New< .UpperHump >Service Create a new < .UpperHump >Service
func New< .UpperHump >Service(db *mgo.Collection) (*< .UpperHump >Service, error) {
	dbRecord := db.Database.C(db.Name + "_record")
	dbRecord.EnsureIndex(mgo.Index{Key: []string{"< .LowerSnake >_id"}})
	return &< .UpperHump >Service{
		db:       db,
		dbRecord: dbRecord,
	}, nil
}

// Create a < .UpperHump > #route:"POST /"#
func (s *< .UpperHump >Service) Create(< .LowerHump > *< .UpperHump >) (err error) {
	return s.db.Insert(&< .UpperHump >WithID{
		< .UpperHump >:   *< .LowerHump >,
		Update: bson.Now(),
	})
}

// Update the < .UpperHump > #route:"PUT /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Update(< .LowerHump >ID bson.ObjectId /* #name:"< .LowerSnake >_id"# */, < .LowerHump > *< .UpperHump >) (err error) {
	if err := s.record(< .LowerHump >ID, < .LowerHump >); err != nil {
		return err
	}

	return s.db.UpdateId(< .LowerHump >ID, bson.D{{"$set", &< .UpperHump >WithID{
		< .UpperHump >:   *< .LowerHump >,
		Update: bson.Now(),
	}}})
}

// Delete the < .UpperHump > #route:"DELETE /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Delete(< .LowerHump >ID bson.ObjectId /* #name:"< .LowerSnake >_id"# */) (err error) {
	if err := s.record(< .LowerHump >ID, nil); err != nil {
		return err
	}

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

// RecordList of the < .UpperHump > record list #route:"GET /record/{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) RecordList(< .LowerHump >ID bson.ObjectId /* #name:"< .LowerSnake >_id"# */, offset, limit int) (< .LowerHump >Records []*< .UpperHump >Record, err error) {
	q := s.dbRecord.Find(bson.D{{"< .LowerSnake >_id", < .LowerHump >ID}}).Skip(offset).Limit(limit)
	err = q.All(&< .LowerHump >Records)
	if err != nil {
		return nil, err
	}
	return < .LowerHump >Records, nil
}

func (s *< .UpperHump >Service) record(< .LowerHump >ID bson.ObjectId, current *< .UpperHump >) error {
	v, err := s.Get(< .LowerHump >ID)
	if err != nil {
		return err
	}
	count, err := s.dbRecord.Find(bson.D{{"< .LowerSnake >_id", v.ID}}).Count()
	if err != nil {
		return err
	}
	return s.dbRecord.Insert(&< .UpperHump >Record{
		< .UpperHump >ID:      v.ID,
		Recent:      &v.< .UpperHump >,
		Current:     current,
		RecentTime:  v.Update,
		CurrentTime: bson.Now(),
		Times:       count + 1,
	})
}

