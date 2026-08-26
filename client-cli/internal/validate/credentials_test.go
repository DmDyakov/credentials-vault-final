// client-cli/internal/validate/credentials_test.go
package validate

import "testing"

func TestCredentials(t *testing.T) {
	tests := []struct {
		name     string
		site     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid credentials",
			site:     "example.com",
			username: "bob",
			password: "secret123",
			wantErr:  false,
		},
		{
			name:     "valid site with subdomain",
			site:     "mail.example.com",
			username: "bob",
			password: "secret123",
			wantErr:  false,
		},
		{
			name:     "valid site with path",
			site:     "example.com/login",
			username: "bob",
			password: "secret123",
			wantErr:  false,
		},
		{
			name:     "empty site",
			site:     "",
			username: "bob",
			password: "secret123",
			wantErr:  true,
		},
		{
			name:     "site without dot",
			site:     "example",
			username: "bob",
			password: "secret123",
			wantErr:  true,
		},
		{
			name:     "site without dot - localhost",
			site:     "localhost",
			username: "bob",
			password: "secret123",
			wantErr:  true,
		},
		{
			name:     "empty username",
			site:     "example.com",
			username: "",
			password: "secret123",
			wantErr:  true,
		},
		{
			name:     "short password",
			site:     "example.com",
			username: "bob",
			password: "12345",
			wantErr:  true,
		},
		{
			name:     "password exactly 6 chars",
			site:     "example.com",
			username: "bob",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "empty password",
			site:     "example.com",
			username: "bob",
			password: "",
			wantErr:  true,
		},
		{
			name:     "all empty",
			site:     "",
			username: "",
			password: "",
			wantErr:  true,
		},
		{
			name:     "site with spaces",
			site:     " example.com ",
			username: "bob",
			password: "secret123",
			wantErr:  false,
		},
		{
			name:     "username with spaces",
			site:     "example.com",
			username: " bob ",
			password: "secret123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Credentials(tt.site, tt.username, tt.password)
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
