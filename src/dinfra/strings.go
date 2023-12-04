package dinfra

func Strings_FirstToLower(str string) string {
	s := []rune(str)
	if len(s) == 0 {
		return str
	}
	if s[0] >= 97 && s[0] <= 122 {
		return str
	}
	if s[0] >= 65 && s[0] <= 90 {
		s[0] = s[0] + 32
	}
	return string(s)
}

func Strings_FirstToUpper(str string) string {
	s := []rune(str)
	if len(s) == 0 {
		return str
	}
	if s[0] >= 65 && s[0] <= 90 {
		return str
	}
	if s[0] >= 97 && s[0] <= 122 {
		s[0] = s[0] - 32
	}
	return string(s)
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}
