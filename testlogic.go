package main

import "fmt"

func main() {
	var temp = make(map[string]interface{})
	temp["total"] = 120
	temp["jan"] = 1
	temp["feb"] = 2
	temp["mar"] = 3
	temp["apr"] = 4
	temp["mei"] = 5
	temp["jun"] = 6
	temp["jul"] = 7
	temp["aug"] = 8
	temp["sep"] = 9
	temp["okt"] = 10
	temp["nov"] = 11
	temp["des"] = 12

	var monthStart, monthend int
	//var monthTemp []string
	monthStart = 6
	monthend = 12



	for key, _ := range temp {
		var exist bool
		if key == "total"{
			continue
		}
		for i := monthStart; i <= monthend; i++ {
			switch i {
			case 1 :
				if key == "jan" {
					exist = true
					break
				}
			case 2 :
				if key == "feb" {
					exist = true
					break
				}
			case 3 :
				if key == "mar" {
					exist = true
					break
				}
			case 4 :
				if key == "apr" {
					exist = true
					break
				}
			case 5 :
				if key == "mei" {
					exist = true
					break
				}
			case 6 :
				if key == "jun" {
					exist = true
					break
				}
			case 7 :
				if key == "jul" {
					exist = true
					break
				}
			case 8 :
				if key == "aug" {
					exist = true
					break
				}
			case 9 :
				if key == "sep" {
					exist = true
					break
				}
			case 10 :
				if key == "okt" {
					exist = true
					break
				}
			case 11:
				if key == "nov" {
					exist = true
					break
				}
			case 12 :
				if key == "des" {
					exist = true
					break
				}

			}
		}
		if !exist {
			//todo delete map key
			delete (temp, key)
		}
	}


	fmt.Println(temp)
}

