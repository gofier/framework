package auth

type IUser interface {
	Scan(userID uint) error
	Value() interface{}
}
