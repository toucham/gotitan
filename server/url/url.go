package url

import "fmt"

type Url struct {
	path string
	url  string
}

func (u *Url) String() string {
	return fmt.Sprintf("%s/%s", u.url, u.path)
}
