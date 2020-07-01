package runner

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func TestOutputManager_ValidAndPrint(t *testing.T) {

	tmpDir := os.TempDir() + "/Test_printAndValidOutputDir"
	_ = fileutil.CreateDirIfNotExists(tmpDir, 0755)
	defer func() { _ = fileutil.RemoveDir(tmpDir) }()

	type args struct {
		setup formula.Setup
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Return empty string when dir is empty",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{}},
					TmpOutputDir: func() string {
						basePath := "/t-rit-return-empty"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						return path
					}(),
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Return only the outputs with printValue",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{
						{
							Name:  "X",
							Print: true,
						},
						{
							Name:  "Y",
							Print: false,
						},
						{
							Name:  "Z",
							Print: true,
						},
					}},
					TmpOutputDir: func() string {
						basePath := "/t-rit-printed"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						_ = ioutil.WriteFile(path+"/x", []byte("1"), 0755)
						_ = ioutil.WriteFile(path+"/y", []byte("2"), 0755)
						_ = ioutil.WriteFile(path+"/z", []byte("3"), 0755)
						return path
					}(),
				},
			},
			want:    "X=1\nZ=3\n",
			wantErr: false,
		},
		{
			name: "Return Red when output dir not have all files",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{
						{
							Name:  "X",
							Print: true,
						},
						{
							Name:  "Y",
							Print: false,
						},
						{
							Name:  "Z",
							Print: true,
						},
					}},
					TmpOutputDir: func() string {
						basePath := "/t-rit-err-all-files"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						_ = ioutil.WriteFile(path+"/x", []byte("1"), 0755)
						_ = ioutil.WriteFile(path+"/z", []byte("3"), 0755)
						return path
					}(),
				},
			},
			want:    prompt.Red("Output dir size is different of outputs array in config.json"),
			wantErr: false,
		},
		{
			name: "Return Red when some output file is missing",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{
						{
							Name:  "X",
							Print: true,
						},
						{
							Name:  "Y",
							Print: false,
						},
						{
							Name:  "Z",
							Print: true,
						},
					}},
					TmpOutputDir: func() string {
						basePath := "/t-rit-err-missing-files"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						_ = ioutil.WriteFile(path+"/x", []byte("1"), 0755)
						_ = ioutil.WriteFile(path+"/z", []byte("3"), 0755)
						_ = ioutil.WriteFile(path+"/w", []byte("3"), 0755)
						return path
					}(),
				},
			},
			want:    prompt.Red("file:Y not found in output dir"),
			wantErr: false,
		},
		{
			name: "Return Err when fail to read dir",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{}},
					TmpOutputDir: func() string {
						basePath := "/not-created-dir"
						return basePath
					}(),
				},
			},
			want:    prompt.Red("Fail to read output dir"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			o := OutputManager{
				writer: &buffer,
			}
			if err := o.ValidAndPrint(tt.args.setup); (err != nil) != tt.wantErr {
				t.Errorf("ValidAndPrint() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := buffer.String(); got != tt.want {
				t.Errorf("printAndValidOutputDir() = %v, want %v", got, tt.want)
			}
		})
	}
}