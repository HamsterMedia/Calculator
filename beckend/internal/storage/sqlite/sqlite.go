package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattn/go-sqlite3"

	"Calculator/internal/storage"
)

type Storage struct {
	db *sql.DB
}

// создаем базу из одной таблицы, в которой будем хранить все-все-все
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS history(
		id INTEGER PRIMARY KEY,
		val TEXT,
		otvet TEXT,
		insDate TIMESTAMP);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
//сораняем в таблицу значение и полученый ответ. А ДАТА сохранится сама
func (s *Storage) SaveURL(val string, otvet string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO history(insDate, val, otvet) VALUES(datetime('now'), ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
// бараба-а-а-ан-н-ная дробь.. отправляем инсерт в базу!
	res, err := stmt.Exec(val, otvet)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
//если ничего не сломалось, то база вернет нам ИД всттавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}
//достаем из таблицы одну запись по ИД
func (s *Storage) GetURL(id string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT id, insDate, val, otvet FROM history WHERE id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resId string
	var resDate string
	var resVal string
	var resOtvet string
// селект выполнился, заполняем переменные
	err = stmt.QueryRow(id).Scan(&resId, &resDate, &resVal, &resOtvet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}
	//да, возвращаем строкой, оркестратор превратит в json
	resUrl := "Долгожданный результат под Номером " + resId + " был посчитан " + resDate + ", считали: " + resVal + " и получили " + resOtvet

	return resUrl, nil
}
//возвращаем все-все-все запси из таблицы
func (s *Storage) GetAll() (string, error) {
type Calc struct {
    id int
    insDate string
    val string
    otvet string
}
//сортируем в обратном порядке, чтобы сверху были свежие записи
	resUrl := ""
	rows, err := s.db.Query("SELECT id, insDate, val, otvet FROM history ORDER BY id DESC")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w",  err)
	}

	calcs := []Calc{}
//селект выполнился, бежим по нему в цикле и бережно сохраняем все, что получили
	for rows.Next(){
	    c := Calc{}
	    err := rows.Scan(&c.id, &c.insDate, &c.val, &c.otvet )
	    if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w",  err)
	    }
	    calcs = append(calcs, c)
	}
	//теперь пробежимся по полученным данным и соберем все в одну строку
	//да, знаю, что коряво, но пока вот так
	for _, c := range calcs{
	  resUrl = resUrl + strconv.Itoa(c.id) + " | " + c.insDate + " | " + c.val + " | " + c.otvet + " | "+ "<br>"
	}
	

	return resUrl, nil
}

