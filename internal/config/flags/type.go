package flags

const (
	runAddressKey    = "a"
	dbUriKey         = "d"
	accrualSystemUrl = "r"
	jwtSecret        = "s"
)

type Options struct {
	runAddress       string
	dbUri            string
	accrualSystemUrl string
	jwtSecret        string
}

func (o *Options) HasRunAddress() bool {
	return o.runAddress != ""
}

func (o *Options) GetRunAddress() string {
	return o.runAddress
}

func (o *Options) HasDBUri() bool {
	return o.dbUri != ""
}

func (o *Options) GetDBUri() string {
	return o.dbUri
}

func (o *Options) HasAccrualSystemUrl() bool {
	return o.accrualSystemUrl != ""
}

func (o *Options) GetAccrualSystemUrl() string {
	return o.accrualSystemUrl
}

func (o *Options) HasJWTSecret() bool {
	return o.jwtSecret != ""
}

func (o *Options) GetJWTSecret() string {
	return o.jwtSecret
}
