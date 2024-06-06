package autoGen

import (
	"testing"
)

func TestCreateTemp(t *testing.T) {
	type args struct {
		autoCode AutoCodeStruct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "order_log",
			args: args{AutoCodeStruct{
				StructName:       "Order",
				TableName:        "order",
				HumpPackageName:  "order",
				Abbreviation:     "repo",
				AutoMoveFile:     true,
				AutoMoveFilePath: "../../app/service/order/internal/data/repo/",
				Fields:           nil,
				Package:          "data",
				Module:           "autoGenExamples",
				ServerName:       "order",
				InterfacePath:    "../../app/service/order/internal/biz/interface.go",
			}},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveAutoCode(tt.args.autoCode); (err != nil) != tt.wantErr {
				t.Errorf("createTemp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
