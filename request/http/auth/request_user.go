package auth

type RequestIUser interface {
	ScanUser() error
	User() IUser
	UserId() (userID uint, err error)
	ScanUserWithJSON() (isAbort bool)
}

type UserNotLoginError struct {
}

func (e UserNotLoginError) Error() string {
	return "user not login"
}

type UserNotExistError struct {
}

func (e UserNotExistError) Error() string {
	return "user not exists"
}

type RequestUser struct {
	c    IAuthContext
	user IUser
}

func (au *RequestUser) User() IUser {
	return au.user
}

func (au *RequestUser) UserId() (userID uint, err error) {
	exist := false
	userID, exist = au.c.AuthClaimID()
	if !exist {
		return 0, UserNotLoginError{}
	}
	return userID, nil
}

func (au *RequestUser) SetContext(c IAuthContext) {
	au.c = c
}

func (au *RequestUser) ScanUser() error {
	if au.user != nil {
		return nil
	}
	user := au.c.IUserModel()
	userID, err := au.UserId()
	if err != nil {
		return err
	}
	if err := user.Scan(userID); err != nil {
		return UserNotExistError{}
	}
	au.user = user
	return nil
}
