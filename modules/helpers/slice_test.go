package helpers

import (
	"reflect"
	"testing"
)

func TestFlatMap(t *testing.T) {
	type args struct {
		i      interface{}
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Plain map return same map",
			args: args{
				i: map[string]string{
					"test": "1",
					"algo": "2",
				},
			},
			want: map[string]string{
				"test": "1",
				"algo": "2",
			},
			wantErr: false,
		},
		{
			name: "Flatten just maps",
			args: args{
				prefix: "user-meta",
				i: map[string]interface{}{
					"pvr": map[string]string{
						"sdk":     "value",
						"another": "value2",
					},
					"algo": "2",
				},
			},
			want: map[string]string{
				"user-meta/pvr.sdk":     "value",
				"user-meta/pvr.another": "value2",
				"user-meta/algo":        "2",
			},
			wantErr: false,
		},
		{
			name: "Flatten arrays",
			args: args{
				prefix: "user-meta",
				i: map[string]interface{}{
					"pvr": map[string]string{
						"sdk":     "value",
						"another": "value2",
					},
					"ips": []interface{}{
						map[string]string{
							"value": "192.168.1.1",
							"name":  "eth0",
						},
						map[string]string{
							"value": "127.0.0.1",
							"name":  "lo",
						},
					},
					"algo": "2",
				},
			},
			want: map[string]string{
				"user-meta/pvr.sdk":     "value",
				"user-meta/pvr.another": "value2",
				"user-meta/ips.len":     "2",
				"user-meta/ips.0.value": "192.168.1.1",
				"user-meta/ips.0.name":  "eth0",
				"user-meta/ips.1.value": "127.0.0.1",
				"user-meta/ips.1.name":  "lo",
				"user-meta/algo":        "2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlatMap(tt.args.i, tt.args.prefix)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlatMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAreIn(t *testing.T) {
	type args struct {
		in  []string
		all []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Are all the values inside an array",
			args: args{
				[]string{
					"1",
					"2",
					"3",
					"4",
				},
				[]string{
					"3",
					"4",
				},
			},
			want: true,
		},
		{
			name: "Are not all the values inside an array",
			args: args{
				[]string{
					"1",
					"2",
					"3",
					"4",
				},
				[]string{
					"3",
					"4",
					"5",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AllIn(tt.args.in, tt.args.all); got != tt.want {
				t.Errorf("AreIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
