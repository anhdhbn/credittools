package util

func Create_acc(code string, start int, end)([]string){
	// end: max is 5
	// start: 0
	var result []string
	year := time.Now().Year()
	month := time.Now().Month()
	if (month < 10) {
		year--
	}

	for i := start; i < end; i++ {
		head := year - i
		for j := 0; j < 2000; j++ {
			head_str := fmt.Sprintf("%v%s", head, code)[2:]
			acc := get_mssv(j, head_str)
			result = append(result, acc)
		}
	}
	return result
}

func get_mssv(i int, head string)(string) {
	if  (i < 10) {
		return fmt.Sprintf("%s000%v", head, i)
	} else if (i < 100) {
		return fmt.Sprintf("%s00%v", head, i)
	} else if (i < 1000) {
		return fmt.Sprintf("%s0%v", head, i)
	} else {
		return fmt.Sprintf("%s%v", head, i)
	}
}
