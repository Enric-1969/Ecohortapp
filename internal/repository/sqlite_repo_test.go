package repository

import (
	"ecohortapp/internal/domain"
	"testing"
	"time"
)

func TestSQLiteRepository_CRUD(t *testing.T) {
	// Usamos SQLite en memoria RAM (:memory:) para no crear archivos .db en el disco
	t.Setenv("DB_PATH", ":memory:")

	cfg := &domain.Config{}
	sqlDB, err := ConnectSQL(cfg)
	if err != nil {
		t.Fatalf("ConnectSQL falló con driver modernc.org/sqlite: %v", err)
	}
	defer sqlDB.Close()

	// 1. Verificar que SetupDB y Migrate() crean las tablas sin errores
	SetupDB(cfg, sqlDB)
	if cfg.DB == nil {
		t.Fatal("SetupDB no inicializó cfg.DB")
	}

	// 2. Probar INSERT (Crear un registro meteorológico)
	nuevo := domain.Registres{
		Data:         time.Now().Truncate(time.Second),
		Precipitacio: 12,
		TempMaxima:   24,
		TempMinima:   14,
		Humitat:      65,
	}

	insertado, err := cfg.DB.InsertRegistre(nuevo)
	if err != nil {
		t.Fatalf("InsertRegistre falló: %v", err)
	}
	if insertado.ID == 0 {
		t.Error("Se esperaba un ID autoincremental mayor que 0")
	}

	// 3. Probar SELECT (Leer el registro recién guardado)
	leido, err := cfg.DB.ObtenirRegistrePerID(int(insertado.ID))
	if err != nil {
		t.Fatalf("ObtenirRegistrePerID falló: %v", err)
	}
	if leido.TempMaxima != 24 || leido.Humitat != 65 {
		t.Errorf("Datos leídos incorrectos: se esperaba TempMax=24 Hum=65, se obtuvo TempMax=%d Hum=%d", leido.TempMaxima, leido.Humitat)
	}

	// 4. Probar DELETE (Borrar el registro y comprobar que ya no existe)
	if err := cfg.DB.BorrarRegistre(insertado.ID); err != nil {
		t.Fatalf("BorrarRegistre falló: %v", err)
	}
	if _, err := cfg.DB.ObtenirRegistrePerID(int(insertado.ID)); err == nil {
		t.Error("El registro borrado no debería existir en la base de datos")
	}
}
