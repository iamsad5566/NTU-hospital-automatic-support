package model

func ToDecimal(str string) int {
	multiple := getMultiple(len(str) - 1)
	decimalIndex := 0
	for len(str) > 0 {
		num := int(str[0]-'A') + 1
		decimalIndex += num * multiple
		str = str[1:]
		multiple /= 26
	}
	return decimalIndex - 1
}

func GetDecimalSlice(strColumns []string) map[int]bool {
	output := make(map[int]bool)
	for _, str := range strColumns {
		output[ToDecimal(str)] = true
	}
	return output
}

func getMultiple(digit int) int {
	output := 1
	for digit > 0 {
		output *= 26
		digit--
	}
	return output
}

func FilterTargetColumn(slice []string, chosen []string) string {
	mp := GetDecimalSlice(chosen)
	str := " ,"
	for i, cell := range slice {
		if mp[i] {
			str += cell + ","
		}
	}
	return str[:len(str)-1]
}

func SliceToStr(slice []string) string {
	output := ""
	for _, str := range slice {
		output += str + ","
	}
	return output[:len(output)-1]
}
