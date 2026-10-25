package repository

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"ecohortapp/internal/domain"

	_ "://github.com"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

// NewSQLiteRepository inicializa la estructura con la conexión activa
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

// ConnectSQL gestiona la connexió física amb la base de dades SQLite usando la ruta del dominio
func ConnectSQL(app *domain.Config) (*sql.DB, error) {
	path := ""

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		if app.InfoLog != nil {
			app.InfoLog.Println("db in:", path)
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SetupDB inicialitza el repositori en la configuración global i executa les migracions
func SetupDB(app *domain.Config, sqlDB *sql.DB) {
	app.DB = NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		if app.ErrorLog != nil {
			app.ErrorLog.Println(err)
		}
		log.Panic(err)
	}
}

func (repo *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS registres(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		data_registre INTEGER NOT NULL,
		precipitacio INTEGER NOT NULL,
		temp_maxima INTEGER NOT NULL,
		temp_minima INTEGER NOT NULL,
		humitat INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS prediccions (
		codi_ine TEXT NOT NULL,
		data_prediccio INTEGER NOT NULL,
		prob_precipitacio INTEGER NOT NULL DEFAULT 0,
		temp_maxima INTEGER NOT NULL DEFAULT 0,
		temp_minima INTEGER NOT NULL DEFAULT 0,
		humitat_relativa INTEGER NOT NULL DEFAULT 0,
		created_at INTEGER NOT NULL,
		PRIMARY KEY (codi_ine, data_prediccio)
	);
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertRegistre(registres domain.Registres) (*domain.Registres, error) {
	stmt := "insert into registres (data_registre, precipitacio, temp_maxima, temp_minima, humitat) values (?,?,?,?,?)"

	res, err := repo.Conn.Exec(stmt, registres.Data.Unix(), registres.Precipitacio, registres.TempMaxima, registres.TempMinima, registres.Humitat)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	registres.ID = id
	return &registres, nil
}

func (repo *SQLiteRepository) ObtenirTotsRegistres() ([]domain.Registres, error) {
	query := "select id, data_registre, precipitacio, temp_maxima, temp_minima, humitat from registres order by data_registre"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tots []domain.Registres
	for rows.Next() {
		var h domain.Registres
		var unixTime int64
		err := rows.Scan(
			&h.ID,
			&unixTime,
			&h.Precipitacio,
			&h.TempMaxima,
			&h.TempMinima,
			&h.Humitat,
		)
		if err != nil {
			return nil, err
		}
		h.Data = time.Unix(unixTime, 0)
		tots = append(tots, h)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tots, nil
}

func (repo *SQLiteRepository) ObtenirRegistrePerID(id int) (*domain.Registres, error) {
	row := repo.Conn.QueryRow("select id, data_registre, precipitacio, temp_maxima, temp_minima, humitat from registres where id = ?", id)

	var h domain.Registres
	var unixTime int64
	err := row.Scan(
		&h.ID,
		&unixTime,
		&h.Precipitacio,
		&h.TempMaxima,
		&h.TempMinima,
		&h.Humitat,
	)
	if err != nil {
		return nil, err
	}

	h.Data = time.Unix(unixTime, 0)
	return &h, nil
}

func (repo *SQLiteRepository) ActualitzarRegistre(id int64, actualitzar domain.Registres) error {
	if id == 0 {
		return errors.New("La id a actualitzar és incorrecte")
	}

	stmt := "Update registres set data_registre = ?, precipitacio = ?, temp_maxima = ?, temp_minima = ?, humitat = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, actualitzar.Data.Unix(), actualitzar.Precipitacio, actualitzar.TempMaxima, actualitzar.TempMinima, actualitzar.Humitat, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrUpdateFailed
	}

	return nil
}

func (repo *SQLiteRepository) BorrarRegistre(id int64) error {
	res, err := repo.Conn.Exec("delete from registres where id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrDeleteFailed
	}

	return nil
}

func (repo *SQLiteRepository) GuardarPrediccioDiaria(codiINE string, dataPrediccio int64, probPrecipitacio, tempMax, tempMin, humitat int) error {
	stmt := `
	INSERT INTO prediccions (codi_ine, data_prediccio, prob_precipitacio, temp_maxima, temp_minima, humitat_relativa, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(codi_ine, data_prediccio) DO UPDATE SET
		prob_precipitacio = excluded.prob_precipitacio,
		temp_maxima = excluded.temp_maxima,
		temp_minima = excluded.temp_minima,
		humitat_relativa = excluded.humitat_relativa,
		created_at = excluded.created_at;
	`
	_, err := repo.Conn.Exec(stmt, codiINE, dataPrediccio, probPrecipitacio, tempMax, tempMin, humitat, time.Now().Unix())
	return err
}
