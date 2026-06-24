package model

type User struct {
	ID       uint64
	Name     string
	Email    string
	Password string
}

type TaskCreator struct {
	UserID    uint64
	UserName  string
	TeamID    uint64
	Name      string
	TaskCount int
	Rank      int
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
