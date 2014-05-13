/*
 * Copyright (c) 2014 ChangZhuo Chen <czchen@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 */
package chinesesegmentation

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"
)

type TrieNode struct {
	children            map[rune]*TrieNode
	isValidSegmentation bool
}

type ChineseSegmentation struct {
	dict *TrieNode
}

type Segmentation struct {
	start    int
	end      int
	isUnique bool
}

func newTrieNode() (this *TrieNode) {
	this = new(TrieNode)
	this.children = make(map[rune]*TrieNode)
	return this
}

func newSegmentation(start int, end int) (this Segmentation) {
	return Segmentation{start, end, false}
}

func New(dict string) (this *ChineseSegmentation, err error) {
	this = new(ChineseSegmentation)
	this.dict = newTrieNode()

	fd, err := os.Open(dict)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	for scanner := bufio.NewScanner(fd); scanner.Scan(); {
		text := scanner.Text()

		comment := strings.Index(text, "#")
		if comment != -1 {
			text = text[:comment]
		}

		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		head := this.dict
		for _, c := range text {
			_, ok := head.children[c]
			if !ok {
				head.children[c] = newTrieNode()
			}
			head = head.children[c]
		}
		head.isValidSegmentation = true
	}

	return this, nil
}

func getRuneArrayFromString(input string) (output []rune) {
	output = make([]rune, 0, utf8.RuneCountInString(input))

	for _, c := range input {
		output = append(output, c)
	}

	return output
}

func (this *ChineseSegmentation) getAllSegmentationFromRune(input []rune) (output []Segmentation) {
	output = make([]Segmentation, 0)

	for i := 0; i < len(input); i++ {
		output = append(output, newSegmentation(i, 1))

		curr, ok := this.dict.children[input[i]]
		if !ok {
			continue
		}

		for j := i + 1; j < len(input); j++ {
			curr, ok = curr.children[input[j]]
			if !ok {
				break
			}

			if curr.isValidSegmentation {
				output = append(output, newSegmentation(i, j-1))
			}
		}
	}

	return output
}

func isUniqueSegmentation(input []Segmentation, index int) bool {
	for i, item := range input {
		if i == index {
			continue
		}

		// contain
		if input[index].start <= item.start &&
			item.end <= input[index].end {
			continue
		}

		if input[index].end <= item.start {
			continue
		}

		if item.end <= input[index].start {
			continue
		}

		return false
	}
	return true
}

func removeUnusedSegmentation(input []Segmentation) (output []Segmentation) {
	removeFlags := make([]bool, len(input))

	/*
	 * Segmentation A contains B means the following conditions are all
	 * true:
	 *
	 *     A.start <= B.start
	 *     B.end <= A.end
	 *
	 * Segmentation A interleave B means one of the following conditions is
	 * true:
	 *
	 *     A.start <= B.start < A.end
	 *     B.start <= A.start < B.end
	 */

	/*
	 * For segmentation A, if all other segmentations are either contained
	 * by A, or are not onterleaved with A, segmentation A is called unique
	 * and all segmentation contained by A will be marked as remove.
	 */
	for i, _ := range input {
		if !isUniqueSegmentation(input, i) {
			continue
		}

		for j, _ := range input {
			if i != j &&
				input[i].start <= input[j].start &&
				input[j].end <= input[i].end {
				removeFlags[j] = true
			}
		}
	}

	output = make([]Segmentation, 0, len(input))

	for i, _ := range input {
		if !removeFlags[i] {
			output = append(output, input[i])
		}
	}

	return output
}

func (this *ChineseSegmentation) GetSegmentation(input string) (segmentation []string) {

	inputRune := getRuneArrayFromString(input)
	allSegs := this.getAllSegmentationFromRune(inputRune)
	_ = removeUnusedSegmentation(allSegs)

	return segmentation
}
