package primitive

import "strconv"

// Identity
type Identity interface {
	Identity() string
	Integer() int64
}

func NewIdentity(v string) (Identity, error) {
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	return identity(n), nil
}

func CreateIdentity(v int64) Identity {
	return identity(v)
}

type identity int64

func (r identity) Identity() string {
	return strconv.FormatInt(int64(r), 10)
}

func (r identity) Integer() int64 {
	return int64(r)
}
