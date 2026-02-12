package yamlconfig_test

import (
	"os"
	"testing"

	"github.com/sculley/yamlconfig"
	"github.com/stretchr/testify/require"
)

type TestConfigStruct struct {
	String string   `yaml:"string"`
	Int    int      `yaml:"int"`
	Bool   bool     `yaml:"bool"`
	Slice  []string `yaml:"slice"`
	Unit   uint     `yaml:"unit"`
	Float  float64  `yaml:"float"`
	Struct struct {
		String string `yaml:"string"`
	} `yaml:"struct"`
}

type TestConfigEmpty struct {
	String string `yaml:"string"`
}

type TestConfigEmptyStruct struct {
	String string `yaml:"string"`
	Struct struct {
		String string `yaml:"string"`
		Int    int    `yaml:"int"`
	} `yaml:"struct"`
}

type TestConfigOmitEmpty struct {
	String string            `yaml:"string"`
	Map    map[string]string `yaml:"map" yamlconfig:"omitempty"`
	Slice  []string          `yaml:"slice" yamlconfig:"omitempty"`
}

type TestConfigOmitEmptyNestedStruct struct {
	String string `yaml:"string"`
	Nested struct {
		Required string `yaml:"string"`
	} `yaml:"nested" yamlconfig:"omitempty"`
}

type TestConfigOmitEmptyNestedStructPointer struct {
	Backend string `yaml:"backend"`
	S3      *struct {
		Bucket string `yaml:"bucket"`
	} `yaml:"s3" yamlconfig:"omitempty"`
	FileSystem *struct {
		Path string `yaml:"path"`
	} `yaml:"file-system" yamlconfig:"omitempty"`
}

func TestConfig(t *testing.T) {
	t.Run("Load Config", func(t *testing.T) {
		cfg := TestConfigStruct{}

		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "config.test.yml")
		if tempConfigFileErr != nil {
			t.Fatal(tempConfigFileErr)
		}
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("string: test\nint: 1\nbool: true\nslice:\n  - foo\n  - bar\nunit: 1\nfloat: 1.0\nstruct:\n  string: test")
		if writeStringErr != nil {
			t.Fatal(writeStringErr)
		}

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		if loadConfigErr != nil {
			t.Fatal(loadConfigErr)
		}

		require.NoError(t, loadConfigErr)

		require.Equal(t, "test", cfg.String)
		require.Equal(t, 1, cfg.Int)
		require.Equal(t, true, cfg.Bool)
		require.Equal(t, []string{"foo", "bar"}, cfg.Slice)
		require.Equal(t, uint(1), cfg.Unit)
		require.Equal(t, 1.0, cfg.Float)
		require.Equal(t, "test", cfg.Struct.String)
	})

	t.Run("Load Config Open File Error", func(t *testing.T) {
		cfg := TestConfigStruct{}
		loadConfigErr := yamlconfig.LoadConfig("nonexistent.yml", &cfg)

		require.Error(t, loadConfigErr)
	})

	t.Run("Load Config Decode Error", func(t *testing.T) {
		cfg := TestConfigStruct{}

		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "invalid_config.test.yml")
		if tempConfigFileErr != nil {
			t.Fatal(tempConfigFileErr)
		}
		defer os.Remove(tempConfigFile.Name())

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)

		require.Error(t, loadConfigErr)
	})

	t.Run("Load Config Must Be Pointer To Struct", func(t *testing.T) {
		cfg := &TestConfigStruct{}

		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "config.test.yml")
		if tempConfigFileErr != nil {
			t.Fatal(tempConfigFileErr)
		}
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("string: test\nint: 1\nbool: true\nslice:\n  - foo\n  - bar\nunit: 1\nfloat: 1.0\nstruct:\n  string: test\n")
		if writeStringErr != nil {
			t.Fatal(writeStringErr)
		}

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)

		require.Error(t, loadConfigErr)
	})

	t.Run("Load Config Missing Config", func(t *testing.T) {
		cfg := TestConfigEmpty{}

		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "config.test.yml")
		if tempConfigFileErr != nil {
			t.Fatal(tempConfigFileErr)
		}
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("int: 1\n")
		if writeStringErr != nil {
			t.Fatal(writeStringErr)
		}

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)

		require.Error(t, loadConfigErr)
	})

	t.Run("Load Config Missing Config Struct", func(t *testing.T) {
		cfg := TestConfigEmptyStruct{}

		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "config.test3.yml")
		if tempConfigFileErr != nil {
			t.Fatal(tempConfigFileErr)
		}
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("string: test\nstruct:\n  string: test\n")
		if writeStringErr != nil {
			t.Fatal(writeStringErr)
		}

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)

		require.Error(t, loadConfigErr)
	})

	t.Run("Load Config With OmitEmpty - All Fields Present", func(t *testing.T) {
		cfg := TestConfigOmitEmpty{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "omit_empty_config_all.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("string: test\nmap:\n  foo: bar\nslice:\n  - item1\n  - item2\n")
		require.NoError(t, writeStringErr)

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.NoError(t, loadConfigErr)

		require.Equal(t, "test", cfg.String)
		require.Equal(t, map[string]string{"foo": "bar"}, cfg.Map)
		require.Equal(t, []string{"item1", "item2"}, cfg.Slice)
	})

	t.Run("Load Config With OmitEmpty - Optional Fields Missing", func(t *testing.T) {
		cfg := TestConfigOmitEmpty{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "omit_empty_config_missing.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("string: test\n")
		require.NoError(t, writeStringErr)

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.NoError(t, loadConfigErr)

		require.Equal(t, "test", cfg.String)
		require.Nil(t, cfg.Map)
		require.Nil(t, cfg.Slice)
	})

	t.Run("Load Config With OmitEmpty - Missing Required Field", func(t *testing.T) {
		cfg := TestConfigOmitEmpty{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "omit_empty_config_missing_required.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("map:\n  foo: bar\nslice:\n  - item1\n")
		require.NoError(t, writeStringErr)

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.Error(t, loadConfigErr)
	})

	t.Run("Load Config With OmitEmpty Nested Struct - Absent Skips Child Validation", func(t *testing.T) {
		cfg := TestConfigOmitEmptyNestedStruct{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "omit_empty_nested_struct.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		// nested absent - child validation skipped (value struct: empty = absent)
		_, writeStringErr := tempConfigFile.WriteString("string: test\n")
		require.NoError(t, writeStringErr)

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.NoError(t, loadConfigErr)

		require.Equal(t, "test", cfg.String)
		require.Equal(t, "", cfg.Nested.Required)
	})

	t.Run("Load Config With OmitEmpty Pointer - Absent Skips Child Validation", func(t *testing.T) {
		cfg := TestConfigOmitEmptyNestedStructPointer{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "omit_empty_pointer_absent.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("backend: s3\n")
		require.NoError(t, writeStringErr)

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.NoError(t, loadConfigErr)

		require.Equal(t, "s3", cfg.Backend)
		require.Nil(t, cfg.S3)
		require.Nil(t, cfg.FileSystem)
	})

	t.Run("Load Config With OmitEmpty Pointer - Present Validates Children", func(t *testing.T) {
		cfg := TestConfigOmitEmptyNestedStructPointer{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "omit_empty_pointer_present_valid.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("backend: s3\ns3:\n  bucket: my-bucket\n")
		require.NoError(t, writeStringErr)

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.NoError(t, loadConfigErr)

		require.Equal(t, "s3", cfg.Backend)
		require.NotNil(t, cfg.S3)
		require.Equal(t, "my-bucket", cfg.S3.Bucket)
	})

	t.Run("Load Config With OmitEmpty Pointer - Present With Missing Required Child Fails", func(t *testing.T) {
		cfg := TestConfigOmitEmptyNestedStructPointer{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "omit_empty_pointer_present_invalid.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		// s3 key present but bucket missing
		_, writeStringErr := tempConfigFile.WriteString("backend: s3\ns3: {}\n")
		require.NoError(t, writeStringErr)

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.Error(t, loadConfigErr)
	})

	t.Run("Load Config with Bool Field - False", func(t *testing.T) {
		cfg := TestConfigStruct{}
		tempConfigFile, tempConfigFileErr := os.CreateTemp("", "config_bool_false.yml")
		require.NoError(t, tempConfigFileErr)
		defer os.Remove(tempConfigFile.Name())

		_, writeStringErr := tempConfigFile.WriteString("string: test\nint: 1\nbool: false\nslice:\n  - foo\n  - bar\nunit: 1\nfloat: 1.0\nstruct:\n  string: test")
		if writeStringErr != nil {
			t.Fatal(writeStringErr)
		}

		loadConfigErr := yamlconfig.LoadConfig(tempConfigFile.Name(), &cfg)
		require.NoError(t, loadConfigErr)

		require.Equal(t, "test", cfg.String)
		require.Equal(t, 1, cfg.Int)
		require.Equal(t, false, cfg.Bool)
		require.Equal(t, []string{"foo", "bar"}, cfg.Slice)
		require.Equal(t, uint(1), cfg.Unit)
		require.Equal(t, 1.0, cfg.Float)
		require.Equal(t, "test", cfg.Struct.String)
	})
}
