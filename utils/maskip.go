package utils

import "strings"

func MaskIP(ip string) string {
	strList := strings.Split(ip , "-")
	switch strings.Contains(strList[0], ":") {
	case true:
		ipGroup := strings.Split(strList[0], ":")
		if len(ipGroup) < 8 {
			var filler []string
			for i := 0; i < 8-len(ipGroup); i++{
				filler =append(filler,"0000")
			}
			for i := range ipGroup{
				if ipGroup[i] == "" {
					ipGroup = append(ipGroup[:i+1], append(filler, ipGroup[i+1:]...)...)
					break
				}
			}
		}
		ipGroup[4] = strings.Repeat("x", len(ipGroup[4]))
		ipGroup[5] = strings.Repeat("x", len(ipGroup[5]))
		strList[0] = strings.ReplaceAll(strings.Join(ipGroup, ":"), ":0000", "")

	default:
		ipGroup := strings.Split(strList[0], ".")
		ipGroup[2] = strings.Repeat("x", len(ipGroup[2]))
		strList[0] = strings.Join(ipGroup, ".")
	}
	return strings.Join(strList, "-")
}
