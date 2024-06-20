package data

import (
	"context"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func Test_getShorturlID(t *testing.T) {
	db, cleanup, err := NewMongoDB("mongodb://localhost:27017", "shorturl")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	type args struct {
		ctx context.Context
		db  *mongo.Database
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				db:  db,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getShorturlID(tt.args.ctx, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("getShorturlID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getShorturlID() = %v, want %v", got, tt.want)
			}
		})
	}
}
