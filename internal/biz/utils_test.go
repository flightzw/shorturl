package biz

import "testing"

func Test_decimalToBase62(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "OK",
			args:    args{num: 0},
			want:    "A",
			wantErr: false,
		},
		{
			name:    "OK",
			args:    args{num: 61},
			want:    "9",
			wantErr: false,
		},
		{
			name:    "OK",
			args:    args{num: 62},
			want:    "BA",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decimalToBase62(tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("decimalToBase62() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decimalToBase62() = %v, want %v", got, tt.want)
			}
		})
	}
}
