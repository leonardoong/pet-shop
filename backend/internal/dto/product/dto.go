package product

type Filter struct {
	CategorySlug string
	Search       string
	MinPrice     *float64
	MaxPrice     *float64
	Sort         string // price_asc | price_desc | newest | oldest
	Page         int
	Limit        int
}
