package main

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"wQueue/db"
	"wQueue/ws"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var (
	c        *gorm.DB
	cons     []ws.WConn
	store    = sessions.NewCookieStore([]byte("SESSION_KEY"))
	upgrader = websocket.Upgrader{}
)

func main() {
	gob.Register(db.User{})

	var err error
	c, err = db.SetupConn()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/about", about)
	r.HandleFunc("/help", help)
	r.HandleFunc("/", home)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/queues/{title}", queue)
	r.HandleFunc("/ws", handleWs)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := &http.Server{
		Addr:              ":1234",
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func home(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["title"] = "Home Page"
	m["info"] = "This is new home info!"
	addLoginToMap(r, m)

	q, err := db.GetAllQueues(c)
	if err != nil {
		log.Print(err)
		return
	}

	sort.Slice(q, func(i, j int) bool {
		return q[i].Open && !q[j].Open
	})

	for i := range q {
		count, err := db.CountQueueitemsInQueue(c, &q[i])
		if err != nil {
			log.Print(err)
			return
		}
		q[i].InQueue = count
	}

	m["queues"] = q

	tpl := template.Must(template.ParseFiles("html/base.html", "html/home.html"))
	err = tpl.ExecuteTemplate(w, "base.html", m)
	if err != nil {
		log.Print(err)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["title"] = "About us"
	addLoginToMap(r, m)

	tpl := template.Must(template.ParseFiles("html/base.html", "html/about.html"))
	err := tpl.ExecuteTemplate(w, "base.html", m)
	if err != nil {
		log.Print(err)
	}
}

func help(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["title"] = "Help/FAQ"
	addLoginToMap(r, m)

	tpl := template.Must(template.ParseFiles("html/base.html", "html/help.html"))
	err := tpl.ExecuteTemplate(w, "base.html", m)
	if err != nil {
		log.Print(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		return
	}

	user := r.Form.Get("username")
	pass := r.Form.Get("password")

	if user != "" && pass != "" {
		correct, u, err := db.Login(c, user, pass)
		if err != nil {
			log.Print(err)
			return
		}

		if correct {
			session, err := store.Get(r, "s")
			if err != nil {
				log.Print(err)
				return
			}
			session.Values["user"] = u
			err = session.Save(r, w)
			if err != nil {
				log.Print(err)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	m := make(map[string]interface{})
	m["title"] = "Login"
	addLoginToMap(r, m)

	tpl := template.Must(template.ParseFiles("html/base.html", "html/login.html"))
	err = tpl.ExecuteTemplate(w, "base.html", m)
	if err != nil {
		log.Print(err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "s")
	if err != nil {
		log.Print(err)
		return
	}
	session.Values["user"] = nil
	err = session.Save(r, w)
	if err != nil {
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Print(err)
		return
	}

	q, err := db.GetQueueByTitle(c, string(p))
	if err != nil {
		log.Print(err)
		return
	}
	wc := ws.WConn{
		Qid:  q.Id,
		Conn: conn,
	}

	cons = append(cons, wc)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			cons = ws.RemoveConnection(cons, conn)
			return
		}
	}
}

func queue(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		return
	}
	location := r.Form.Get("location")
	comment := r.Form.Get("comment")
	leave := r.Form.Get("leave")
	action := r.Form.Get("action")

	vars := mux.Vars(r)
	title := vars["title"]

	q, err := db.GetQueueByTitle(c, title)
	if err != nil {
		log.Print(err)
		return
	}

	qs, err := db.GetQueueitemsInQueue(c, &q)
	if err != nil {
		log.Print(err)
		return
	}

	session, err := store.Get(r, "s")
	if err != nil {
		log.Print(err)
		return
	}

	if session.Values["user"] == nil {
		m := make(map[string]interface{})
		m["title"] = q.Title
		m["message"] = q.Message
		m["open"] = q.Open
		m["qitems"] = qs
		addLoginToMap(r, m)

		funcs := template.FuncMap{"add": func(a, b int) int { return a + b }}
		tpl := template.Must(template.New("queuePage").Funcs(funcs).ParseFiles("html/base.html", "html/queue.html"))
		tpl = template.Must(tpl.ParseFiles("html/wsScript.html", "html/admin.html", "html/closedForm.html", "html/openForm.html", "html/qTable.html"))
		err := tpl.ExecuteTemplate(w, "base.html", m)
		if err != nil {
			log.Print(err)
		}
		return
	}

	if action != "" {
		switch action {
		case "Close":
			err := db.AdminCloseQueue(c, &q)
			if err != nil {
				log.Print(err)
				return
			}
			go ws.ReloadWs(cons, &q)

		case "Open":
			err := db.AdminOpenQueue(c, &q)
			if err != nil {
				log.Print(err)
				return
			}
			go ws.ReloadWs(cons, &q)

		case "Remove":
			removeID := r.Form.Get("removeId")
			rInt, err := strconv.Atoi(removeID)
			if err != nil {
				log.Print(err)
				return
			}

			userToRemove, err := db.GetUserByID(c, rInt)
			if err != nil {
				log.Print(err)
				return
			}

			removeQi, err := db.GetQueueitem(c, &userToRemove, &q)
			if err != nil {
				log.Print(err)
				return
			}

			err = db.RemoveFromQueue(c, &removeQi)
			if err != nil {
				log.Print(err)
				return
			}

			for i := range qs {
				if qs[i].Id == removeQi.Id {
					qs = append(qs[:i], qs[i+1:]...)
					break
				}
			}
			go ws.ReloadWs(cons, &q)

		case "Text":
			newText := r.Form.Get("text")
			err := db.SetQueueMessage(c, &q, newText)
			if err != nil {
				log.Print(err)
				return
			}
			go ws.ReloadWs(cons, &q)

		case "Message":
			newMessage := r.Form.Get("message")
			go ws.MessageWs(cons, &q, newMessage)
		}
	}

	userNotInQueue := true
	var u db.User
	if session.Values["user"] != nil {
		u = session.Values["user"].(db.User)
	} else {
		return
	}
	for i := range qs {
		if qs[i].Userid == u.Id {
			userNotInQueue = false
		}
	}

	if location != "" && comment != "" {
		if userNotInQueue {
			err := db.AddToQueue(c, &u, &q, location, comment)
			if err != nil {
				log.Print(err)
				return
			}
			qi := db.Queueitem{
				Userid:   u.Id,
				Queueid:  q.Id,
				Location: location,
				Active:   true,
				Comment:  comment,
			}
			qs = append(qs, qi)
			userNotInQueue = false
			ws.ReloadWs(cons, &q)
		}
	}

	if leave == "1" {
		if !userNotInQueue {
			qi, err := db.GetQueueitem(c, &u, &q)
			if err != nil {
				log.Print(err)
				return
			}

			err = db.RemoveFromQueue(c, &qi)
			if err != nil {
				log.Print(err)
				return
			}

			for i := range qs {
				if qs[i].Id == qi.Id {
					qs = append(qs[:i], qs[i+1:]...)
					break
				}
			}
			userNotInQueue = true
			ws.ReloadWs(cons, &q)
		}
	}

	m := make(map[string]interface{})
	err = db.UserIsAdminOfQueue(c, &u, &q)
	if err != nil {
		m["isadmin"] = false
	} else {
		m["isadmin"] = true
	}
	m["inqueue"] = !userNotInQueue
	m["title"] = q.Title
	m["message"] = q.Message
	m["open"] = q.Open
	m["qitems"] = qs
	addLoginToMap(r, m)

	funcs := template.FuncMap{"add": func(a, b int) int { return a + b }}
	tpl := template.Must(template.New("queuePage").Funcs(funcs).ParseFiles("html/base.html", "html/queue.html"))
	tpl = template.Must(tpl.ParseFiles("html/wsScript.html", "html/admin.html", "html/closedForm.html", "html/openForm.html", "html/qTable.html"))
	err = tpl.ExecuteTemplate(w, "base.html", m)
	if err != nil {
		log.Print(err)
	}
}

func addLoginToMap(r *http.Request, m map[string]interface{}) {
	session, err := store.Get(r, "s")
	if err != nil {
		log.Print(err)
		return
	}

	if session.Values["user"] != nil {
		m["log"] = "Logout " + session.Values["user"].(db.User).Username
		m["logPath"] = "logout"
		m["loggedin"] = true
	} else {
		m["log"] = "Login"
		m["logPath"] = "login"
		m["loggedin"] = false
	}
}
