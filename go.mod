module unscrabble

go 1.15

replace example.com/unscrabble => ../GoUnscrabble

require (
	example.com/unscrabble v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.4
	github.com/stretchr/testify v1.6.1
	golang.org/x/text v0.3.0
	golang.org/x/tools v0.0.0-20190425150028-36563e24a262
)
