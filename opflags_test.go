package cloner

import (
	"fmt"
	"testing"
)

/**
 *
 * @Author AiTao
 * @Date 2023/10/7 4:42
 * @Url
 **/

func TestOpFlags_Has(t *testing.T) {
	f1 := AllFields
	// All:     1 1
	// Public:  0 1
	// Private: 1 0
	fmt.Println(f1.Has(OnlyPrivateField))
	fmt.Println(f1.Has(OnlyPublicField))
	fmt.Println(f1.Has(AllFields))

	fmt.Println(f1 == OnlyPrivateField)
	fmt.Println(f1 == OnlyPublicField)
	fmt.Println(f1 == AllFields)

	f1.Clear(OnlyPrivateField)
	fmt.Println(f1)
}
