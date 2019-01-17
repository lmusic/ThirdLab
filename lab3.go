package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var _hash string
var c chan string

func main() {
	_type := os.Args[1]
	n,_ := strconv.Atoi(os.Args[2])
	if _type == "serv" {
		Serv()
	}
	if _type == "client" {
		for i:=0; i < n; i++  {
			go Client()
		}
	}
	stop:="stop"
	fmt.Fscan(os.Stdin, &stop)
}
func Serv(){
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		log.Printf("accepted connection from %v", conn.RemoteAddr())
		go handle(conn)
	}
}
func handle(conn net.Conn){
	message, _ := bufio.NewReader(conn).ReadString('\n')
	arr := strings.Split(message,"&")
	if len(arr)>0 {
		_hash = arr[1]
	}
	key := arr[0]
	returnedKey := next_session_key(key)
	conn.Write([]byte(returnedKey + "\n"))
	for i:=0; i<10; i++ {
		returnedKey = next_session_key(returnedKey)
		message, _ = bufio.NewReader(conn).ReadString('\n')
		arr := strings.Split(message,"&")
		if len(arr)>1 {
			_hash = arr[1]
		}
		key = arr[0][0:len(arr[0])-1]
		if key == returnedKey {
			returnedKey = next_session_key(returnedKey)
			conn.Write([]byte(returnedKey + "\n"))
		}
	}
	conn.Close()
	println("finished")
}
func Client (){
	rand.Seed(time.Now().UnixNano())
	set_hash_str()
	key := get_session_key()

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")

	fmt.Fprintf(conn, key +"&"+_hash+"\n")

	for i:=0;i<10 ;i++  {
		returnedStr, _ := bufio.NewReader(conn).ReadString('\n')
		returnedKey := returnedStr[0:len(returnedStr)-1]
		key = next_session_key(key)
		if key == returnedKey {
			println("finished")
			key = next_session_key(key)
			fmt.Fprintf(conn, key +"\n")
		} else{
			println("error")
		}
	}
}
func get_session_key() string {
	rand.Seed(time.Now().UnixNano())
	var result string
	for i:=0;i<10 ;i++  {
		var elem int
		elem = 9*rand.Intn(10000)+1
		result += strconv.Itoa(elem)[0:1]
	}
	return result
}

func set_hash_str() {
	var result string
	for i:=0;i<5 ;i++  {
		var elem int
		elem = 6*rand.Intn(10000)+1
		result += strconv.Itoa(elem)[0:1]
	}
	_hash = result
}

func next_session_key(session_key string) string{
	result := 0
	for i:= 0;i<5;i++{
		id := _hash[i:i+1]
		t_result,_ := strconv.Atoi(_calc_hash(session_key, id))
		result += t_result
	}
	return ("0000000000"+strconv.Itoa(result))[len("0000000000"+strconv.Itoa(result))-10:]
}
func _calc_hash(session_key string, val string) string {
	intVal,_ := strconv.ParseInt(val,10,0)
	result := ""
	if intVal == 1{
		param,_ := strconv.Atoi(session_key[0:5])
		expr := "00" + strconv.Itoa(param%97)
		return expr[len(expr)-2:]
	}
	if intVal == 2{
		for i:= len(session_key)-1; i>=0; i-- {
			result += session_key[i:i+1]
		}
		return result
	}
	if intVal == 3{
		return session_key[len(session_key)-5:]+session_key[0:5]
	}
	if intVal == 4{
		num := 0
		for i:=1; i<9;i++  {
			val,_ := strconv.Atoi(session_key[i:i+1])
			num += val +41
		}
		return strconv.Itoa(num)
	}
	if intVal == 5{
		var num int
		for i:=0; i<len(session_key); i++  {
			ord,_:= utf8.DecodeRuneInString(session_key[i:i+1])
			resInt := ord^43
			if resInt>=48 && resInt<=56{
				symbol := string(resInt)
				val,_:= strconv.ParseInt(symbol,10,0)
				num += int(val)
			} else{
				num += int(resInt)
			}
		}
		return strconv.Itoa(num)
	}
	p1,_:= strconv.Atoi(session_key)
	p2,_:= strconv.Atoi(val)
	result = strconv.Itoa(p1+p2)
	return result
}
