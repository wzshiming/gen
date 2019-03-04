// Code generated; Do not regenerate the overwrite after editing.

package < .Package >

import (
	"errors"
)

// < .UpperHump >WithID is < .UpperHump > with ID
type < .UpperHump >WithID struct {
	ID int `json:"< .LowerSnake >_id,string"`
	< .Original >
}

// < .UpperHump >Service #path:"/< .LowerSnake >/"#
type < .UpperHump >Service struct {
	datas []*< .UpperHump >WithID
}

// New< .UpperHump >Service Create a new < .UpperHump >Service
func New< .UpperHump >Service() (*< .UpperHump >Service, error) {
	return &< .UpperHump >Service{}, nil
}

// Create a < .UpperHump > #route:"POST /"#
func (s *< .UpperHump >Service) Create(< .LowerHump > *< .Original >) (err error) {
	< .LowerHump >ID := len(s.datas) + 1
	data := &< .UpperHump >WithID {
		ID: < .LowerHump >ID ,
		< .Original >: *< .LowerHump >,
	}
	s.datas = append(s.datas, data)
	return nil
}

// Update the < .UpperHump > #route:"PUT /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Update(< .LowerHump >ID int /* #name:"< .LowerSnake >_id"# */, < .LowerHump > *< .Original >) (err error) {
	if 0 >= < .LowerHump >ID || < .LowerHump >ID > len(s.datas) || s.datas[< .LowerHump >ID-1] == nil {
		return errors.New("id does not exist")
	}
	v := s.datas[< .LowerHump >ID-1]
	v.< .Original > = *< .LowerHump >
	return nil
}

// Delete the < .UpperHump > #route:"DELETE /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Delete(< .LowerHump >ID int /* #name:"< .LowerSnake >_id"# */) (err error) {
	if 0 >= < .LowerHump >ID || < .LowerHump >ID > len(s.datas) || s.datas[< .LowerHump >ID-1] == nil {
		return errors.New("id does not exist")
	}
	s.datas[< .LowerHump >ID-1] = nil
	return nil
}

// Get the < .UpperHump > #route:"GET /{< .LowerSnake >_id}"#
func (s *< .UpperHump >Service) Get(< .LowerHump >ID int /* #name:"< .LowerSnake >_id"# */) (< .LowerHump > *< .UpperHump >WithID, err error) {
	if 0 >= < .LowerHump >ID || < .LowerHump >ID > len(s.datas) || s.datas[< .LowerHump >ID-1] == nil {
		return nil, errors.New("id does not exist")
	}
	return s.datas[< .LowerHump >ID-1], nil
}

// List of the < .UpperHump > #route:"GET /"#
func (s *< .UpperHump >Service) List(offset, limit int) (< .LowerHump >s []*< .UpperHump >WithID, err error) {
	off := 0
	lim := 0
	for _, v := range s.datas {
		if v != nil {
			if offset != 0 &&  off != offset {
				off++
				continue
			}
			if limit == 0 || lim != limit{
				lim++
				< .LowerHump >s = append(< .LowerHump >s, v)
				if lim == limit {
					break
				}
				continue
			}
		}
	}
	return < .LowerHump >s, nil
}
