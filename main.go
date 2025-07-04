package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type BankAccount struct {
	AccHolderName string `json:"acc_holder_name"`
	AccNo string `json:"acc_no"`
	IFSCcode string `json:"ifsc_code"`
	Email string `json:"email"`
	Password string `json:"password"`
	Balance string `json:"balance"`
}

type AddAmount struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Amount string `json:"amount"`
}
type DebitAmount struct{
	Email string `json:"email"`
	Password string `json:"password"`
	Amount string `json:"amount"`
}

func main() {
	db,err:=sql.Open("postgres" , "postgresql://free_d1iq_user:KHOULvBoVAL5elQcbXqE6ZwcRad8nCDr@dpg-d1iecrali9vc73fu5ueg-a.oregon-postgres.render.com/free_d1iq")
	if err !=nil{
		fmt.Println("error")
		return
	}
	fmt.Println("connected")
	defer db.Close()

	http.HandleFunc("/create_account" , func(w http.ResponseWriter , r *http.Request){
		var account BankAccount
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			fmt.Println("json error")
			return
		}
		query := `INSERT INTO bank_account VALUES($1,$2,$3,$4,$5,$6)`
		_,err = db.Exec(query,account.AccHolderName,account.AccNo,account.IFSCcode,account.Email,account.Password,account.Balance)
		if err != nil {
			fmt.Println("query error")
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{"message":"account created successfully"})
		if err != nil {
			fmt.Println("encode error")
			return
		}
	})

	http.HandleFunc("/add_amount",func(w http.ResponseWriter, r*http.Request){
		var userAddAmountDetails AddAmount
		err := json.NewDecoder(r.Body).Decode(&userAddAmountDetails)
		if err != nil {
			fmt.Println("json decode error")
			return
		}

		var databaseAddAmountDetails AddAmount
		query1 := `SELECT password,balance FROM bank_account WHERE email=$1`
		err = db.QueryRow(query1,userAddAmountDetails.Email).Scan(&databaseAddAmountDetails.Password,&databaseAddAmountDetails.Amount)
		if err != nil {
			fmt.Println("query exec error")
			return
		}

		if userAddAmountDetails.Password != databaseAddAmountDetails.Password {
			fmt.Println("incorrect password")
			return
		}

		var balance , amount int
		balance , err = strconv.Atoi(databaseAddAmountDetails.Amount)
		if err != nil {
			fmt.Println("error converting balance")
			return
		}

		amount , err = strconv.Atoi(userAddAmountDetails.Amount)
		if err != nil {
			fmt.Println("error converting amount")
			return
		}

		totalBalance := strconv.Itoa(amount+balance)

		query2 := `UPDATE bank_account SET balance=$1 WHERE email=$2`
		_,err = db.Exec(query2,totalBalance,userAddAmountDetails.Email)
		if err != nil {
			fmt.Println("error updating balance")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message":"amount credited successfully"})
	})

	http.HandleFunc("/debit_amount",func(w http.ResponseWriter,r *http.Request){
		var debitAmountdetails DebitAmount
		err:=json.NewDecoder(r.Body).Decode(&debitAmountdetails)
		if err!=nil{
			fmt.Println("json debitamount error")
			return
		}
		var databaseDebitAmountDetails DebitAmount
		query1:=`SELECT password,balance FROM bank_account WHERE email=$1`
		err=db.QueryRow(query1,debitAmountdetails.Email).Scan(&databaseDebitAmountDetails.Password,&databaseDebitAmountDetails.Amount)
		if err!=nil{
			fmt.Println("query exec error2")
			return
		}
		if debitAmountdetails.Password!=databaseDebitAmountDetails.Password{
			fmt.Println("incorrect password")
		}
		var balance,amount int
		balance,err = strconv.Atoi(databaseDebitAmountDetails.Amount)
		if err!=nil{
			fmt.Println("error in converting balance")
			return
		}
		amount,err= strconv.Atoi(debitAmountdetails.Amount)
		if err!=nil{
			fmt.Println("error in converting amount")
			return
		}
	
		
		if balance < amount {
			json.NewEncoder(w).Encode(map[string]string{"message":"insufficient funds"})
			return
		}
		totalBalance:=strconv.Itoa(balance-amount)
		query2:= `UPDATE bank_account SET balance=$1 WHERE email=$2`
		_,err=db.Exec(query2,totalBalance,debitAmountdetails.Email)
			if err !=nil{
				fmt.Println("qurey update error")
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"mesage":"Amount debited Succesfuly"})
		})

	err = http.ListenAndServe(":8000" , nil)
	if err != nil {
		fmt.Println("server error",err)
	}
}
