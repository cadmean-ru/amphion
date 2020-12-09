package common

type Platform struct {
	pString string
	name    string
	os      string
}

func (p Platform) GetName() string {
	return p.name
}

func (p Platform) ToString() string {
	return p.pString
}

func PlatformFromString(pString string) Platform {
	return Platform{
		name: pString,
		os: "pipi",
	}
}