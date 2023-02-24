package gorm

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/backendmaster/simple_bank/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBClient struct {
	Client *gorm.DB
}

func (m *DBClient) Connect(config util.Config) {
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("can't not connect to database ")
		panic(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{})
	m.Client = gormDB
}
func (m *DBClient) Disconnect() {
	db, err := m.Client.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("db.DB() failed")
	}
	db.Close()
}

// type User struct {
// 	// gorm.Model
// 	Username          string
// 	HashedPassword    string
// 	FullName          string
// 	Email             string
// 	PasswordChangedAt time.Time
// 	CreatedAt         time.Time
// }

// func DoSomeThing(db *DBClient) {

// 	// db.Client.AutoMigrate(&User{})

// 	password := util.RandomString(10)
// 	hashPassword, _ := util.HashedPassword(password)
// 	user := User{
// 		Username:       "Jack",
// 		HashedPassword: hashPassword,
// 		FullName:       util.RandomOwnerName(),
// 		Email:          util.RandomEmail(),
// 	}
// 	db.Client.Create(&user)
// 	// fmt.Printf("id:%v,username:%v",result.)
// }
