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
	"path/filepath"
	"runtime"
	"testing"
)

func Test_getRuneArrayFromString(t *testing.T) {
	output := getRuneArrayFromString("測試")

	if len(output) != 2 {
		t.Fatal("len(output) shall be 2")
	}

	if output[0] != 28204 {
		t.Errorf("output[0] = %d shall be 28204", output[0])
	}

	if output[1] != 35430 {
		t.Errorf("output[0] = %d shall be 35430", output[1])
	}
}

func compareSegmentationRange(x []Segmentation, y []Segmentation) bool {
	if len(x) != len(y) {
		return false
	}

	for i, _ := range x {
		if x[i].start != y[i].start {
			return false
		}

		if x[i].end != y[i].end {
			return false
		}
	}

	return true
}

func Test_getAllSegmentationFromRune(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Cannot get current filename")
	}

	testData := filepath.Join(filepath.Dir(filename), "test", "data", "dict")

	this, err := New(testData)
	if err != nil {
		t.Fatal("Cannot create ChineseSegmentation")
	}

	res := this.getAllSegmentationFromRune(getRuneArrayFromString("自由和平等"))

	expected := []Segmentation{
		{0, 1},
		{0, 2}, // 自由
		{1, 2},
		{2, 3},
		{2, 4}, // 和平
		{3, 4},
		{3, 5}, // 平等
		{4, 5},
	}

	if compareSegmentationRange(res, expected) {
		t.Fatal("res is not expected value")
	}
}

func Test_isUniqueSegmentation(t *testing.T) {
	input := []Segmentation{
		{0, 1}, // removed
		{0, 2}, // unique
		{1, 2}, // removed
		{2, 3}, // keep
		{2, 4}, // not unique
		{3, 4}, // removed
		{3, 5}, // not unique
		{4, 5}, // keep
	}

	if !isUniqueSegmentation(input, 1) {
		t.Error("input[1] shall be unique")
	}

	if isUniqueSegmentation(input, 4) {
		t.Error("input[4] shall not be unique")
	}

	if isUniqueSegmentation(input, 6) {
		t.Error("input[6] shall not be unique")
	}
}

func Test_removeUnusedSegmentation(t *testing.T) {
	input := []Segmentation{
		{0, 1}, // removed
		{0, 2}, // unique
		{1, 2}, // removed
		{2, 3}, // keep
		{2, 4}, // not unique
		{3, 4}, // removed
		{3, 5}, // not unique
		{4, 5}, // keep
	}

	expected := []Segmentation{
		{0, 2}, // unique
		{2, 3}, // keep
		{2, 4}, // not unique
		{3, 4}, // removed
		{3, 5}, // not unique
		{4, 5}, // keep
	}

	output := removeUnusedSegmentation(input)
	if !compareSegmentationRange(output, expected) {
		t.Error("output is not expected value")
	}
}
