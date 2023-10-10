package clone

/**
 *
 * @Author AiTao
 * @Date 2023/10/7 4:42
 * @Url
 **/

const (
	Invalid          OpFlags = iota
	OnlyPublicField  OpFlags = 1 << (iota - 1) // only copy the exported fields.
	OnlyPrivateField                           // only copy fields are not exported.
	DeepString
	DeepFunc
	DeepArray
	AllFields = OnlyPublicField | OnlyPrivateField // copy all fields.
)

type OpFlags uint

func (o *OpFlags) Has(flags OpFlags) bool {
	return (*o & flags) != 0
}

func (o *OpFlags) Add(flags OpFlags) OpFlags {
	*o = *o | flags
	return *o
}

func (o *OpFlags) Clear(flags OpFlags) OpFlags {
	*o = *o & ^flags
	return *o
}
