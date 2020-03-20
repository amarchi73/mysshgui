package gui

import (
	"fmt"
	"log"
	"mygitlab/mysshgui/db"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gtk"
)

// looks like handlers can literally be any function or method
func b1Clicked() {
	fmt.Println("b1 clicked")
}

func b2Clicked() {
	fmt.Println("b2 clicked")
}

func b3Clicked() {
	fmt.Println("b3 clicked")
}

// you just place them in a map that names the signals, then feed the map to the builder
var signals = map[string]interface{}{
	"popupSave":     popupSave,
	"popupDel":      popupDel,
	"modalOK":       modalOK,
	"modalNo":       modalNo,
	"B3":            b3Clicked,
	"nuovoElemento": nuovoElemento,
}

func modalOK(b *gtk.Button) {
	d, _ := builder.GetObject("dialog1")

	var newAction db.Azioni
	newAction.Action, _ = curAction.action.GetText()
	newAction.Short, _ = curAction.short.GetText()
	db.InserisciAzione(newAction, curAction.shortOld, curAction.bottone, curAction.azione)
	d.(*gtk.Dialog).Hide()
}
func modalNo(b *gtk.Button) {
	d, _ := builder.GetObject("dialog1")
	d.(*gtk.Dialog).Hide()
}

func click(btt *gtk.Button) {
	//s, _ := btt.GetLabel()
	//fmt.Println("click %s aaaa", s)
	//btt.Destroy()

	nome, _ := btt.GetName()
	fmt.Println(fmt.Sprintf("Bottone %s", nome))
	a, _ := builder.GetObject("actionEdit")
	actn := azioni[nome].Action
	a.(*gtk.Entry).SetText(actn)
	shr := azioni[nome].Short
	b, _ := builder.GetObject("labelEdit")
	b.(*gtk.Entry).SetText(shr)
	//close(sigDel)

	curAction.short = b.(*gtk.Entry)
	curAction.action = a.(*gtk.Entry)
	curAction.shortOld = shr
	curAction.bottone = btt
}

func popupDel(b *gtk.Button) {
	nome, _ := b.GetLabel()
	lID, _ := b.GetName()
	lbl, _ := builder.GetObject("popUpLabel")
	lbl.(*gtk.Label).SetText("Eliminare?" + nome + "--" + lID)
	fmt.Println("click %s", "dialog")
	d, _ := builder.GetObject("dialog1")
	d.(*gtk.Dialog).Show()

	curAction.azione = "elimina"
}

func popupSave(b *gtk.Button) {
	nome, _ := b.GetLabel()
	lID, _ := b.GetName()
	lbl, _ := builder.GetObject("popUpLabel")
	lbl.(*gtk.Label).SetText("Salvare?" + nome + "--" + lID)
	fmt.Println("click %s", "dialog")
	d, _ := builder.GetObject("dialog1")
	d.(*gtk.Dialog).Show()

	curAction.azione = "salva"
}

func nuovoElemento() {
	fmt.Println("nuovoElemento")

	obj, _ := builder.GetObject("window")
	lb1, _ := builder.GetObject("listbox1")

	btt, _ := gtk.ButtonNewWithMnemonic("myButton")
	btt.Connect("clicked", click, btt)
	var a db.Azioni
	a.Short = "xxx"
	a.Action = "---"
	fmt.Println(a)
	nomeButt := fmt.Sprintf("Bottone%s", a.Short)
	azioni[nomeButt] = a
	btt.SetName(nomeButt)
	btt.SetLabel(fmt.Sprintf("xxx%s", a.Short))
	lb1.(*gtk.ListBox).Add(btt)
	fmt.Println(nomeButt)
	btt.SetLabel("xxx")

	wnd := obj.(*gtk.Window)
	wnd.ShowAll()
}

/* ------- TreeView --------- */
// Appends single value to the TreeView's model
func AppendToList(value string, c int) {
	i := ListStore.Append()
	ListStore.SetValue(i, 0, value)
	ListStore.SetValue(i, 1, value+"xx")
}

// Appends several values to the TreeView's model
func AppendMultipleToList(values ...string) {
	for _, v := range values {
		AppendToList(v, 0)
	}
}

func listaClick(selection *gtk.TreeSelection) {
	var iter *gtk.TreeIter
	var model gtk.ITreeModel
	var ok bool
	model, iter, ok = selection.GetSelected()
	if ok {
		tpath, err := model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			log.Printf("treeSelectionChangedCB: Could not get path from model: %s\n", err)
			return
		}
		log.Printf("treeSelectionChangedCB: selected path: %s\n", tpath)
		ListStore.SetValue(iter, 0, "xxx")
	}
}

/* ------- FINE TV ---------- */

var builder *gtk.Builder
var app *gtk.Application
var err error

var sigDel = make(chan string, 1)
var azioni map[string]db.Azioni
var curAction struct {
	short    *gtk.Entry
	action   *gtk.Entry
	shortOld string
	azione   string
	bottone  *gtk.Button
}
var ListStore *gtk.ListStore

func InitGui() {
	const appID = "com.retc3.mytest"
	app, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatalln("Couldn't create app:", err)
	}

	db.OpenDB()
	act := make(map[int]db.Azioni)
	act = db.ElencoAzioni()
	azioni = make(map[string]db.Azioni)

	// It looks like all builder code must execute in the context of `app`.
	// If you try creating the builder inside the main function instead of
	// the `app` "activate" callback, then you will get a segfault
	app.Connect("activate", func() {
		// Use this instead if you have your glade XML in a separate file
		// builder, err := gtk.BuilderNewFromFile("mytest.glade")

		if 1 == 0 {
			builder, err = gtk.BuilderNew()
			if err != nil {
				log.Fatalln("Couldn't make builder:", err)
			}
			err = builder.AddFromString(gladeTemplate)
		} else {
			builder, err = gtk.BuilderNewFromFile("finestra.glade")
		}
		if err != nil {
			log.Fatalln("Couldn't make builder:", err)
		}
		builder.ConnectSignals(signals)

		obj, _ := builder.GetObject("window")
		lb1, _ := builder.GetObject("listbox1")

		for i, _ := range act {
			btt, _ := gtk.ButtonNewWithMnemonic("myButton")
			btt.Connect("clicked", click, btt)
			nomeButt := fmt.Sprintf("BottoneNome %s", act[i].Short)
			azioni[nomeButt] = act[i]
			btt.SetName(nomeButt)
			btt.SetLabel(fmt.Sprintf("%s", act[i].Short))

			lb1.(*gtk.ListBox).Add(btt)
		}

		/*btn, _ := gtk.LabelNew("ciao")*/
		/*tv1, _ := builder.GetObject("treeView1")
		for i := 1; i < 4; i++ {
			col, _ := gtk.TreeViewColumnNew()
			nomecol := fmt.Sprintf("Titolo %d", i)
			col.SetTitle(nomecol)
			tv1.(*gtk.TreeView).AppendColumn(col)
		}
		col := tv1.(*gtk.TreeView).GetColumn(1)
		store, _ := gtk.store*/

		//TreeView, _ := gtk.TreeViewNew() //builder.GetObject("treeView1")
		TreeView, _ := builder.GetObject("treeView1")
		//Entry, _ := gtk.EntryNew()
		ListStore, _ = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
		// TreeView properties
		{
			renderer, _ := gtk.CellRendererTextNew()
			column, _ := gtk.TreeViewColumnNewWithAttribute("Value", renderer, "text", 0)
			TreeView.(*gtk.TreeView).AppendColumn(column)
			column1, _ := gtk.TreeViewColumnNewWithAttribute("Value1", renderer, "text", 1)
			TreeView.(*gtk.TreeView).AppendColumn(column1)
		}
		TreeView.(*gtk.TreeView).SetModel(ListStore)
		AppendMultipleToList("Go", "Docker", "CockroachDB")

		sel, _ := TreeView.(*gtk.TreeView).GetSelection()
		sel.Connect("changed", listaClick)

		wnd := obj.(*gtk.Window)
		wnd.ShowAll()
		app.AddWindow(wnd)
	})

	//app.Run(os.Args)
	var noArgs []string
	app.Run(noArgs)
}
