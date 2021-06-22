package swag


/**
   执行 swag init --generalInfo .\routers\http.go 生成docs
   https://www.ctolib.com/swaggo-swag.html
 */
type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}


