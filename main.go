package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"mygitlab/mysshgui/gui"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	_ "github.com/mattn/go-sqlite3"
)

var gItems []string

func main() {
	addMode := flag.Bool("add", false,
		"addition mode")
	rmMode := flag.Bool("rm", false,
		"remove mode")
	shortVal := flag.String("s", "", "azione")
	actionVal := flag.String("a", "", "comando")
	guiGo := flag.Bool("gui", false, "attiva gui")

	flag.Parse()

	if len(flag.Args()) > 0 {
		arg0 := flag.Args()[0]
		if *shortVal == "" && flag.Args()[0] != "" {
			*shortVal = arg0
		}
	}
	//fmt.Fprintf(os.Stderr,
	//	"Hai inserito: %s %s %b\n", *shortVal, *actionVal, *addMode)

	//fmt.Println("Short: ", *shortVal)

	db, err :=
		sql.Open("sqlite3", dbPath())
	panicOnErr(err)
	defer db.Close()

	_, err = os.Stat(dbPath())
	if os.IsNotExist(err) {
		create(db)
	}

	//dir, err := os.Getwd()
	//panicOnErr(err)

	if *guiGo {
		gui.InitGui()
	} else if *addMode {
		dirInsert(db, *shortVal, *actionVal)
	} else if *shortVal != "" {
		azione := trovaAzione(db, *shortVal)
		if azione != "" {
			esegui(azione)
		} else {
			fmt.Println(*shortVal + " NON TROVATO")
		}
	} else {

		items, azioni := dirList(db)
		var n int
		if len(items) > 20 {
			n = 20
		} else {
			n = len(items)
		}
		gItems = items

		prompt := promptui.Select{
			Label:    "Scegli una azione",
			Items:    items,
			Searcher: cerca,
			Size:     n,
		}

		ind, result, err := prompt.Run()
		panicOnErr(err)
		if ind == 0 {
			fmt.Println("ESCO")
			return
		}
		if *rmMode {
			prompt := promptui.Select{
				Label: "Confermi eliminazione di " + result + "?",
				Items: []string{"no", "si"},
			}

			_, rs, err := prompt.Run()
			panicOnErr(err)

			if rs == "si" {
				ok := rmVal(db, azioni[ind]["short"])
				if ok {
					println("Fatto")
				}
			} else {
				println("Operazione annullata")
			}
		} else {
			fmt.Fprintf(os.Stderr,
				"a s d %s =%d== %s\n", result, ind, azioni[ind]["action"])
			esegui(azioni[ind]["action"])
			/* if err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Print(out)
			} */
			//fmt.Printf("Command finished with error: %v", err)
		}
	}
}

func esegui(c string) {

	var args = strings.Split(c, " ")

	binary, lookErr := exec.LookPath(args[0])
	if lookErr != nil {
		panic(lookErr)
	}

	//args := []string{"telnet", "towel.blinkenlights.nl"}
	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func cerca(k string, i int) bool {
	log.Println(i)
	if strings.Contains(gItems[i], k) {
		return true
	}
	return false
}
func trovaAzione(db *sql.DB, short string) string {
	rows, err := db.Query(`SELECT action FROM
     comandi WHERE short=?`, short)
	panicOnErr(err)
	var s string
	rows.Next()
	rows.Scan(&s)
	return s
}

func dirList(db *sql.DB) ([]string, []map[string]string) {
	items := []string{}
	var azioni []map[string]string
	//var azione map[string]string

	rows, err := db.Query(`SELECT short, action, date FROM
     comandi ORDER BY short ASC`)
	panicOnErr(err)

	//usr, err := user.Current()
	//panicOnErr(err)
	items = append(items, "--ESCI--")
	azioni = append(azioni, make(map[string]string))

	for rows.Next() {
		//fmt.Println(rows)
		var s string
		var t string
		var d string
		err = rows.Scan(&s, &t, &d)
		panicOnErr(err)
		items = append(items, s+" => "+t)
		azione := make(map[string]string)
		azione["short"] = s
		azione["action"] = t
		azioni = append(azioni, azione)
	}

	if len(items) == 0 {
		//items = append(items, usr.HomeDir)
	} else if len(items) > 1 {
		//items = items[1:] // skip first
	}

	return items, azioni
}

func create(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE comandi
     (short text, action text, date text)`)
	panicOnErr(err)

	_, err = db.Exec(`CREATE UNIQUE INDEX
     idx ON comandi (short)`)
	panicOnErr(err)
}

func rmVal(db *sql.DB, short string) bool {
	stmt, err := db.Prepare(`DELETE FROM comandi WHERE short=?`)
	panicOnErr(err)

	_, err = stmt.Exec(short)
	panicOnErr(err)
	return true
}

func dirInsert(db *sql.DB, s string, t string) {
	stmt, err := db.Prepare(`REPLACE INTO
     comandi(short, action, date)
     VALUES(?, ?, datetime('now'))`)
	panicOnErr(err)

	_, err = stmt.Exec(s, t)
	panicOnErr(err)
}

func dbPath() string {
	var dbFile = ".myssh.db"

	usr, err := user.Current()
	panicOnErr(err)
	return path.Join(usr.HomeDir, dbFile)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
