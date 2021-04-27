package controllers

import "errors"

var (
	//ErrItemNotFound ...
	ErrItemNotFound = errors.New("item not found")
	//ErrEmptyRequestBody ...
	ErrEmptyRequestBody = errors.New("empty request body")
	//ErrEmptyItem ...
	ErrEmptyItem = errors.New("empty item")
)
