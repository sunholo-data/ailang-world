package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"modernc.org/sqlite"
)

var errCommitAfterDurability = errors.New("injected error after real COMMIT")

const durableFaultDriverName = "world_iter197_durable_fault"

var (
	durableFault             = &durableFaultDriver{}
	registerDurableFaultOnce sync.Once
)

type durableFaultDriver struct {
	armed   atomic.Bool
	reached atomic.Bool
}

func (d *durableFaultDriver) Open(name string) (driver.Conn, error) {
	c, err := (&sqlite.Driver{}).Open(name)
	if err != nil {
		return nil, err
	}
	return &durableFaultConn{Conn: c, fault: d}, nil
}

type durableFaultConn struct {
	driver.Conn
	fault *durableFaultDriver
}

func (c *durableFaultConn) Begin() (driver.Tx, error) {
	tx, err := c.Conn.Begin()
	if err != nil {
		return nil, err
	}
	return &durableFaultTx{Tx: tx, fault: c.fault}, nil
}

func (c *durableFaultConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &durableFaultTx{Tx: tx, fault: c.fault}, nil
}

type durableFaultTx struct {
	driver.Tx
	fault *durableFaultDriver
}

func (tx *durableFaultTx) Commit() error {
	err := tx.Tx.Commit()
	if err != nil {
		return err
	}
	if tx.fault.armed.Swap(false) {
		tx.fault.reached.Store(true)
		return errCommitAfterDurability
	}
	return nil
}

func TestDriverCommitErrorAfterRealCommitIsUncertainAndReconciles(t *testing.T) {
	s := openFileStore(t)
	c := journalCommitFixture(t, s, "driver-fault")
	if _, _, err := s.AppendIntent(context.Background(), "driver-fault", testCommitIntent("driver-fault", c)); err != nil {
		t.Fatal(err)
	}
	path := s.lock.dbPath
	// sql.Register panics on a second registration of one name, so the driver
	// is registered once per process and its flags are reset per run
	// (go test -count=N re-runs this test in the same process).
	registerDurableFaultOnce.Do(func() { sql.Register(durableFaultDriverName, durableFault) })
	fault := durableFault
	fault.armed.Store(false)
	fault.reached.Store(false)
	replacement, err := sql.Open(durableFaultDriverName, path)
	if err != nil {
		t.Fatal(err)
	}
	replacement.SetMaxOpenConns(1)
	if _, err := replacement.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = replacement.Close()
		t.Fatal(err)
	}
	if err := s.db.Close(); err != nil {
		_ = replacement.Close()
		t.Fatal(err)
	}
	s.db = replacement // openFileStore's cleanup closes this DB and releases the original lock.
	fault.armed.Store(true)
	err = s.Commit(context.Background(), c)
	var uncertain *UncertainError
	if !fault.reached.Load() {
		t.Fatal("fault did not reach the real COMMIT")
	}
	if !errors.As(err, &uncertain) || !errors.Is(uncertain.Cause, errCommitAfterDurability) {
		t.Fatalf("err = %v; want uncertain with post-COMMIT sentinel", err)
	}
	rc, ok, err := s.GetReceipt(context.Background(), "driver-fault")
	if err != nil || !ok || rc.State != ReceiptResolved {
		t.Fatalf("receipt = %v, found = %v, err = %v; want resolved", rc.State, ok, err)
	}
}

func TestCommitBodyUsesCallerContext(t *testing.T) {
	path := filepath.Join(repoRootFromCaller(t), "host/store/store.go")
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, source, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Commit" {
			continue
		}
		body := string(source[fset.Position(fn.Body.Pos()).Offset:fset.Position(fn.Body.End()).Offset])
		if !strings.Contains(body, "selectedHeadTx(ctx, tx)") {
			t.Error("Commit head read must use caller ctx")
		}
		if !strings.Contains(body, "tx.ExecContext(ctx,\n\t\t`INSERT OR IGNORE INTO worlds") {
			t.Error("Commit world insert must use ExecContext with caller ctx")
		}
		return
	}
	t.Fatal("Commit not found")
}
