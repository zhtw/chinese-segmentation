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

func Test_New(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Cannot get current filename")
	}

	testData := filepath.Join(filepath.Dir(filename), "test", "data", "dict")

	_, err := New(testData)
	if err != nil {
		t.Fatal("Cannot create ChineseSegmentation")
	}
}

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
