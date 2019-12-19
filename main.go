package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"

	_ "github.com/go-sql-driver/mysql"
	stomp "github.com/go-stomp/stomp"
)

type Log struct {
	TimeStamp string `json:"timestamp"`
	Command   string `json:"command"`
	Number    string `json:"number"`
}

func newSQL() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:sql101@tcp(localhost:3306)/demodb")

	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newREDIS() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "",
		DB:          0,
		PoolSize:    80,
		PoolTimeout: time.Minute,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Redis Fail Ping")
		return nil, err
	}
	result, err := client.Exists("LOG").Result()
	if err != nil {
		return nil, err
	} else if result == 0 {
		err = client.Set("LOG", "", 0).Err()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func newAMQ() (*stomp.Conn, error) {
	conn, err := stomp.Dial("tcp", "localhost:61613", stomp.ConnOpt.HeartBeat(0, 0))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

var bads = [...]string{"02", "20", "04", "40", "30", "03", "006", "70", "60", "06", "90", "08", "000", "64", "46", "22", "44", "66", "414", "616", "646", "545", "54", "45", "41", "14", "940"}

func main() {
	var wg sync.WaitGroup
	fmt.Println("Demo redis sql activemq")
	db, err := newSQL()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	client, err := newREDIS()
	if err != nil {
		panic(err.Error())
	}
	defer client.Close()
	conn, err := newAMQ()
	if err != nil {
		panic(err.Error())
	}
	defer conn.Disconnect()
	wg.Add(1)
	wg.Add(1)
	go receiveRequest(conn, db, client, wg)
	go receiveLogRequest(conn, client, wg)
	wg.Wait()
	log.Println("FINISH")
}

func receiveLogRequest(conn *stomp.Conn, client *redis.Client, wg sync.WaitGroup) {
	defer wg.Done()
	sub, err := conn.Subscribe("/topic/logreq", stomp.AckAuto)
	if err != nil {
		panic(err.Error())
	}
	for {
		msg := <-sub.C
		var smsg string = string(msg.Body)
		if smsg == "" {
			continue
		}
		log.Println(smsg)
		result, err := client.Get("LOG").Result()
		if err != nil {
			log.Println(err)
		}
		result = "[" + result[1:] + "]"
		err = conn.Send(
			"/topic/logres/", // destination
			"text/plain",     // content-type
			[]byte(result))   // body
		if err != nil {
			log.Println(err)
		}

	}

	err = sub.Unsubscribe()
	if err != nil {
		panic(err.Error())
	}
}

func receiveRequest(conn *stomp.Conn, db *sql.DB, client *redis.Client, wg sync.WaitGroup) {
	defer wg.Done()
	sub, err := conn.Subscribe("/topic/request", stomp.AckAuto)
	if err != nil {
		panic(err.Error())
	}

	for {
		msg := <-sub.C
		var smsg string = string(msg.Body)
		if smsg == "" && len(smsg) < 13 {
			continue
		}
		log.Println(smsg)
		logToRedis(client, smsg)

		val, err := client.Get(smsg).Result()
		if err != nil {
			switch smsg[0:3] {
			case "csr":
				responsevalue := checkSimple(smsg[3:], db)
				err = conn.Send(
					"/topic/response/"+smsg, // destination
					"text/plain",            // content-type
					[]byte(responsevalue))   // body
				if err != nil {
					log.Println(err)
				}
				err = client.Set(smsg, responsevalue, 0).Err()
				if err != nil {
					fmt.Println(err)
				}
			default:
				err = conn.Send(
					"/topic/response/"+smsg,   // destination
					"text/plain",              // content-type
					[]byte("Invalid command")) // body
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			err = conn.Send(
				"/topic/response/"+smsg, // destination
				"text/plain",            // content-type
				[]byte(val))             // body
			if err != nil {
				log.Println(err)
			}
		}

	}

	err = sub.Unsubscribe()
	if err != nil {
		panic(err.Error())
	}

}

func logToRedis(client *redis.Client, msg string) {

	objectjson, err := json.Marshal(Log{TimeStamp: time.Now().UTC().String(), Command: msg[:3], Number: msg[3:]})
	if err != nil {
		panic(err)
	}
	err = client.Append("LOG", ","+string(objectjson)).Err()
	if err != nil {
		panic(err)
	}
}

func checkSimple(number string, db *sql.DB) string {
	result, err := db.Query(fmt.Sprintf("SELECT simple FROM fortunetable WHERE number='%s'", number))

	if err != nil {
		panic(err.Error())
	}

	data := make(map[string]string)
	for result.Next() {
		var simple_r string
		err = result.Scan(&simple_r)
		if err != nil {
			panic(err.Error())
		}
		data["result"] = simple_r
		jsons, err := json.Marshal(data)
		if err != nil {
			panic(err.Error())
		}
		return string(jsons)
	}

	//row not exist add to db
	res := calSimple(number)
	_, err = db.Exec(fmt.Sprintf("INSERT INTO `fortunetable`(`number`, `simple`) VALUES ('%s','%s')", number, res))

	if err != nil {
		panic(err.Error())
	}

	data["result"] = res
	jsons, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}
	return string(jsons)
}

func calSimple(number string) string {
	sub := number[3:]
	for _, s := range bads {
		if strings.Contains(sub, s) {
			return "Bad"
		}
	}
	return "Good"
}
