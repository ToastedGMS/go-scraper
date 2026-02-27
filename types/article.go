package types

type Article struct {
	Title, URL, Img, Date, Source string
}

type Story struct {
	Headline string    `json:"headline"`
	Articles []Article `json:"article"`
}
