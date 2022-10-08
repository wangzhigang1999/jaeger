package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"os"
	"testing"
)

func Test_parserUrl(t *testing.T) {
	println(parseUrl("http://10.1.62.178:30001/customers/iumpng6WmrGbIuszE2-cAK8M81PdmW52"))
	println(parseUrl("http://catalogue/catalogue/634053fbf5754700071142fc"))
	println(parseUrl("http://orders/orders"))
}

func Benchmark_parserUrl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseUrl("http://10.1.62.178:30001/customers/iumpng6WmrGbIuszE2-cAK8M81PdmW52")
	}
}

func readJson() []model.Span {
	filePtr, err := os.Open("span.json")
	if err != nil {
		fmt.Printf("文件打开失败 [Err:%s]\n", err.Error())
		return nil
	}
	defer filePtr.Close()
	var info []model.Span
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		return nil
	}

	return info
}
