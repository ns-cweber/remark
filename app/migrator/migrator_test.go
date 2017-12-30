package migrator

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMigrator_RemoveOldBackupFiles(t *testing.T) {
	loc := "/tmp/remark-backups.test"
	defer os.RemoveAll(loc)

	os.MkdirAll(loc, 0700)
	for i := 1; i <= 10; i++ {
		fname := fmt.Sprintf("%s/backup-site1-201712%02d.gz", loc, i)
		err := ioutil.WriteFile(fname, []byte("blah"), 0600)
		assert.Nil(t, err)
	}

	removeOldBackupFiles(loc, "site1", 3)
	ff, err := ioutil.ReadDir(loc)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(ff), "should keep 3 files only")
	assert.Equal(t, "backup-site1-20171208.gz", ff[0].Name())
	assert.Equal(t, "backup-site1-20171209.gz", ff[1].Name())
	assert.Equal(t, "backup-site1-20171210.gz", ff[2].Name())
}

func TestMigrator_MakeBackup(t *testing.T) {
	loc := "/tmp/remark-backups.test"
	defer os.RemoveAll(loc)
	os.MkdirAll(loc, 0700)

	fname, err := makeBackup(&mockExporter{}, loc, "site1")
	assert.NoError(t, err)
	expFile := fmt.Sprintf("/tmp/remark-backups.test/backup-site1-%s.gz", time.Now().Format("20060102"))
	assert.Equal(t, expFile, fname)

	fi, err := os.Lstat(expFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(52), fi.Size())
}

type mockExporter struct{}

func (mock *mockExporter) Export(w io.Writer, siteID string) error {
	w.Write([]byte("some export blah blah 1234567890"))
	return nil
}
