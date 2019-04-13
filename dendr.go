package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type fileStats struct {
	size  int64
	mtime time.Time
}

func openFileDatabase(fileName string) (*fileDatabase, error) {
	fdb := &fileDatabase{}

	db, err := sql.Open("sqlite3", "file:"+fileName)
	if err != nil {
		panic(err)
	}
	fdb.db = db

	for _, fn := range []func() error{
		fdb.prepareTable,
		fdb.prepareStatements,
	} {
		if err := fn(); err != nil {
			panic(err)
		}
	}

	return fdb, nil
}

type fileDatabase struct {
	db   *sql.DB
	qget *sql.Stmt
	qset *sql.Stmt
}

func (fdb *fileDatabase) Close() error {
	for _, c := range []io.Closer{
		fdb.qset,
		fdb.qget,
		fdb.db,
	} {
		if c != nil {
			if err := c.Close(); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

const (
	file_            = "file"
	path_            = "path"
	size_            = "size"
	mtime_           = "mtime"
	createtablefile_ = "create table if not exists " + file_ + " (" + path_ + " primary key, " + size_ + " integer, " + mtime_ + " timestamp)"
	get_             = "select " + size_ + ", " + mtime_ + " from " + file_ + " where " + path_ + " = :" + path_
	set_             = "insert into " + file_ + " (" + path_ + ", " + size_ + ", " + mtime_ + ") values (:" + path_ + ", :" + size_ + ", :" + mtime_ + ") on conflict(" + path_ + ") do update set " + size_ + "=excluded." + size_ + ", " + mtime_ + "=excluded." + mtime_
)

func (fdb *fileDatabase) prepareTable() error {
	_, err := fdb.db.Exec(createtablefile_)
	if err != nil {
		panic(err)
	}

	return nil
}

func (fdb *fileDatabase) prepareOneStatement(s **sql.Stmt, q string) error {
	if *s != nil {
		if err := (*s).Close(); err != nil {
			panic(err)
		}
	}
	if stmt, err := fdb.db.Prepare(q); err != nil {
		panic(err)
	} else {
		*s = stmt
	}
	return nil
}

func (fdb *fileDatabase) prepareStatements() error {
	for _, v := range []struct {
		h **sql.Stmt
		q string
	}{
		{&fdb.qget, get_},
		{&fdb.qset, set_},
	} {
		if err := fdb.prepareOneStatement(v.h, v.q); err != nil {
			panic(err)
		}
	}
	return nil
}

func (fdb *fileDatabase) get(path string) (fileStats, bool, error) {
	var stats fileStats
	npath := sql.Named(path_, path)
	row := fdb.qget.QueryRow(npath)
	err := row.Scan(&stats.size, &stats.mtime)
	if err == sql.ErrNoRows {
		return fileStats{}, false, nil
	}
	if err != nil {
		panic(err)
	}
	return stats, true, nil
}

func (fdb *fileDatabase) set(path string, stats fileStats) error {
	npath := sql.Named(path_, path)
	nsize := sql.Named(size_, stats.size)
	nmtime := sql.Named(mtime_, stats.mtime)
	result, err := fdb.qset.Exec(npath, nsize, nmtime)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}
	return nil
}

func main() {
	fdb, err := openFileDatabase("test.db")
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getcwd: %v\n", err)
		return
	}

	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return nil
		}
		newStats := fileStats{info.Size(), info.ModTime()}
		oldStats, ok, err := fdb.get(path)
		if err != nil {
			return err
		}
		if ok {
			sameSize := oldStats.size == newStats.size
			sameMtime := oldStats.mtime.Equal(newStats.mtime)
			if sameSize && sameMtime {
				// do nothing?
			} else {
				fmt.Print("=")
				if sameSize {
					fmt.Print(".")
				} else {
					fmt.Print("s")
				}
				if sameMtime {
					fmt.Print(".")
				} else {
					fmt.Print("m")
				}
				fmt.Println("  ", path)
			}
		} else {
			fmt.Println("+sm  ", path)
		}
		err = fdb.set(path, newStats)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking: %v\n", err)
		return
	}
}
