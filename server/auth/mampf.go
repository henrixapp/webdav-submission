//Communication with the mampf server
package auth

type MaMpfClient interface {
	//Validate checks if a user is a valid MaMpf user
	Validate(username, password string) (User, error)
}

//User is a common User
type User interface {
}

type MampfParams struct {
	APIURL string
	Token  string
}

type MaMpfClientImpl struct{}
type UserImpl struct{}

func (m MaMpfClientImpl) Validate(username, password string) (User, error) {
	return UserImpl{}, nil
}
