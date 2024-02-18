package storage

//оно тут просто так. так надо. когда-нибудь здесь будет нормальный обработчик ошибок.
import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)
