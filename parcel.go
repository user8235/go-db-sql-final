package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	//_, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :created_at)",
	result, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES ( :client, :status, :address, :created_at)",
		//sql.Named("number,", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Println("Error during 'Add' operation:", err)
		return 0, nil
	}
	// верните идентификатор последней добавленной записи
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}

	return int(lastID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = :id",
		sql.Named("id", number))
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println("Error during 'Get' operation:")
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :id", sql.Named("id", client))
	if err != nil {
		fmt.Println("Error during 'GetByClient' operation:", err)
		return nil, err
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var (
			num   int
			cl    int
			stat  string
			addr  string
			creat string
		)

		err := rows.Scan(&num, &cl, &stat, &addr, &creat)
		if err != nil {
			fmt.Println("Error during 'GetByClient' operation:", err)
			return nil, err
		}
		res = append(res, Parcel{Number: num, Client: cl, Status: stat, Address: addr, CreatedAt: creat})
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :id",
		sql.Named("status", status),
		sql.Named("id", number))
	if err != nil {
		fmt.Println("Error during 'SetStatus' operation:", err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	// var status string
	// err := s.db.QueryRow("SELECT status FROM parcel WHERE number = :id", sql.Named("id", number))//.Scan(&address)
	// if err != nil {
	// 	fmt.Println("Error during 'SetAddress' operation:", err)
	// 	return err
	// }
	// if status != ParcelStatusRegistered {
	// 	fmt.Println("Parcel status is NOT registered!")
	// 	return nil
	// }

	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :id AND status = :status",
		sql.Named("address", address),
		sql.Named("id", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		fmt.Println("Error during 'SetAddress' operation:", err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var status string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE id = :id", sql.Named("id", number)).Scan(&status)
	if err != nil {
		fmt.Println("Error during 'Delete' operation:", err)
		return err
	}
	if status != ParcelStatusRegistered {
		fmt.Println("Parcel status is NOT registered! Delet is deny!")
		return nil
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE id = :id", sql.Named("id", number))
	if err != nil {
		fmt.Println("Error during 'Delete' operation:", err)
		return err
	}
	return nil
}
