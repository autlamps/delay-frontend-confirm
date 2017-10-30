package confirm

import (
	"database/sql"
	"fmt"

	"net/http"

	"log"

	"github.com/autlamps/delay-frontend-confirm/data"
	"github.com/gorilla/mux"
)

type Conf struct {
	DBURL string
}

type Env struct {
	Users data.UserStore
}

func Create(c Conf) (*mux.Router, error) {
	db, err := sql.Open("postgres", c.DBURL)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to Database. Err: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Cannot ping to Database. Err: %v", err)
	}

	env := Env{
		Users: data.InitUserService(db),
	}

	r := mux.NewRouter()
	r.HandleFunc("/confirm/{user_id}", env.confirm)

	return r, nil
}

func (e *Env) confirm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id := vars["user_id"]

	user, err := e.Users.GetUser(user_id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed"))
		return
	}

	if user.EmailConfirmed {
		w.WriteHeader(http.StatusNotFound)
	} else {
		err := e.Users.ConfEmail(user_id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}
}
