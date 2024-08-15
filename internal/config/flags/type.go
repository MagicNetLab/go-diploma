package flags

const (
	runAddressKey    = "a"
	dbUriKey         = "d"
	accrualSystemUrl = "r"
)

type Options struct {
	runAddress       string
	dbUri            string
	accrualSystemUrl string
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
