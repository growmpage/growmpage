package growhelper

import (
	"net/http/httptest"
	"testing"
)

// func TestGitPrivateVersion(t *testing.T) {
// 	type args struct {
// 		remote bool
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		wantLength int
// 	}{
// 		{
// 			args:       args{remote: true},
// 			wantLength: 7,
// 		},
// 		{
// 			args:       args{remote: false},
// 			wantLength: 7,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GitPrivateVersion(tt.args.remote); len(got) != tt.wantLength {
// 				t.Errorf("GitPrivateVersion() = %v, want %v, err: %v", len(got), tt.wantLength, got)
// 			}
// 		})
// 	}
// }

//TODO: enable if public
// func TestGitPublicVersion(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want string
// 	}{
// 		{
// 			want: "v1.0.0",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GitPublicVersion(); got != tt.want {
// 				t.Errorf("GitPublicVersion() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestGitPublicUpgrade(t *testing.T) {
// 	type args struct {
// 		w       *httptest.ResponseRecorder
// 		version string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			args: args{w: httptest.NewRecorder(), version: "v0.5.1"},
// 			want: "upgraded from v0.5.1 to v1.0.0, reloading growmpage...",
// 		},
// 		{
// 			args: args{w: httptest.NewRecorder(), version: "v0.5.0"},
// 			want: "upgraded from v0.5.0 to v1.0.0, reloading growmpage...", //would normaly not happen
// 		},
// 		{
// 			args: args{w: httptest.NewRecorder(), version: "v1.0.0"},
// 			want: "binary allready up to date: v1.0.0",
// 		},
// 		{
// 			args: args{w: httptest.NewRecorder(), version: "v0.5.3"},
// 			want: "upgraded from v0.5.3 to v1.0.0, reloading growmpage...",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			GitPublicUpgrade(tt.args.w, tt.args.version)
// 			if tt.args.w.Body.String() != tt.want {
// 				t.Errorf("got %v want %v",
// 					tt.args.w.Body.String(), tt.want)
// 			}
// 		})
// 	}
// }

func TestGitPublicInstall(t *testing.T) {
	type args struct {
		w       *httptest.ResponseRecorder
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{w: httptest.NewRecorder(), version: "v0.5.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GithubPublicInstall(tt.args.w, tt.args.version)
		})
	}
}

// Could destroy git history
// func TestGitPrivateReset(t *testing.T) {
// 	type args struct {
// 		w *httptest.ResponseRecorder
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{args: args{w: httptest.NewRecorder()}, want: "allready up to date: "},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			GitPrivateReset(tt.args.w)
// 			if !strings.Contains(tt.args.w.Body.String(), tt.want) {
// 				t.Errorf("handler returned unexpected body: got %v want %v",
// 					tt.args.w.Body.String(), tt.want)
// 			}
// 		})
// 	}
// }
