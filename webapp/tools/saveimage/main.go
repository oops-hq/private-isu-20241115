package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const imageFilePath = "../../public/image"

var (
	db *sqlx.DB
)

func saveImage(id string, blob []byte, mime string) (string, error) {
	ext := ""
	switch mime {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	}

	name := path.Join(imageFilePath, id+ext)
	file, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(blob)
	if err != nil {
		return "", err
	}

	return name, nil
}

func main() {
	host := os.Getenv("ISUCONP_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("ISUCONP_DB_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to read DB port number from an environment variable ISUCONP_DB_PORT.\nError: %s", err.Error())
	}
	user := os.Getenv("ISUCONP_DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("ISUCONP_DB_PASSWORD")
	if password == "" {
		password = "root"
	}
	dbname := os.Getenv("ISUCONP_DB_NAME")
	if dbname == "" {
		dbname = "isuconp"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)

	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s.", err.Error())
	}
	defer db.Close()

	var images []struct {
		ID      int    `db:"id"`
		Imgdata []byte `db:"imgdata"`
		Mime    string `db:"mime"`
	}
	err = db.Select(&images, "SELECT id, imgdata, mime FROM posts")
	if err != nil {
		log.Fatalf("failed to select images: %s", err.Error())
	}
	for _, img := range images {
		log.Printf("Processing image %d", img.ID)
		fileName, err := saveImage(strconv.Itoa(img.ID), img.Imgdata, img.Mime)
		if err != nil {
			log.Fatalf("failed to save image: %s", err.Error())
		}
		fn := path.Base(fileName)
		log.Printf("Saved image %d as %s", img.ID, fn)
		db.Exec("UPDATE posts SET image_file_name=? WHERE id = ?", fn, img.ID)
	}
}
