package algo

import "fmt"

func computeLPS(pat string) (lps []int) {
	M := len(pat)
	lps = make([]int, M) // lps[0] = 0

	l := 0 // length of the previous longest prefix suffix

	// the loop calculates lps[i] for i = 1 to M-1
	for i := 1; i < M; i++ {
		for {
			if pat[i] == pat[l] {
				l++
				break
			}

			if l == 0 {
				break
			}

			l = lps[l-1]
		}
		lps[i] = l
	}
	return lps
}

func KMPSearch(txt, pat string) (int, error) {
	M, N := len(pat), len(txt)

	// Preprocess the pattern that will hold the longest prefix suffix values for pattern
	lps := computeLPS(pat)

	for i, j := 0, 0; i < N; i++ {
		for {
			if pat[j] == txt[i] {
				j++

				if j == M {
					//	fmt.Printf("Found pattern at index %d \n", i-j+1)
					return i - j + 1, nil
					//j = lps[j-1]
				}
				break
			}

			if j > 0 {
				j = lps[j-1]
			} else {
				break
			}
		}
	}
	return 0, fmt.Errorf("pattern not found")
}
