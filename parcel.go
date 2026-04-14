package main

import (
	"database/sql"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	//добавление данных с помощью запроса insert
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		log.Println("не удалось добавить данные в таблицу")
		return 0, err

	}

	//получаем иддентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("не удалось получить идентификатор последней добавленной записи")
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	//создание экмземпляра структуры
	p := Parcel{}

	// получение строки в соответствии с запросом
	row := s.db.QueryRow("SELECT number, client, status, address, created_at from parcel where number = :number", sql.Named("number", number))

	//заполнение экземпляра структуры полученными данными
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		log.Println("не удалось записать данные по параметру number")
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	//объявление слайса структур, куда будут записываться данные
	var res []Parcel

	//получение строк в соответствии с запросом
	rows, err := s.db.Query("SELECT number, client, status, address, created_at from parcel where client = :client", sql.Named("client", client))
	if err != nil {
		log.Println("не удалось получить данные по параметру Client")
		return []Parcel{}, err
	}

	defer rows.Close()

	//итерируемся по этим строкам
	for rows.Next() {
		//создали экземпляр структуры Parcel
		p := Parcel{}

		//записываем в структуру Parcel полученные данные
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			log.Println("не удалось записать данные в структуру Parcel")
			return nil, err
		}

		//заполняем срез такими структурами с помощью append
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		log.Println("ошибка при итерации")
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	//обновляем статус с помощью функции update
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		log.Println("не удалось обновить статус по заданному номеру")
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	//упрощение логики
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		log.Println("не удалось обновить адрес")
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	//упрощение логики
	_, err := s.db.Exec("DELETE FROM parcel where number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		log.Println("не удалость удалить строку")
		return err
	}

	return nil
}
