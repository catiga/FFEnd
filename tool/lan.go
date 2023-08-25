package tool

var (
	suppLans = [...]string{"zh-CN", "en"}
)

func IsSupportLan(lan string) bool {
	for _, v := range suppLans {
		if v == lan {
			return true
		}
	}
	return false
}
