package http

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/danhilltech/100kb.golang/pkg/utils"
)

type URLRequest struct {
	Url           string
	Domain        string
	LastAttemptAt int64
	Status        int64
	ContentType   string
	Etag          string
	LastModified  string

	Response *http.Response
}

var mu sync.Mutex

func getURLRequestFromDB(newUrl string, db *sql.DB) (*URLRequest, error) {
	mu.Lock()
	defer mu.Unlock()

	var url string
	var domain string
	var lastAttemptAt sql.NullInt64
	var status sql.NullInt64
	var contentType sql.NullString
	var etag sql.NullString
	var lastModified sql.NullString

	err := db.QueryRow("SELECT url, domain, lastAttemptAt, status, contentType, etag, lastModified FROM url_requests WHERE url = ?", newUrl).Scan(&url, &domain, &lastAttemptAt, &status, &contentType)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &URLRequest{
		Url:           url,
		Domain:        domain,
		LastAttemptAt: lastAttemptAt.Int64,
		Status:        status.Int64,
		ContentType:   contentType.String,
		Etag:          etag.String,
		LastModified:  lastModified.String,
	}, nil

}

func (urlRequest *URLRequest) Save(txn *sql.DB) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := txn.Exec(`INSERT INTO 
	url_requests(url, domain, lastAttemptAt, status, contentType, etag, lastModified) 
	VALUES(?, ?, ?, ?, ?, ?, ?) 
	ON CONFLICT(url) DO UPDATE SET domain=excluded.domain, lastAttemptAt=excluded.lastAttemptAt, status=excluded.status, contentType=excluded.contentType, etag=excluded.etag, lastModified=excluded.lastModified;`,
		urlRequest.Url,
		urlRequest.Domain,
		utils.NullInt64(urlRequest.LastAttemptAt),
		utils.NullInt64(urlRequest.Status),
		utils.NullString(urlRequest.ContentType),
		utils.NullString(urlRequest.Etag),
		utils.NullString(urlRequest.LastModified))

	return err

}
