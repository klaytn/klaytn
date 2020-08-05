package kas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

var models = []interface{}{
	&Tx{},
}

type SuiteRepository struct {
	suite.Suite
	repo *repository
}

func setTestDatabase(t *testing.T, mysql *gorm.DB) {
	// Drop previous test database if possible.
	if err := mysql.Exec("DROP DATABASE IF EXISTS test").Error; err != nil {
		if !strings.Contains(err.Error(), "database doesn't exist") {
			t.Fatal("Unexpected error happened!", "err", err)
		}
	}
	// Create new test database.
	if err := mysql.Exec("CREATE DATABASE test DEFAULT CHARACTER SET UTF8").Error; err != nil {
		t.Fatal("Error while creating test database", "err", err)
	}
	// Use test database
	if err := mysql.Exec("USE test").Error; err != nil {
		t.Fatal("Error while setting test database", "err", err)
	}

	// Auto-migrate data model from model.DataModels
	if err := mysql.AutoMigrate(models...).Error; err != nil {
		t.Fatal("Error while auto migrating data models", "err", err)
	}
}

func (s *SuiteRepository) SetupSuite() {
	id := "root"
	password := "root"

	mysql, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/?parseTime=True", id, password))
	if err != nil {
		s.T().Log("Failed connecting to mysql", "id", id, "password", password, "err", err)
		s.T().Skip()
	}

	setTestDatabase(s.T(), mysql)
	s.repo = &repository{db: mysql}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SuiteRepository))
}
