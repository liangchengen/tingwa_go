package main

import (
	"database/sql"
	"fmt"
	"github.com/antchfx/htmlquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)


func main()  {
	var(
		db *sql.DB
		err error
		id_sql int
		id int
		name string
		play int
		musiclike int
		url string

	)

	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/web?charset=utf8")

	rows, err := db.Query("SELECT number FROM number")
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&id_sql)
		checkErr(err)
	}

	for id=id_sql ; id < 120000; id++{
		name,play,musiclike,url = httpGet(id)

		if (play != 0 || name != "" || url != "" || musiclike != 0) {
			stmt,err := db.Prepare("INSERT INTO music SET id=?,name=?,play=?,url=?,musiclike=?")
			checkErr(err)
			_,err =stmt.Exec(id,name,play,url,musiclike)
			checkErr(err)
		}

		stmt, err := db.Prepare("update number set number=?")
		checkErr(err)
		_, err = stmt.Exec(id)
		checkErr(err)
	}
	db.Close()

}

func httpGet(id int)(name string,play int,like int,url string) {
	json_like,err :=http.Get("https://www.itingwa.com/?c=event&m=get_like&id=" + strconv.Itoa(id))
	checkErr(err)
	defer json_like.Body.Close()

	body_byte,err := ioutil.ReadAll(json_like.Body)
	checkErr(err)
	str := string(body_byte)
	list_bool := strings.Contains(str, "null")

	if list_bool {
		fmt.Printf("\033[1;31;40m%s\033[0m\n",strconv.Itoa(id) + "  null")
		return
	}

	htm, err := http.Get("https://www.itingwa.com/listen/"+ strconv.Itoa(id))
	checkErr(err)
	defer htm.Body.Close()
	htm_byte,err := ioutil.ReadAll(htm.Body)

	doc, err := htmlquery.Parse(strings.NewReader(string(htm_byte)))
	checkErr(err)

	name_xpath := htmlquery.FindOne(doc, "//h1/text()")
	play_xpath := htmlquery.FindOne(doc, "//ul/li[3]/font/text()")
	url_xpath := htmlquery.FindOne(doc, "//div[@id='tw_player']/@init-data")
	if (name_xpath == nil || play_xpath == nil || url_xpath == nil) {
		return
	}

	like,_ = strconv.Atoi(gjson.Get(str, "total").String())
	play,_ = strconv.Atoi(htmlquery.InnerText(play_xpath))
	name = strings.TrimSpace(htmlquery.InnerText(name_xpath))
	url = htmlquery.InnerText(url_xpath)
	println(strconv.Itoa(id)+ "  " + name)
	return name,play,like,url

}

func checkErr(err error)  {
	if err != nil {
		println(err)
	}
}
