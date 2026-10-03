package storage

import (
	"database/sql"
	"os"
	"project/helper"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type UserStorage interface {
	CreateUser(id string, firstName string, lastName string, phone string) error
	CreateCar(id, model, car_type, user_id string) error
	GetUser(id string) (User, error)
	GetUsers() ([]User, error)
	UpdateUser(id string, firstName string, lastName string, phone string) error
	UpdateCar(id string, model string, car_type string, user_id string) error
	DeleteUser(id string) error
}

func NewStorage(db *sql.DB) UserStorage {
	return &Storage{DB: db}
}

type Storage struct {
	DB *sql.DB
}

func DB() (*sql.DB, error) {

	if err := godotenv.Load("internals/storage/.env"); err != nil {
		return nil, err
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}

func (d *Storage) CreateUser(
	id string,
	firstName string,
	lastName string,
	phone string,
) error {
	query := `INSERT INTO users (id, first_name, last_name, phone)
	VALUES($1, $2, $3, $4)`

	_, err := d.DB.Exec(query, id, firstName, lastName, phone)
	if err != nil {
		return err
	}

	return nil
}

func (d *Storage) CreateCar(id, model, car_type, user_id string) error {
	query := `INSERT INTO cars(id, model, car_type, user_id)
	VALUES($1, $2, $3, $4)`

	_, err := d.DB.Exec(query, id, model, car_type, user_id)

	if err != nil {
		return err
	}

	return nil
}

// INNER JOIN показывает только те строки, для которых есть совпадение в обеих таблицах

// SELF JOIN особенно полезен, когда строки одной таблицы связаны друг с другом

// CROSS JOIN = каждая строка × каждая строка

//LEFT JOIN говорит:

//«Покажи ВСЕ строки из левой таблицы, даже если справа нет совпадения

func (d *Storage) GetUser(id string) (User, error) {

	var carID sql.NullString
	var carModel sql.NullString
	var carType sql.NullString
	var carUserId sql.NullString

	query := `SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.phone,
	users.created_at,
	users.updated_at,
	cars.id,
	cars.model,
	cars.car_type,
	cars.user_id
	FROM users
	LEFT JOIN cars 
	ON users.id = cars.user_id
	WHERE users.id = $1
	`

	rows, err := d.DB.Query(query, id)

	if err != nil {
		return User{}, err
	}

	defer rows.Close()

	var user User

	for rows.Next() {

		var car Car

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.CreatedAt,
			&user.UpdatedAt,
			&carID,
			&carModel,
			&carType,
			&carUserId,
		); err != nil {
			return User{}, err
		}

		if carID.Valid {
			car.ID = carID.String
			car.Model = carModel.String
			car.CarType = carType.String
			car.UserId = carUserId.String
			user.Car = append(user.Car, car)
		}

	}

	if err := rows.Err(); err != nil {
		return User{}, err
	}

	return user, nil
}

func (d *Storage) GetUsers() ([]User, error) {

	var carID sql.NullString
	var carModel sql.NullString
	var carType sql.NullString
	var carUserId sql.NullString

	// здеь мы обьявляем sql стоя связять наш таблицу с этим
	// надо быть аккурантным с JOIN
	query := `SELECT * FROM users JOIN cars ON users.id = cars.user_id`

	// чтобы получить много строк мы испольщуем query not query.Row
	rows, err := d.DB.Query(query)
	if err != nil {
		return []User{}, err
	}
	// здесь нужна закрыть так как нам не нужна бесполезный открытый окно который память
	defer rows.Close()

	// мы обьявлем новый переменную который с типом юзер слайс
	var users []User

	// мы обьявлем новый переменную который с типом юзер
	var currentUser User

	// query() всегда нужна Next функция это делает итерацию
	for rows.Next() {

		// чтобы получить данные мы создаем два следующий перемен и с scan функцией нужна переменная
		var user User

		var car Car

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.CreatedAt,
			&user.UpdatedAt,
			&carID,
			&carModel,
			&carType,
			&carUserId); err != nil {
			return []User{}, err
		}

		if carID.Valid {
			car.ID = carID.String
			car.Model = carModel.String
			car.CarType = carType.String
			car.UserId = carUserId.String
		}

		if currentUser.ID == "" {
			currentUser = user
			currentUser.Car = append(currentUser.Car, car)

		} else {
			if currentUser.ID == user.ID {

				currentUser.Car = append(currentUser.Car, car)

			} else {
				users = append(users, currentUser)

				currentUser = user

				currentUser.Car = append(currentUser.Car, car)

			}

		}

	}

	users = append(users, currentUser)

	if err = rows.Err(); err != nil {
		return []User{}, err
	}

	return users, nil

}

func (d *Storage) UpdateUser(id string, firstName string, lastName string, phone string) error {

	query := `
	UPDATE users 
	SET first_name = $1,
		last_name = $2,
		phone = $3,
		updated_at = $4
	WHERE id = $5`

	updatedAt := time.Now()

	result, err := d.DB.Exec(query, firstName, lastName, phone, updatedAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return helper.UserNotFound
	}

	return nil
}

func (d *Storage) UpdateCar(id string, model string, car_type string, user_id string) error {

	query := `UPDATE cars
			  SET model = $1,
			  	car_type = $2
			  WHERE id = $3
			  AND user_id = $4`

	result, err := d.DB.Exec(query, model, car_type, id, user_id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return helper.CarNotFound
	}

	return nil

}

func (d *Storage) DeleteUser(id string) error {

	query1 := `DELETE FROM cars WHERE user_id = $1`

	query2 := `DELETE FROM users WHERE id = $1`

	_, err := d.DB.Exec(query1, id)
	if err != nil {
		return err
	}

	result, err := d.DB.Exec(query2, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return helper.UserNotFound
	}

	return nil

}
