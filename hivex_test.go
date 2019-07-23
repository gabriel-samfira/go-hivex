// Copyright 2019 Cloudbase Solutions SRL
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hivex

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func openHive(name string, t *testing.T) *Hivex {
	pth := filepath.Join("testdata", name)
	if _, err := os.Stat(pth); err != nil {
		panic(fmt.Sprintf("missing test hive %s", name))
	}
	hive, err := NewHivex(pth, WRITE)
	if err != nil {
		t.Errorf("Failed to open hive: %q\n", err)
	}
	return hive
}

func closeHive(h *Hivex, t *testing.T) {
	err := h.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpenClose(t *testing.T) {
	hive := openHive("minimal", t)
	closeHive(hive, t)
}

func TestRLenValue(t *testing.T) {
	hive := openHive("rlenvalue_test_hive", t)
	root, err := hive.Root()
	assert.NilError(t, err)
	child, err := hive.NodeGetChild(root, "ModerateValueParent")
	assert.NilError(t, err)
	assert.Assert(t, child != 0)

	moderateValue, err := hive.NodeGetValue(child, "33Bytes")
	assert.NilError(t, err)
	length, offset, err := hive.NodeValueDataCellOffset(moderateValue)
	assert.NilError(t, err)
	assert.Equal(t, length, int64(37))
	assert.Equal(t, offset, int64(8712))
}

type value struct {
	id   int64
	name string
}

type node struct {
	name   string
	id     int64
	values []value
}

func TestSpecialCharacters(t *testing.T) {
	hive := openHive("special", t)
	root, err := hive.Root()
	assert.NilError(t, err)
	children, err := hive.NodeChildren(root)
	assert.NilError(t, err)

	names := map[string]node{}
	for _, child := range children {
		childName, err := hive.NodeName(child)
		assert.NilError(t, err)
		_, ok := names[childName]
		assert.Assert(t, !ok)

		values, err := hive.NodeValues(child)
		assert.NilError(t, err)

		vals := make([]value, len(values))
		for idx := range vals {
			data, err := hive.NodeValueKey(values[idx])
			assert.NilError(t, err)
			vals[idx].id = values[idx]
			vals[idx].name = data
		}

		names[childName] = node{
			name:   childName,
			id:     child,
			values: vals,
		}
	}

	val1, ok := names["abcd_äöüß"]
	assert.Assert(t, ok)
	assert.Assert(t, len(val1.values) == 1)
	assert.Assert(t, val1.values[0].name == "abcd_äöüß")

	val2, ok := names["zero\x00key"]
	assert.Assert(t, ok)
	assert.Assert(t, len(val2.values) == 1)
	assert.Assert(t, val2.values[0].name == "zero\x00val")

	val3, ok := names["weird™"]
	assert.Assert(t, ok)
	assert.Assert(t, len(val3.values) == 1)
	assert.Assert(t, val3.values[0].name == "symbols $£₤₧€")
}

func setValues(t *testing.T, hive *Hivex) int64 {
	root, err := hive.Root()
	assert.NilError(t, err)
	assert.Assert(t, root != 0)

	addedA, err := hive.NodeAddChild(root, "A")
	assert.NilError(t, err)
	assert.Assert(t, addedA != 0)

	addedB, err := hive.NodeAddChild(root, "B")
	assert.NilError(t, err)
	assert.Assert(t, addedB != 0)

	getB, err := hive.NodeGetChild(root, "B")
	assert.NilError(t, err)
	assert.Assert(t, getB != 0)

	vals := []HiveValue{
		HiveValue{
			Type:  RegBinary,
			Key:   "Key1",
			Value: []byte("ABC"),
		},
		HiveValue{
			Type:  RegBinary,
			Key:   "Key2",
			Value: []byte("DEF"),
		},
	}
	ret, err := hive.NodeSetValues(getB, vals)
	assert.NilError(t, err)
	assert.Assert(t, ret == 0)

	getKey1, err := hive.NodeGetValue(getB, "Key1")
	assert.NilError(t, err)
	assert.Assert(t, getKey1 != 0)

	valType, length, err := hive.NodeValueType(getKey1)
	assert.NilError(t, err)
	assert.Assert(t, valType == 3)
	assert.Assert(t, length == 3)

	valType, valVal, err := hive.ValueValue(getKey1)
	assert.NilError(t, err)
	assert.Assert(t, valType == 3)
	assert.Assert(t, string(valVal) == "ABC")

	getKey2, err := hive.NodeGetValue(getB, "Key2")
	assert.NilError(t, err)
	assert.Assert(t, getKey2 != 0)

	valType2, length, err := hive.NodeValueType(getKey2)
	assert.NilError(t, err)
	assert.Assert(t, valType2 == 3)
	assert.Assert(t, length == 3)

	valType2, valVal, err = hive.ValueValue(getKey2)
	assert.NilError(t, err)
	assert.Assert(t, valType2 == 3)
	assert.Assert(t, string(valVal) == "DEF")

	return getB

}

func TestWrite(t *testing.T) {
	hive := openHive("minimal", t)
	setValues(t, hive)
}

func TestSetValue(t *testing.T) {
	hive := openHive("minimal", t)
	getB := setValues(t, hive)

	value1 := HiveValue{
		Type:  RegBinary,
		Key:   "Key3",
		Value: []byte("GHI"),
	}

	val1Ret, err := hive.NodeSetValue(getB, value1)
	assert.NilError(t, err)
	assert.Assert(t, val1Ret == 0)

	val1Get, err := hive.NodeGetValue(getB, "Key3")
	assert.NilError(t, err)
	assert.Assert(t, val1Get != 0)

	valType1, valVal1, err := hive.ValueValue(val1Get)
	assert.NilError(t, err)
	assert.Assert(t, valType1 == 3)
	assert.Assert(t, string(valVal1) == "GHI")

	value2 := HiveValue{
		Type:  RegBinary,
		Key:   "Key1",
		Value: []byte("JKL"),
	}

	val2Ret, err := hive.NodeSetValue(getB, value2)
	assert.NilError(t, err)
	assert.Assert(t, val2Ret == 0)

	val2Get, err := hive.NodeGetValue(getB, "Key1")
	assert.NilError(t, err)
	assert.Assert(t, val2Get != 0)

	valType2, valVal, err := hive.ValueValue(val2Get)
	assert.NilError(t, err)
	assert.Assert(t, valType2 == 3)
	assert.Assert(t, string(valVal) == "JKL")
}
