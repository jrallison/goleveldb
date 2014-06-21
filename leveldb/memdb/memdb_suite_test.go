package memdb

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jrallison/goleveldb/leveldb/testutil"
)

func TestMemdb(t *testing.T) {
	testutil.RunDefer()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Memdb Suite")
}
