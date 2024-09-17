package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Zelvalna/go_final_project/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

// InitDB инициализирует соединение с базой данных и создает таблицу, если она не существует
func InitDB() (*sqlx.DB, error) {
	dbFile := "./scheduler.db"

	// Получаем путь к базе данных из переменной окружения
	dbPath := os.Getenv("TODO_DBFILE")
	if len(dbPath) > 0 {
		dbFile = dbPath
	}
	// Вывод пути к базе данных для проверки
	log.Println("Путь к базе данных:", dbFile)

	// Проверка наличия файла базы данных
	_, err := os.Stat(dbFile)
	if err != nil && os.IsNotExist(err) {
		// Файл базы данных не существует, создаем новый
		file, err := os.Create(dbFile)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	// Открываем соединение с базой данных
	db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу, если она не существует
	err = createTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// createTable создает таблицу `scheduler`, если она не существует, и индекс по дате
func createTable(db *sqlx.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT(128)
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
        `)

	return err
}

// InsertTask добавляет новую задачу в базу данных
func InsertTask(task model.Task) (int, error) {
	if db == nil {
		return 0, errors.New("database not initialized")
	}
	// Вставляем задачу в таблицу

	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// ReadTasks читает все задачи из базы данных, ограничивая результат 15 записями
func ReadTasks() ([]model.Task, error) {
	var tasks []model.Task

	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date")
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

// SearchTasks ищет задачи по заголовку или комментарию
func SearchTasks(search string) ([]model.Task, error) {
	var tasks []model.Task

	search = fmt.Sprintf("%%%s%%", search)
	rows, err := db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date",
		sql.Named("search", search))
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

// SearchTasksByDate ищет задачи по дате
func SearchTasksByDate(date string) ([]model.Task, error) {
	var tasks []model.Task

	rows, err := db.Query("SELECT * FROM scheduler WHERE date = :date",
		sql.Named("date", date))
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

// ReadTaskById читает задачу по ID
func ReadTaskById(id string) (model.Task, error) {
	var task model.Task

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// UpdateTask обновляет задачу по ID
func UpdateTask(task model.Task) (model.Task, error) {

	result, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return model.Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}

	if rowsAffected == 0 {
		return model.Task{}, errors.New("failed to update")
	}

	return task, nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete")
	}

	return err
}
