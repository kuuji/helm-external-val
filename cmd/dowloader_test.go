package cmd

import (
	"errors"
	"testing"
)

func TestParseUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name              string
		args              args
		wantProtocol      string
		wantNamespace     string
		wantConfigMapName string
		wantKey           string
		wantErr           error
	}{
		{
			name:              "Should return default namespace",
			args:              args{"cm://helm-values"},
			wantProtocol:      "cm",
			wantNamespace:     "default",
			wantConfigMapName: "helm-values",
			wantKey:           "values.yaml",
			wantErr:           nil,
		},
		{
			name:              "Should return namespace and name",
			args:              args{"cm://kuuji/helm-values"},
			wantProtocol:      "cm",
			wantNamespace:     "kuuji",
			wantConfigMapName: "helm-values",
			wantKey:           "values.yaml",
			wantErr:           nil,
		},
		{
			name:              "Should return namespace and name",
			args:              args{"cm://kuuji/helm-values/values-key"},
			wantProtocol:      "cm",
			wantNamespace:     "kuuji",
			wantConfigMapName: "helm-values",
			wantKey:           "values-key",
			wantErr:           nil,
		},
		{
			name:              "Missing config should fail",
			args:              args{"cm://"},
			wantProtocol:      "cm",
			wantNamespace:     "",
			wantConfigMapName: "",
			wantKey:           "",
			wantErr:           errors.New("no config provided after protocol"),
		},
		{
			name:              "Bad url should fail",
			args:              args{"weird"},
			wantProtocol:      "weird",
			wantNamespace:     "",
			wantConfigMapName: "",
			wantKey:           "",
			wantErr:           errors.New(":// missing after protocol"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProtocol, gotNamespace, gotConfigMapName, gotKey, err := ParseUrl(tt.args.url)
			if err != nil && errors.Is(err, tt.wantErr) {
				t.Errorf("ParseUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotProtocol != tt.wantProtocol {
				t.Errorf("ParseUrl() gotProtocol = %v, want %v", gotProtocol, tt.wantProtocol)
			}
			if gotNamespace != tt.wantNamespace {
				t.Errorf("ParseUrl() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotConfigMapName != tt.wantConfigMapName {
				t.Errorf("ParseUrl() gotConfigMapName = %v, want %v", gotConfigMapName, tt.wantConfigMapName)
			}
			if gotKey != tt.wantKey {
				t.Errorf("ParseUrl() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
