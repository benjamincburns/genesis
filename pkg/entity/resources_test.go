/*
	Copyright 2019 whiteblock Inc.
	This file is a part of the genesis.

	Genesis is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Genesis is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package entity

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResources_GetMemory_Successful(t *testing.T) {
	var tests = []struct {
		res      Resources
		expected int64
	}{
		{res: Resources{
			Cpus:   "",
			Memory: "45",
		}, expected: int64(45)},
		{res: Resources{
			Cpus:   "",
			Memory: "1",
		}, expected: int64(1)},
		{res: Resources{
			Cpus:   "",
			Memory: "92233720368547",
		}, expected: int64(92233720368547)},
		{res: Resources{
			Cpus:   "",
			Memory: "3gb",
		}, expected: int64(3000000000)},
		{res: Resources{
			Cpus:   "",
			Memory: "6KB",
		}, expected: int64(6000)},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			num, err := tt.res.GetMemory()
			assert.NoError(t, err)

			assert.Equal(t, num, tt.expected)
		})
	}
}

func TestResources_GetMemory_Unsuccessful(t *testing.T) {
	var tests = []struct {
		res Resources
	}{
		{res: Resources{
			Cpus:   "",
			Memory: "45.46",
		}},
		{res: Resources{
			Cpus:   "",
			Memory: "35273409857203948572039458720349857",
		}},
		{res: Resources{
			Cpus:   "",
			Memory: "s",
		}},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := tt.res.GetMemory()

			assert.NotEqual(t, err, nil)
		})
	}
}

func TestResources_NoLimits(t *testing.T) {
	var tests = []struct {
		res Resources
		expected bool
	}{
		{
			res: Resources{
				Memory: "",
				Cpus: "",
			},
			expected: true,
		},
		{
			res: Resources{
				Memory: "5gb",
				Cpus: "",
			},
			expected: false,
		},
		{
			res: Resources{
				Memory: "",
				Cpus: "5",
			},
			expected: false,
		},
		{
			res: Resources{
				Memory: "4gb",
				Cpus: "6",
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.res.NoLimits(), tt.expected)
		})
	}
}

func TestResources_NoCPULimits(t *testing.T) {
	var tests = []struct {
		res Resources
		expected bool
	}{
		{
			res: Resources{
				Cpus: "",
			},
			expected: true,
		},
		{
			res: Resources{
				Cpus: "5",
			},
			expected: false,
		},
		{
			res: Resources{
				Cpus: " ",
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.res.NoCPULimits(), tt.expected)
		})
	}
}

func TestResources_NoMemoryLimits(t *testing.T) {
	var tests = []struct {
		res Resources
		expected bool
	}{
		{
			res: Resources{
				Memory: "",
			},
			expected: true,
		},
		{
			res: Resources{
				Memory: " ",
			},
			expected: false,
		},
		{
			res: Resources{
				Memory: "5GB",
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.res.NoMemoryLimits(), tt.expected)
		})
	}
}
