package card

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"log"
	"strconv"
)

type Transaction struct {
	XMLName xml.Name `xml:"transaction"`
	UserId int64 `json:"user_id" xml:"user_id"`
	Sum    int64 `json:"sum" xml:"sum"`
	MCC    string `json:"mcc" xml:"mcc"`
}

type Transactions struct {
	XMLName xml.Name `xml:"transactions"`
	Transactions []Transaction `xml:"transaction"`
}

func MakeTransactions(userId int64) []Transaction {
	const usersCount = 5
	const transactionsCount = 2
	const transactionAmount = 1_00
	transactions := make([]Transaction, usersCount*transactionsCount)
	x := Transaction{
		UserId: userId,
		Sum:    transactionAmount,
		MCC:    "5411",
	}
	y := Transaction{
		UserId: userId,
		Sum:    transactionAmount,
		MCC:    "5812",
	}
	z := Transaction{
		UserId: 2,
		Sum:    transactionAmount,
		MCC:    "5812",
	}

	for index := range transactions {
		switch index % 100 {
		case 0:
			transactions[index] = x
		case 20:
			transactions[index] = y
		default:
			transactions[index] = z
		}
	}
	return transactions
}


func ExportCSV(transactions []Transaction) ([]byte, error) {
	buffer := &bytes.Buffer{}
	w := csv.NewWriter(buffer)

	for _, t := range transactions {
		record := []string{
			strconv.FormatInt(t.UserId, 10),
			strconv.FormatInt(t.Sum, 10),
			TranslateMCC(t.MCC),
		}

		err := w.Write(record)
		if err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buffer.Bytes(), nil
}


func ExportJson(transactions []Transaction) ([]byte, error) {
	encoded, err := json.Marshal(transactions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return encoded, err
}

func ExportXML(transactions []Transaction) ([]byte, error) {
	encoded, err := xml.Marshal(Transactions{
		Transactions: transactions,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	encoded = append([]byte(xml.Header), encoded...)

	return encoded, nil
}