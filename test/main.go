package main
import (
	"fmt"
)
func main()  {
	str := []string{"I","like","Golang"}

	i := 0


	for k, v := range str{
		str[1] = "66666"
		str = append(str, "good")
		fmt.Println(k, v)
	}

	fmt.Println(len(str))




	for i < len(str) {
		str[1] = "66666"
		str = append(str, "good")
		fmt.Println(i, str[i])
		i++
		if i == 10 {
			return
		}
	}


	for {
		str = append(str, "good")
		fmt.Println(str[i])
		i++
		if i == 10 {
			return
		}
	}

}