package profile

type Profile struct {
	ID     string `json:"-" bson:"_id,omitempty"`
	UserID string `json:"userId,omitempty" bson:"userId,omitempty"`
	Value  string `json:"-" bson:"value,omitempty"`

	FullName *string `json:"fullName" bson:"-"`
}
