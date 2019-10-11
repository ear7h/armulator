package is

import "strconv"

func immStr(n int) string {
	return "#" + strconv.Itoa(n)
}

func regStr(is64 bool, n byte) string {
	if n > 31 {
		return "INCORECT ENCODING FOR REGISTER: " + strconv.Itoa(int(n))
	}

	switch {
	case is64 && n == 31:
		return "SP"
	case is64:
		return "X" + strconv.Itoa(int(n))
	case !is64 && n == 31:
		return "WSP"
	case !is64:
		return "W" + strconv.Itoa(int(n))
	default:
		panic("unrechable")
	}
}
