package chat

import (
    "html/template"
    "net/http"
	"time"
    "appengine"
    "appengine/channel"
    "appengine/user"
    "fmt"
)

func init() {
    http.HandleFunc("/main", main)
    http.HandleFunc("/talk", talk)
    http.HandleFunc("/opened", opened)
}

type ChatMessage struct {
    UserX			string
    MessageX		string
    MessageDateX	time.Time
}

var allUsers map[string]string = make(map[string]string)

var mainTemplate = template.Must(template.ParseFiles("chat/main.html"))

func opened(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    u := user.Current(c) 
    allUsers[u.ID] = u.ID
    for _, uId := range allUsers {
    	fmt.Println("opened",uId);
    }
    fmt.Println("created channel");
    fmt.Println(r.FormValue("toto"));
}

func main(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    u := user.Current(c) // assumes 'login: required' set in app.yaml
    
    //check user
    if u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Location", url)
        w.WriteHeader(http.StatusFound)
        return
    }
    
    
    //allUsers := append(allUsers, u.String())
    fmt.Println("test");
    //key := r.FormValue("chatkey")
    

    //tok, err := channel.Create(c, u.ID+key)
    tok, err := channel.Create(c, u.ID)
   
    if err != nil {
        http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
        c.Errorf("channel.Create: %v", err)
        return
    }

    err = mainTemplate.Execute(w, map[string]string{
        "token":    tok,
        "me":       u.ID,
    })
        //"chat_key": key,
    if err != nil {
        c.Errorf("mainTemplate: %v", err)
    }
}

func talk(w http.ResponseWriter, r *http.Request) {
 	fmt.Println("talking ?");
    c := appengine.NewContext(r)

    // Get the user and their proposed move.
    u := user.Current(c)
    msg := r.FormValue("msg")
//    if err != nil {
//        http.Error(w, "Invalid msg", http.StatusBadRequest)
//        return
//    }
    //key := r.FormValue("chatkey")

    g := ChatMessage {
    	UserX : u.String(),
    	MessageX: msg,
        MessageDateX:    time.Now(),
    }

    // Send the chat state to both clients.
    fmt.Println("sending talk with channel ?")
    //msgJson, _ := json.Marshal(g)
    //fmt.Println("chatmsg>%s",msgJson)
    for _, uId := range allUsers {
        fmt.Println("sending to ",uId," msg :",g)
    
        err := channel.SendJSON(c, uId, g)
        //err := channel.SendJSON(c, uId+key, g)
        if err != nil {
            c.Errorf("sending chat: %v", err)
        }
    }
}


