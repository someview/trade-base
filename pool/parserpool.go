package pool

import "github.com/valyala/fastjson"

var fastParsePool = &fastjson.ParserPool{}

func GetParser() *fastjson.Parser {
	return fastParsePool.Get()
}
func PutParser(p *fastjson.Parser) {
	fastParsePool.Put(p)
}
