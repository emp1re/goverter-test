package generator

import (
	"github.com/emp1re/goverter-test/xtype"
)

func validateMethods(lookup map[xtype.Signature]*generatedMethod) error {
	//fmt.Println("lookup", lookup)
	//for _, genMethod := range lookup {
	//	if genMethod.Explicit && len(genMethod.RawFieldSettings) > 0 {
	//		isTargetStructPointer := genMethod.Target.Pointer && genMethod.Parameters.Target.PointerInner.Struct
	//		if !genMethod.Target.Struct && !isTargetStructPointer {
	//			return fmt.Errorf("Invalid struct field mapping on method:\n    %s\n\nField mappings like goverter:map or goverter:ignore may only be set on struct or struct pointers.\nSee https://goverter.jmattheis.de/guide/configure-nested", genMethod.ID)
	//		}
	//	}
	//}
	return nil
}
