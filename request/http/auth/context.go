package auth

// IAuthContext 定义auth上下文接口
type IAuthContext interface {
	AuthClaimID() (ID uint, exist bool)
	IUserModel() IUser
}
