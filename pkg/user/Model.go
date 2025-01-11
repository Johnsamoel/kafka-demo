package user

// User represents a user document in MongoDB
type User struct {
	UserId     string `bson:"userId,omitempty" json:"userId"`
	FirstName  string `bson:"first_name" json:"first_name" validate:"required,min=2,max=50"`
	LastName   string `bson:"last_name" json:"last_name" validate:"required,min=2,max=50"`
	Email      string `bson:"email" json:"email" validate:"required,email"`
	Password   string `bson:"password,omitempty" json:"-" validate:"required,min=8"`
	IsNotified bool   `bson:"isNotified" json:"isNotified"`
}
