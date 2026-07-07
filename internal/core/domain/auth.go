package domain

type RegisterUserPayload struct {
	Username  string
	FirstName string
	LastName  *string
	Bio       *string
	Password  string
}

func NewRegisterUserPayload(
	username string,
	firstName string,
	lastName *string,
	bio *string,
	password string,
) RegisterUserPayload {
	return RegisterUserPayload{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Bio:       bio,
		Password:  password,
	}
}

func (p *RegisterUserPayload) Validate() error {
	if err := ValidateUsername(p.Username); err != nil {
		return err
	}

	if err := ValidateFirstName(p.FirstName); err != nil {
		return err
	}

	if p.LastName != nil {
		if err := ValidateLastName(*p.LastName); err != nil {
			return err
		}
	}

	if p.Bio != nil {
		if err := ValidateBio(*p.Bio); err != nil {
			return err
		}
	}

	if err := ValidatePassword(p.Password); err != nil {
		return err
	}

	return nil
}
