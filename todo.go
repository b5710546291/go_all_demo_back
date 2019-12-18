package main

// func getAllTBI(db *sql.DB, client *redis.Client) string {
// 	result, err := db.Query(fmt.Sprintf("SELECT * FROM fortunetable WHERE indepth='tbi'"))

// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	var phones []Phone
// 	for result.Next() {
// 		var phone Phone
// 		err = result.Scan(&phone.Number, &phone.Simple, &phone.Indepth)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		phones = append(phones, phone)
// 	}
// 	jsons, err := json.Marshal(phones)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return string(jsons)
// }

// func requestInspectNumber(number string, db *sql.DB, client *redis.Client) {
// 	result, err := db.Query(fmt.Sprintf("SELECT indepth FROM fortunetable WHERE number='%s'", number))
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	for result.Next() {
// 		var s string
// 		err = result.Scan(&s)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		if s == "" {
// 			_, err = db.Exec(fmt.Sprintf("UPDATE `fortunetable` SET `indepth`='tbi' WHERE number='%s'", number))

// 			if err != nil {
// 				panic(err.Error())
// 			}
// 		}
// 		return
// 	}
// 	res := calSimple(number)
// 	_, err = db.Exec(fmt.Sprintf("INSERT INTO `fortunetable`(`number`, `simple`, `indepth`) VALUES ('%s','%s','tbi')", number, res))

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	return
// }

// func getNumber(number string, db *sql.DB, client *redis.Client) string {
// 	result, err := db.Query(fmt.Sprintf("SELECT * FROM fortunetable WHERE number='%s'", number))
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	for result.Next() {
// 		var phone Phone
// 		err = result.Scan(&phone.Number, &phone.Simple, &phone.Indepth)
// 		if err != nil {
// 			panic(err.Error())
// 		}

// 		jsons, err := json.Marshal(phone)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		return string(jsons)
// 	}
// 	return "NotExist"
// }

// func getAllNumber(db *sql.DB, client *redis.Client) string {
// 	result, err := db.Query(fmt.Sprintf("SELECT * FROM fortunetable"))

// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	var phones []Phone
// 	for result.Next() {
// 		var phone Phone
// 		err = result.Scan(&phone.Number, &phone.Simple, &phone.Indepth)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		phones = append(phones, phone)
// 	}
// 	jsons, err := json.Marshal(phones)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return string(jsons)
// }

// func checkIndepth(number string, db *sql.DB) string {
// 	result, err := db.Query(fmt.Sprintf("SELECT indepth FROM fortunetable WHERE number='%s'", number))
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	data := make(map[string]string)
// 	for result.Next(){
// 		var indepth_r string
// 		err = result.Scan(&indepth_r))
// 		if err != nil{
// 			panic(err.Error())
// 		}
// 		if indepth_r == ""{
// 			_, err = db.Exec(fmt.Sprintf("UPDATE `fortunetable` SET `indepth`='tbi' WHERE number='%s'", number))

// 			if err != nil {
// 				panic(err.Error())
// 			}
// 			data["result"] = "tbi"
// 		} else {
// 			data["result"] = indepth_r
// 		}
// 		jsons, err := json.Marshal(data)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		return string(jsons)

// 	}
// 	//row not exist add to db
// 	res := calSimple(number)
// 	_, err = db.Exec(fmt.Sprintf("INSERT INTO `fortunetable`(`number`, `simple`, `indepth`) VALUES ('%s','%s','tbi')", number, res))

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	data["result"] = "tbi
// 	jsons, err := json.Marshal(data)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return string(jsons)

// }
