package engine

type Query struct {
	QueryString string `form:"q"`
	Page        int    `form:"page"`
}
