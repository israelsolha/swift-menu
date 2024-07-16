package main

import "fmt"

func Testing(i int) error {
	if i == 1 {
		return fmt.Errorf("Error")
	}
	// if i == 2 {
	// 	return fmt.Errorf("Error")
	// }
	// if i == 3 {
	// 	return fmt.Errorf("Error")
	// }
	return nil
}
