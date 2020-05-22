package utils
const PENDING = 0
const ACCEPTED = 1
const DECLINED = 2
const BlOCK = 3
const SUBSCRIBE = 4
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
