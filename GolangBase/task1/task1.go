package main

import (
	"fmt"
	"strconv"
)

// 136. 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。
// 找出那个只出现了一次的元素。可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，
// 例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素。
func pickOne(nums []int) int {
	var resultMap = make(map[int]int)
	for _, num := range nums {
		resultMap[num]++
	}
	for num, count := range resultMap {
		if count == 1 {
			return num // 返回只出现一次的数字
		}
	}
	return -1 // 如果没有找到，返回一个特殊值
}

// 判断一个整数是否是回文数
func palindromic(num int) bool {
	var str = strconv.Itoa(num)
	var result = true
	for i := 0; i < len(str)/2; i++ {
		if str[i] != str[len(str)-i-1] {
			result = false
			break
		}
	}
	return result
}

func main() {
	// 1.只出现一次的数字
	var nums = []int{2, 2, 1, 3, 4, 4, 3}
	var result = pickOne(nums)
	fmt.Println(result)

	// 2.判断一个整数是否是回文数
	var num = 12321
	var result2 = palindromic(num)
	fmt.Println(result2)
}
