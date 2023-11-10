/*
go env -w GOOS=linux
go build
go env -w GOOS=windows

acceptHeader := r.Header.Get("Accept")
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	bolt "go.etcd.io/bbolt"
)

func main() {
	aid := Aid()
	fmt.Println(aid)
	fmt.Println(Aid10(aid))
	fmt.Println(Aid16(aid))
	fmt.Println(Aid36(aid))
	return

	// 数据库初始化
	db, _ := bolt.Open("a.db", 0600, nil)
	db.Batch(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("^"))
		return nil
	})
	defer db.Close()

	// 模板加载
	tmpl := template.Must(template.ParseFiles("a.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			db.Batch(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("^"))
				key := []byte(strconv.FormatInt(123456789, 12))
				fmt.Printf("key: %s\n", key)
				b.Put(key, []byte("abc"))
				//res := b.Get(key)
				c := b.Cursor()
				sk, _ := c.Seek([]byte(strconv.FormatInt(1, 12)))
				fmt.Printf("%s\n", sk)
				skk, _ := c.Next()
				fmt.Printf("%s\n", skk)
				return nil
			})

			fmt.Fprintf(w, "fuck")
			return
		} else if r.URL.Path == "/" {
			fmt.Println(Aid36(aid))
			// 模板渲染
			err := tmpl.Execute(w, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			//fmt.Print(r.URL.Query())
			switch r.URL.Path[1] {
			case '@':
				fmt.Fprintf(w, "User")
			case '&':
				fmt.Fprintf(w, "Tag")
			default:
				fmt.Fprintf(w, "Page")
			}
		}
	})

	http.ListenAndServe(":3001", nil)

	/*
		db.Batch(func(tx *bolt.Tx) error {
			tx.CreateBucketIfNotExists([]byte("test"))
			b := tx.Bucket([]byte("test"))
			key := []byte(strconv.FormatInt(123456789, 12))
			fmt.Printf("key: %s\n", key)
			b.Put(key, []byte("abc"))
			//res := b.Get(key)
			c := b.Cursor()
			sk, _ := c.Seek([]byte(strconv.FormatInt(1, 12)))
			fmt.Printf("%s\n", sk)
			skk, _ := c.Next()
			fmt.Printf("%s\n", skk)
			return nil
		})
	*/

}
