package algo

func grayCode(n int) []int {
	ans := []int{0}
	for i := 0; i < n; i++ {
		for j := len(ans) - 1; j >= 0; j-- {
			ans = append(ans, ans[j]|1<<i)
		}
	}
	return ans
}
