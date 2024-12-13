package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type database struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

type server struct {
	Address string
	Port    string
}

var (
	Debug                        bool
	Server                       server
	Database                     database
	CryptoKey                    string
	JwtSecretKey                 string
	DefaultDataPerPage           int
	AuthTokenExpiredMinute       int
	FileBaseUrl                  string
	StorageDirectory             string
	PhotoDirectory               string
	PhotoincRunningLimit         int64
	PhotoUploadMaxSizeByte       int64
	PhotoUploadAllowedExtensions []string
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Failed load env Err: " + err.Error())
		panic(err)
	}

	Debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		fmt.Println("Failed parse DEBUG Err: " + err.Error())
		panic(err)
	}
	Server = server{
		Address: os.Getenv("SERVER_ADDRESS"),
		Port:    os.Getenv("SERVER_PORT"),
	}

	Database = database{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	hasher := md5.New()
	hasher.Write([]byte(os.Getenv("CRYPTO_KEY")))
	CryptoKey = hex.EncodeToString(hasher.Sum(nil))

	hasher.Write([]byte(os.Getenv("JWT_SECRET_KEY")))
	JwtSecretKey = hex.EncodeToString(hasher.Sum(nil))

	DefaultDataPerPage, err = strconv.Atoi(os.Getenv("DEFAULT_DATA_PER_PAGE"))
	if err != nil {
		fmt.Println("Failed parse DEFAULT_DATA_PER_PAGE Err: " + err.Error())
		panic(err)
	}

	AuthTokenExpiredMinute, err = strconv.Atoi(os.Getenv("AUTH_TOKEN_EXPIRED_MINUTES"))
	if err != nil {
		fmt.Println("Failed parse AUTH_TOKEN_EXPIRED_MINUTES Err: " + err.Error())
		panic(err)
	}

	FileBaseUrl = os.Getenv("FILE_BASE_URL")
	StorageDirectory = os.Getenv("STORAGE_DIRECTORY")
	PhotoDirectory = os.Getenv("PHOTO_DIRECTORY")

	PhotoincRunningLimit, err = strconv.ParseInt(os.Getenv("PHOTOINC_RUNNING_LIMIT"), 10, 64)
	if err != nil {
		fmt.Println("Failed parse PHOTOINC_RUNNING_LIMIT Err: " + err.Error())
		panic(err)
	}

	PhotoUploadMaxSizeByte, err = strconv.ParseInt(os.Getenv("PHOTO_UPLOAD_MAX_SIZE_BYTE"), 10, 64)
	if err != nil {
		fmt.Println("Failed parse PHOTO_UPLOAD_MAX_SIZE_BYTE Err: " + err.Error())
		panic(err)
	}

	PhotoUploadAllowedExtensions = strings.Split(os.Getenv("PHOTO_UPLOAD_ALLOWED_EXTENSIONS"), ",")

}
