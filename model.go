package main

type record struct {
	LineNo    int
	Ip        string
	Domain    string
	IsDisplay bool
	IsComment bool
	Valid     bool
	Raw       string
}

type records []*record
