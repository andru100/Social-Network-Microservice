package model


type Likes struct {
	Username string `bson:"Username" json:"Username"`
	Profpic  string `bson:"Profpic" json:"Profpic"`
}

type MongoFields struct {
	Key            string     `json:"Key"`
	ID 			   string 	  `bson:"_id" json:"ID"`
	Username       string     `json:"Username"`
	Password       string     `json:"Password"`
	Email          string     `json:"Email"`
	Bio            string     `json:"Bio"`
	Profpic        string     `json:"Profpic"`
	Photos         []string   `json:"Photos"`
	LastCommentNum int32      `json:"LastCommentNum"`
	Posts          []PostData `json:"Posts"`
}

type MsgCmts struct {
	Username string `bson:"Username" json:"Username"`
	Comment  string `bson:"Comment" json:"Comment"`
	Profpic  string `bson:"Profpic" json:"Profpic"`
}

type PostData struct {
	Username    string    `bson:"Username" json:"Username"`
	SessionUser string    `bson:"SessionUser" json:"SessionUser"`
	MainCmt     string    `bson:"MainCmt" json:"MainCmt"`
	PostNum     int32     `bson:"PostNum" json:"PostNum"`
	Time        string    `bson:"Time" json:"Time"`
	TimeStamp   int64     `bson:"TimeStamp" json:"TimeStamp"`
	Date        string    `bson:"Date" json:"Date"`
	Comments    []MsgCmts `bson:"Comments" json:"Comments"`
	Likes       []Likes   `bson:"Likes" json:"Likes"`
}
