package userattribute

type DataType string

const (
	DataTypeBool         DataType = "bool"
	DataTypeInt                   = "int"
	DataTypeString                = "string"
	DataTypeStringSlice           = "string_slice"
	DataTypeEmailAddress          = "email_address"
	DataTypePhoneNumber           = "phone_number"
	DataTypeDate                  = "date"
)

type UserAttributeSpec struct {
	Name        string
	Description string
	DataType    DataType
}
