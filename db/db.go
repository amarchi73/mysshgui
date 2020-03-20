package db

import (
	"database/sql"
	"fmt"
	"os/user"
	"path"

	"github.com/gotk3/gotk3/gtk"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err error
var dbFile = ".myssh.db"

type Azioni struct {
	Short  string
	Action string
	Date   string
}

// Ciao stampa ciao
func Ciao(s string) string {
	return s + "aaaa"
}

// OpenDB apre la connessione al DB
func OpenDB() {
	usr, err := user.Current()
	panicOnErr(err)
	dbpath := path.Join(usr.HomeDir, dbFile)

	fmt.Println(dbpath)

	dbb, errr := sql.Open("sqlite3", dbpath)
	panicOnErr(errr)

	db = dbb
}

// ElencoAzioni estrae tutte le righe
func ElencoAzioni() map[int]Azioni {
	rows, err := db.Query(`SELECT action, short, date FROM comandi`)
	panicOnErr(err)
	items := make(map[int]Azioni)
	var item Azioni
	var i int = 0

	for rows.Next() {
		//fmt.Println(rows)
		var action string
		var short string
		var date string
		rows.Scan(&action, &short, &date)
		item.Action = action
		item.Short = short
		item.Date = date
		items[i] = item

		fmt.Println("====" + items[i].Short)

		i++
	}

	return items
}

// InserisciAzione aggiorna o inserisce un'azione
func InserisciAzione(a Azioni, oldval string, btn *gtk.Button, azione string) {
	stmt, err := db.Prepare(`DELETE FROM comandi WHERE short=?`)
	panicOnErr(err)

	_, err = stmt.Exec(oldval)
	panicOnErr(err)

	if azione == "elimina" {
		btn.Destroy()
		return
	}
	stmt, err1 := db.Prepare(`REPLACE INTO
     comandi(short, action, date)
     VALUES(?, ?, datetime('now'))`)
	panicOnErr(err1)

	fmt.Println(a)

	_, err = stmt.Exec(a.Short, a.Action)
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
