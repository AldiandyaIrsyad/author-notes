package auth

type User struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	Email     string `bson:"email" json:"email"`
	Username  string `bson:"username" json:"username"`
	Password  string `bson:"password" json:"-"` // hashed password
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`
}

type AuthToken struct {
	Token string `json:"token"`
}
