package validate

import "testing"

func TestCard(t *testing.T) {
	tests := []struct {
		name    string
		brand   string
		bank    string
		number  string
		holder  string
		expiry  string
		cvv     string
		wantErr bool
	}{
		{
			name:    "empty brand",
			brand:   "",
			bank:    "sberbank",
			number:  "4111111111111111",
			holder:  "IVAN IVANOV",
			expiry:  "12/25",
			cvv:     "123",
			wantErr: true,
		},
		{
			name:    "empty bank",
			brand:   "visa",
			bank:    "",
			number:  "4111111111111111",
			holder:  "IVAN IVANOV",
			expiry:  "12/25",
			cvv:     "123",
			wantErr: true,
		},
		{
			name:    "empty holder",
			brand:   "visa",
			bank:    "sberbank",
			number:  "4111111111111111",
			holder:  "",
			expiry:  "12/25",
			cvv:     "123",
			wantErr: true,
		},
		{
			name:    "invalid expiry format",
			brand:   "visa",
			bank:    "sberbank",
			number:  "4111111111111111",
			holder:  "IVAN IVANOV",
			expiry:  "1225",
			cvv:     "123",
			wantErr: true,
		},
		{
			name:    "invalid expiry month",
			brand:   "visa",
			bank:    "sberbank",
			number:  "4111111111111111",
			holder:  "IVAN IVANOV",
			expiry:  "13/25",
			cvv:     "123",
			wantErr: true,
		},
		{
			name:    "short cvv",
			brand:   "visa",
			bank:    "sberbank",
			number:  "4111111111111111",
			holder:  "IVAN IVANOV",
			expiry:  "12/25",
			cvv:     "12",
			wantErr: true,
		},
		{
			name:    "long cvv",
			brand:   "visa",
			bank:    "sberbank",
			number:  "4111111111111111",
			holder:  "IVAN IVANOV",
			expiry:  "12/25",
			cvv:     "12345",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Card(tt.brand, tt.bank, tt.number, tt.holder, tt.expiry, tt.cvv)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
