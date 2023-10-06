package cloner

/**
 * @Author AiTao
 * @Date 2023/6/15 10:48
 **/

type Cloneable interface {
	// Clone shallow clone creates a new object and copies the field values of the source object to the new object.
	//
	// If the source object contains reference type data, the shallow clone copies only the reference, not the reference object itself.
	//
	// That is, the source and clone objects will share the same reference object (only the top-level structure of the
	// object is copied, while the nested object is still a reference to the original object)
	Clone() any

	// DeepClone deep cloning creates a new object, which recursively copies the field values of the source object and all its
	// nested objects, ensuring that the cloned object is completely independent and does not have any shared reference objects.
	DeepClone() any
}
