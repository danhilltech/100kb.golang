package http

import (
	"database/sql"
	"net/http"

	"github.com/danhilltech/100kb.golang/pkg/utils"
)

type URLRequest struct {
	Url           string
	Domain        string
	LastAttemptAt int64
	Status        string
	ContentType   string

	Response *http.Response
}

func getURLRequestFromDB(newUrl string, txn *sql.Tx) (*URLRequest, error) {
	var url string
	var domain string
	var lastAttemptAt sql.NullInt64
	var status sql.NullString
	var contentType sql.NullString

	err := txn.QueryRow("SELECT url, domain, lastAttemptAt, status, contentType FROM url_requests WHERE url = ?", newUrl).Scan(&url, &domain, &lastAttemptAt, &status, &contentType)

	if err != nil {
		return nil, err
	}

	return &URLRequest{
		Url:           url,
		Domain:        domain,
		LastAttemptAt: lastAttemptAt.Int64,
		Status:        status.String,
		ContentType:   status.String,
	}, nil

}

func (urlRequest *URLRequest) Save(txn *sql.Tx) error {
	_, err := txn.Exec(`INSERT INTO 
	url_requests(url, domain, lastAttemptAt, status, contentType) 
	VALUES(?, ?, ?, ?, ?) 
	ON CONFLICT(url) DO UPDATE SET domain=excluded.domain, lastAttemptAt=excluded.lastAttemptAt, status=excluded.status, contentType=excluded.contentType;`,
		urlRequest.Url,
		urlRequest.Domain,
		utils.NullInt64(urlRequest.LastAttemptAt),
		utils.NullString(urlRequest.Status),
		utils.NullString(urlRequest.ContentType))

	return err

}
