package utils

import "strconv"

type RGB struct {
	Red   int
	Green int
	Blue  int
}

func Hex2RGB(hex string) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   int(values >> 16),
		Green: int((values >> 8) & 0xFF),
		Blue:  int(values & 0xFF),
	}

	return rgb, nil
}
