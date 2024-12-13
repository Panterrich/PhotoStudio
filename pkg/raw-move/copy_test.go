package rawmove //nolint:testpackage

import (
	"testing"
)

func TestModifyFileName(t *testing.T) {
	type args struct {
		fileName string
		camera   string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1DX2",
			args: args{
				fileName: "AF0I3540.JPG",
				camera:   "1DX2",
			},
			want: "1DX2_3540.JPG",
		},
		{
			name: "5DSR",
			args: args{
				fileName: "279A0236.JPG",
				camera:   "5DSR",
			},
			want: "5DSR_0236.JPG",
		},
		{
			name: "5DM3",
			args: args{
				fileName: "1L5A2878.JPG",
				camera:   "5DM3",
			},
			want: "5DM3_2878.JPG",
		},
		{
			name: "D4",
			args: args{
				fileName: "DCS_0010.JPG",
				camera:   "D4",
			},
			want: "D4_0010.JPG",
		},
		{
			name: "Renama 1",
			args: args{
				fileName: "DCS_0010 (1).JPG",
				camera:   "D4",
			},
			want: "D4_0010_1.JPG",
		},
		{
			name: "Rename 2",
			args: args{
				fileName: "DCS_0010(1).JPG",
				camera:   "D4",
			},
			want: "D4_0010_1.JPG",
		},
		{
			name: "5DSR",
			args: args{
				fileName: "279A0236(1).JPG",
				camera:   "5DSR",
			},
			want: "5DSR_0236_1.JPG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := modifyFileName(tt.args.fileName, tt.args.camera); got != tt.want {
				t.Errorf("modifyFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
