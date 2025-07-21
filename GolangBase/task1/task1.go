package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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

/**
* 编写一个函数来查找字符串数组中的最长公共前缀。
* 如果不存在公共前缀，返回空字符串 ""。
 */
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, str := range strs {
		for strings.Index(str, prefix) != 0 {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

/*
*
* 给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。
* 这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。

* 将这个大整数加 1，并返回结果的数字数组。
 */
func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	return append([]int{1}, digits...)
}

/**
* 给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，返回删除后数组的新长度。
* 元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。
 */
func removeDuplicates(nums []int) int {
	var result = 0
	for i := 0; i < len(nums); i++ {
		if i == 0 || nums[i] != nums[i-1] {
			nums[result] = nums[i]
			result++
		}
	}
	return result
}

/**
* 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
* 请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。
 */
func merge(intervals [][]int) [][]int {
	var result = make([][]int, 0)
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	for _, interval := range intervals {
		if len(result) == 0 || result[len(result)-1][1] < interval[0] {
			result = append(result, interval)
		}
		result[len(result)-1][1] = max(result[len(result)-1][1], interval[1])
	}
	return result
}

/**
* 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数，并返回它们的数组下标。
* 你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。
* 你可以按任意顺序返回答案。
 */	
func twoSum(nums []int, target int) []int {
	var result = make([]int, 0)
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				result = append(result, i, j)
			}
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

	// 3. 最长公共前缀
	var strs = []string{"flower", "flow", "flight"}
	var result3 = longestCommonPrefix(strs)
	fmt.Println(result3)

	// 4.数组操作、进位处理
	var digits = []int{9, 9, 9}
	var result4 = plusOne(digits)
	fmt.Println(result4)

	// 5. 合并区间
	var intervals = [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	var result5 = merge(intervals)
	fmt.Println(result5)

	// 6. 两数之和
	var nums2 = []int{2, 7, 11, 15}
	var target = 9
	var result6 = twoSum(nums2, target)
	fmt.Println(result6)
}
