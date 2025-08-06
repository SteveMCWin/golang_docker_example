package persons

import (
	"os"
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"log"
)

var person_db_opened = false

type Db struct {
	db *sql.DB
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IpAddress string `json:"ip_address"`
}

func (base *Db) InitDb() error {
	if person_db_opened == true {
		log.Println("WARNING: tried to open person.db more than once")
		return nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir("./persons/")
    if err != nil {
        log.Fatal(err)
    }
 
    for _, e := range entries {
            log.Println(e.Name())
    }

	spellfix_relative_path := "/persons/spellfix.so"

	log.Println("spellfix full path:", dir+spellfix_relative_path)

	sql.Register("sqlite3_with_extension",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				dir + spellfix_relative_path,
			},
		},
	)

	db_path := "persons/person.db"
	base.db, err = sql.Open("sqlite3_with_extension", db_path)
	if err != nil {
		return err
	}

	person_db_opened = true

	return nil
}

func (base *Db) GetPersons() ([]Person, error) {
	rows, err := base.db.Query("select id, first_name, last_name, email, ip_address from people")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singular_person := Person{}

		err = rows.Scan(
			&singular_person.Id,
			&singular_person.FirstName,
			&singular_person.LastName,
			&singular_person.Email,
			&singular_person.IpAddress,
		)

		if err != nil {
			return nil, err
		}

		people = append(people, singular_person)
	}

	return people, nil
}

func (base *Db) GetPersonById(id string) (Person, error) {
	stmt, err := base.db.Prepare("select id, first_name, last_name, email, ip_address from people where id=?")
	if err != nil {
		return Person{}, err
	}

	p := Person{}

	sqlErr := stmt.QueryRow(id).Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.IpAddress)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Person{}, nil
		}
		return Person{}, sqlErr
	}

	return p, nil
}

func (base *Db) AddPerson(newPerson Person) (bool, error) {
	tx, err := base.db.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("insert into people (first_name, last_name, email, ip_address) values (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.IpAddress)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (base *Db) UpdatePerson(newPerson Person, id int) (bool, error) {
	tx, err := base.db.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.IpAddress, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (base *Db) DeletePerson(id int) (bool, error) {
	tx, err := base.db.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from people where id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (base *Db) FindPeopleByName(name string) ([]*Person, error) {

	rows, err := base.db.Query("select id, first_name from spellfix_people inner join people on word = first_name where word match ?", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matches := make([]*Person, 0)

	for rows.Next() {
		person := &Person{}
		err = rows.Scan(&person.Id, &person.FirstName)
		if err != nil {
			return nil, err
		}

		matches = append(matches, person)
	}

	return matches, nil
}
