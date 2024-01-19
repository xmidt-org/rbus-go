// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	errUnknown := errors.New("unknown error")
	tests := []struct {
		description string
		url         string
		want        Client
		expectedErr error
	}{
		// Success case
		{
			description: "Valid args",
			url:         "unix://file",
			want: Client{
				network: "unix",
				address: "/file",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			got, err := New(Config{
				URL: tc.url,
			})

			if tc.expectedErr != nil {
				assert.Nil(got)
				if errors.Is(tc.expectedErr, errUnknown) {
					assert.Error(err)
				} else {
					assert.ErrorIs(err, tc.expectedErr)
				}
				return
			}
			require.NotNil(got)
			assert.Equal(&tc.want, got)
		})
	}
}
