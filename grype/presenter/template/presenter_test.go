package template

import (
	"bytes"
	"flag"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anchore/go-testutils"
	"github.com/anchore/grype/grype/presenter/internal"
)

var update = flag.Bool("update", false, "update the *.golden files for template presenters")

func TestPresenter_Present(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	templateFilePath := path.Join(workingDirectory, "./test-fixtures/test.template")

	pb := internal.GeneratePresenterConfig(t, internal.ImageSource)

	templatePresenter := NewPresenter(pb, templateFilePath)

	var buffer bytes.Buffer
	if err := templatePresenter.Present(&buffer); err != nil {
		t.Fatal(err)
	}

	actual := buffer.Bytes()

	if *update {
		testutils.UpdateGoldenFileContents(t, actual)
	}
	expected := testutils.GetGoldenFileContents(t)

	assert.Equal(t, string(expected), string(actual))
}

func TestPresenter_SprigDate_Fails(t *testing.T) {
	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	// this template has the generic sprig date function, which is intentionally not supported for security reasons
	templateFilePath := path.Join(workingDirectory, "./test-fixtures/test.template.sprig.date")

	pb := internal.GeneratePresenterConfig(t, internal.ImageSource)

	templatePresenter := NewPresenter(pb, templateFilePath)

	var buffer bytes.Buffer
	err = templatePresenter.Present(&buffer)
	require.ErrorContains(t, err, `function "now" not defined`)
}
