package enums

type AuthOperation string

const (
	AuthSignUp        AuthOperation = "SignUp"
	AuthSignIn        AuthOperation = "SignIn"
	AuthSignOut       AuthOperation = "SignOut"
	AuthResetPassword AuthOperation = "ResetPassword"
	AuthRequestReset  AuthOperation = "RequestReset"
)
