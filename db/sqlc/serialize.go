package db

import "encoding/json"

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (fr *FriendRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(fr)
}

func (fr *FriendRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, fr)
}

func (f *Friendship) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Friendship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}

func (g *Group) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Group) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, g)
}

func (gm *GroupMember) MarshalBinary() ([]byte, error) {
	return json.Marshal(gm)
}

func (gm *GroupMember) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, gm)
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
