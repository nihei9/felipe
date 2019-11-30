package definitions

import "errors"

var (
	errorVersionIsMissing                   = errors.New("`version` must be specified")
	errorKindIsMissing                      = errors.New("`kind` must be specified")
	errorKindIsNotComponents                = errors.New("`kind` must be `components`")
	errorKindIsNotFaces                     = errors.New("`kind` must be `faces`")
	errorComponentsHasNoComponent           = errors.New("`components` must contain at least one content")
	errorComponentsHasEmptyComponent        = errors.New("`components[]` includes empty components")
	errorComponentNameIsMissing             = errors.New("`components[].name` must be specified")
	errorComponentHasEmptyDependency        = errors.New("`dependencies[]` includes empty components")
	errorDependencyNameIsMissing            = errors.New("`dependencies[].name` must be specified")
	errorFacesHasNoFace                     = errors.New("`faces` must contain at least one face")
	errorFacesHasEmptyFace                  = errors.New("`faces[]` includes empty faces")
	errorFaceTargetIsMissing                = errors.New("`faces[].targets` is not specified")
	errorFaceTargetIsEmpty                  = errors.New("`faces[].targets` is empty")
	errorFaceMatchLabelsTargetHasEmptyEntry = errors.New("`faces[].targets.match_labels[]` includes empty entries")
	errorFaceAttributesHasNoAttribute       = errors.New("`faces[].attributes[]` must contain at least one attribute")
	errorFaceAttributesHasEmptyAttribute    = errors.New("`faces[].attributes[]` includes empty attributes")
)
