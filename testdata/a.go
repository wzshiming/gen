package a

import (
	"time"
)

//  UUID #format:"uuid"#
type UUID string

const (
	demoUUID UUID = "0000-0000-0000-0000"
)

type State uint8

const (
	_ State = 1 << iota
	A
	B
	C
	D
	E
)

type Kind string

const (
	K1 Kind = Kind("K1") + "888"
	K2 Kind = K3 + "K2"
	K3 Kind = K1 + "K3"
)

type Users []*User

type User struct {
	UUID     UUID      `json:"uuid,omitempty"`
	Name     string    `json:"name,omitempty"`
	PWD      string    `json:"pwd,omitempty"`
	CreateAt time.Time `json:"create_at,omitempty"`
	UserInfo
}

type UserInfo struct {
	Age   uint  `json:"age,omitempty"`
	State State `json:"state,omitempty"`
	Kind  Kind  `json:"kind,omitempty"`
}

// CreateUser #route:"PUT /user"#
func CreateUser(u *User, k Kind) (u0 *User, err error) {
	// TODO
	return nil, nil
}

// UpdateUser #route:"POST /user"#
func UpdateUser(u *User, hello string) (u0 *User, err error) {
	// TODO
	return nil, nil
}

// DeleteUser #route:"DELETE /user"#
func DeleteUser(userID uint64, b bool) (err error) {
	// TODO
	return nil
}

// GetUser #route:"GET /user/{userID}"#
func GetUser(userID uint64) (u *User, err error) {
	// TODO
	return nil, nil
}

// ListUser #route:"GET /user"#
func ListUser(offfset, limit uint64, state State) (us Users, err error) {
	// TODO
	return nil, nil
}
