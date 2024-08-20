package flags

const (
	runAddressKey    = "a"
	dbURIKey         = "d"
	accrualSystemURL = "r"
	jwtSecret        = "s"
)

type Options struct {
	runAddress       string
	dbURI            string
	accrualSystemURL string
	jwtSecret        string
}

func (o *Options) HasRunAddress() bool {
	return o.runAddress != ""
}

func (o *Options) GetRunAddress() string {
	return o.runAddress
}

func (o *Options) HasDBUri() bool {
	return o.dbURI != ""
}

func (o *Options) GetDBUri() string {
	return o.dbURI
}

func (o *Options) HasAccrualSystemURL() bool {
	return o.accrualSystemURL != ""
}

func (o *Options) GetAccrualSystemURL() string {
	return o.accrualSystemURL
}

func (o *Options) HasJWTSecret() bool {
	return o.jwtSecret != ""
}

func (o *Options) GetJWTSecret() string {
	return o.jwtSecret
}
