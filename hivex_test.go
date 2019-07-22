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
	"testing"
)

func TestOpenClose(t *testing.T) {
	hive, err := NewHivex("testdata/minimal", WRITE)
	if err != nil {
		t.Errorf("Error: %q\n", err)
	}
	err = hive.Close()
	if err != nil {
		t.Errorf("Error: %q\n", err)
	}
}

// TODO(gsamfira): add more tests
