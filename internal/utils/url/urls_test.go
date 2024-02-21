package url

import "testing"

func TestIncludeBasicAuth(t *testing.T) {
	tests := []struct {
		name string
		url  string
		user string
		pass string
		want string
	}{
		{
			name: "valid url",
			url:  "https://example.com",
			user: "user",
			pass: "pass",
			want: "https://user:pass@example.com",
		},
		{
			name: "invalid url",
			url:  "example.com",
			user: "user",
			pass: "pass",
			want: "https://user:pass@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IncludeBasicAuth(tt.url, tt.user, tt.pass)
			if got != tt.want {
				t.Errorf("IncludeBasicAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
