package replacer

import (
	"fmt"
	"strings"

	"github.com/oarkflow/pkg/docx/replacer/algo"
)

func (d *Document) RemoveBlockByID(i int) error {
	if len(d.WP) < i {
		return fmt.Errorf("remove block failed. len []wp < i")

	}
	part1 := d.WP[:i]
	part2 := d.WP[i+1:]
	var r []WP
	r = append(r, part1...)
	r = append(r, part2...)
	d.WP = r
	return nil
}
func (d *Document) InsertBlockAfterBlockByID(i int, wp WP) error {
	if len(d.WP) < i {
		return fmt.Errorf("insert block failed. len []wp < i")

	}
	part1 := d.WP[:i]
	part2 := d.WP[i:]
	var r []WP
	r = append(r, part1...)
	r = append(r, wp)
	r = append(r, part2...)
	d.WP = r
	return nil
}
func (d *Document) GetCopyBlockByTag(pattern string) (WP, error) {
	id, err := d.GetBlockIDByTag(pattern)
	if err != nil {
		return WP{}, err
	}
	var wp WP
	CopyWithOption(&wp, &d.WP[id], Option{IgnoreEmpty: false, DeepCopy: true})
	arr, err := ExtractWPToArrayTextString(wp)
	if err != nil {
		return WP{}, err
	}
	rebArr, _, err := RebuildBlocks(pattern, arr)
	if err != nil {
		return WP{}, err
	}

	wpNew, err := BuildArrayTextStringToWP(wp, rebArr)
	if err != nil {
		return WP{}, err
	}

	return wpNew, nil
}
func (d *Document) ReplaceTextByTag(pattern string, text string) error {
	id, err := d.GetBlockIDByTag(pattern)
	if err != nil {
		return err
	}
	wp := d.WP[id]
	arr, err := ExtractWPToArrayTextString(wp)
	if err != nil {
		return err
	}
	rebArr, patID, err := RebuildBlocks(pattern, arr)
	if err != nil {
		return err
	}
	rebArr[patID] = strings.Replace(rebArr[patID], pattern, Screening(text), 1)
	wpNew, err := BuildArrayTextStringToWP(wp, rebArr)
	if err != nil {
		return err
	}
	wp = wpNew
	return nil
}

func (d *Document) GetBlockIDByTag(tag string) (int, error) {
	for i, WPItem := range d.WP { // итерация всех параграфов
		var body string
		for _, token := range WPItem.Body {
			if token.Tag == "w:r" {
				//	fmt.Printf("i %d z %d\n", i, z)
				res, err := wpParser(token.Body)
				if err != nil {
					return 0, err
				}
				for _, wtTag := range res {
					if wtTag.Tag == "w:t" {
						body += wtTag.Body
					}
				}
			}
		}
		// поиск, есть ли в параграфе паттерн
		if strings.Contains(body, tag) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("tag not found")
}

// RebuildBlocks (pattern string, source []string) (expectedArray []string, blockIDWithTag int, err)
func RebuildBlocks(pattern string, source []string) ([]string, int, error) {
	start, stop, shift, err := algo.FindMatchInArray(source, pattern)
	if err != nil {
		return []string{}, 0, err
	}
	if start == stop { // если паттерн находится только в одном элементе
		return source, start, nil
	}
	// ______________________________
	// fmt.Printf("start: %d| stop: %d| shift: %d\n", start, stop, shift)
	// variables block_______
	var newSource []string
	var firstElem string
	patternLen := len(pattern)
	patternLenCounter := 0
	// closed[variables block]
	for i, item := range source {
		if start <= i && i <= stop {
			// блок если в item находится часть pattern
			if stop == i {
				// работа с последним блоком, в котором есть паттерн
				shift := patternLen - patternLenCounter
				firstElem += item[:shift]
				newSource = append(newSource, item[shift:])
				continue
			}

			if i == start {
				// работа с первым блоком, в котором есть паттерн
				if shift == 0 {
					firstElem += item
					patternLenCounter += len(item)
					newSource = append(newSource, "")
				} else {
					firstElem += item[shift:]
					patternLenCounter += len(item[shift:])
					newSource = append(newSource, item[:shift])
				}
				continue
				// closed[работа с первым блоком, в котором есть паттерн]
			}
			// ____работа с средними блоками, в которых есть паттерн
			firstElem += item
			patternLenCounter += len(item)
			newSource = append(newSource, "")
			// ____closed[работа с средними блоками, в которых есть паттерн]
			// closed[блок если в item находится часть pattern]
		} else {
			newSource = append(newSource, item) // сохранение остальных беспаттерновых блоков без изменений
		}
	}
	newSource[start] += firstElem
	//	fmt.Printf("newSource: [%+v] patternLen: [%d], plCounter: [%d]\n", newSource, patternLen, patternLenCounter)
	return newSource, start, nil
}

func ExtractWPToArrayTextString(wp WP) ([]string, error) {
	var body []string
	for _, item := range wp.Body {
		switch item.Tag { // извлекает из параграфа только блоки с текстом
		case "w:r":
			// работа с тегом w:r________
			wrTokens, err := wpParser(item.Body) // парсит содержимое <w:r> на отдельные токены
			if err != nil {
				return []string{}, err
			}
			// ____работа с токенами тега "w:r" и их фильтрация
			for _, wrToken := range wrTokens {
				switch wrToken.Tag {
				case "w:t": // фильтруем только токены w:t. в них содержится текст
					//	fmt.Printf("w:t: [%s]\n", wrToken.Body)
					body = append(body, wrToken.Body)
				}
			}
			// ____
			// closed[работа с тегом w:r]
		}
	}
	return body, nil
}
func BuildArrayTextStringToWP(wp WP, bodyStrings []string) (WP, error) {
	counter := 0
	for r, item := range wp.Body {
		switch item.Tag { // извлекает из параграфа только блоки с текстом
		case "w:r":
			// работа с тегом w:r________
			wrTokens, err := wpParser(item.Body) // парсит содержимое <w:r> на отдельные токены
			if err != nil {
				return WP{}, err
			}
			var wrBody string // заготовка тела wr для упаковки
			// ____работа с токенами тега "w:r" и их фильтрация
			for i, wrToken := range wrTokens {
				switch wrToken.Tag {
				case "w:t": // фильтруем только токены w:t. в них содержится текст

					if len(bodyStrings) < counter+1 {
						return WP{}, fmt.Errorf("len(bodyStrings)mpore than len([]array(<w:t>))")
					}
					wrToken.Body = bodyStrings[counter]
					wrTokens[i] = wrToken
					counter++

				}
				wrBody += AtomicWPTokensToString(wrToken)
			}
			// ____closed[работа с токенами тега "w:r" и их фильтрация]
			// ____упаковка body токена w:r
			wp.Body[r].Body = wrBody
			// ____closed[упаковка body токена w:r]
			// closed[работа с тегом w:r]
		}
	}
	return wp, nil
}
