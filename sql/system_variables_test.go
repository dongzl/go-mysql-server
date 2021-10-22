// Copyright 2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-errors.v1"
)

var newConn = SystemVariable{
	Name:    "max_connections",
	Scope:   SystemVariableScope_Global,
	Dynamic: true,
	Type:    NewSystemIntType("max_connections", 1, 100000, false),
	Default: int64(1000),
}

var newTimeout = SystemVariable{
	Name:    "net_write_timeout",
	Scope:   SystemVariableScope_Both,
	Dynamic: true,
	Type:    NewSystemIntType("net_write_timeout", 1, 9223372036854775807, false),
	Default: int64(1),
}

var newUnknown = SystemVariable{
	Name:    "net_write_timeout",
	Scope:   SystemVariableScope_Both,
	Dynamic: true,
	Type:    NewSystemIntType("net_write_timeout", 1, 9223372036854775807, false),
	Default: int64(1),
}

func TestInitSystemVariablesWithDefaults(t *testing.T) {

	tests := []struct {
		name             string
		persistedGlobals []SystemVariable
		err              *errors.Kind
		expectedCmp      []SystemVariable
	}{
		{
			name:             "set max_connections",
			persistedGlobals: []SystemVariable{newConn},
			expectedCmp:      []SystemVariable{newConn},
		}, {
			name:             "set two variables",
			persistedGlobals: []SystemVariable{newConn, newTimeout},
			expectedCmp:      []SystemVariable{newConn, newTimeout},
		}, {
			name: "bad type",
			persistedGlobals: []SystemVariable{{
				Name:    "max_connections",
				Scope:   SystemVariableScope_Global,
				Dynamic: true,
				Type:    NewSystemIntType("max_connections", 1, 100000, false),
				Default: "1000",
			}},
			expectedCmp: []SystemVariable{{
				Name:    "max_connections",
				Scope:   SystemVariableScope_Global,
				Dynamic: true,
				Type:    NewSystemIntType("max_connections", 1, 100000, false),
				Default: "1000",
			}},
			err: nil, // TODO: nothing is stopping us from setting incorrect types currently
		}, {
			name:             "unknown system variable",
			persistedGlobals: []SystemVariable{newUnknown},
			expectedCmp:      []SystemVariable{newUnknown},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := InitSystemVariables(test.persistedGlobals)
			if test.err != nil {
				assert.True(t, test.err.Is(err))
				return
			} else {
				assert.NoError(t, err)
			}

			for i, sysVar := range test.persistedGlobals {
				cmp, _, _ := SystemVariables.GetGlobal(sysVar.Name)
				assert.Equal(t, test.expectedCmp[i], cmp)
			}
		})
	}
}
