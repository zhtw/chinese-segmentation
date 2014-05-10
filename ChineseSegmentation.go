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
)

type TrieNode struct {
	children            map[rune]*TrieNode
	isValidSegmentation bool
}

type ChineseSegmentation struct {
	dict *TrieNode
}

func newTrieNode() (this *TrieNode) {
	this = new(TrieNode)
	this.children = make(map[rune]*TrieNode)
	return this
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
