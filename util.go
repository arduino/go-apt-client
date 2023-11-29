package apt

import "os"

func FileExists(path string) bool {
	st, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if st.IsDir() == true {
		return false
	}
	return true
}
