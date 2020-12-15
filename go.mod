module unscrabble

go 1.15

replace example.com/unscrabble => ../my_module

require (
	example.com/unscrabble v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1
)
