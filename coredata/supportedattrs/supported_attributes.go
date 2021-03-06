package supportedattrs

import (
	"github.com/MarcGrol/userautomation/core/userattribute"
)

const (
	FirstName      = "first_name"
	FullName       = "full_name"
	EmailAddress   = "email_address"
	PhoneNumber    = "phone_number"
	DateBorn       = "date_born"
	Age            = "age"
	IsAdmin        = "is_admin"
	IsModerator    = "is_moderator"
	IsRegularUser  = "is_regular_user"
	LastLoginDate  = "last_login_date"
	LoginCount     = "login_count"
	LastPostDate   = "last_post_date"
	PostCount      = "post_count"
	MemberOfGroups = "member_of_groups"
)

var attributes = map[string]userattribute.UserAttributeSpec{
	FirstName: {
		Name:        FirstName,
		Description: "First name of a user",
		DataType:    userattribute.DataTypeString,
	},
	FullName: {
		Name:        FullName,
		Description: "Full name of a user",
		DataType:    userattribute.DataTypeString,
	},
	EmailAddress: {
		Name:        EmailAddress,
		Description: "Email address of a user",
		DataType:    userattribute.DataTypeEmailAddress,
	},
	PhoneNumber: {
		Name:        PhoneNumber,
		Description: "International phone number of a user",
		DataType:    userattribute.DataTypePhoneNumber,
	},
	DateBorn: {
		Name:        DateBorn,
		Description: "Date born",
		DataType:    userattribute.DataTypeDate,
	},
	Age: {
		Name:        Age,
		Description: "Current age of user",
		DataType:    userattribute.DataTypeInt,
	},
	IsAdmin: {
		Name:        IsAdmin,
		Description: "Indicates if a user is an admin",
		DataType:    userattribute.DataTypeBool,
	},
	IsModerator: {
		Name:        IsModerator,
		Description: "Indicates if a user is an admin",
		DataType:    userattribute.DataTypeBool,
	},
	IsRegularUser: {
		Name:        IsRegularUser,
		Description: "Indicates if a user is a regular user",
		DataType:    userattribute.DataTypeBool,
	},
	LastLoginDate: {
		Name:        LastLoginDate,
		Description: "Date oof last login",
		DataType:    userattribute.DataTypeDate,
	},
	LoginCount: {
		Name:        LoginCount,
		Description: "Amount of times that this user has logged in",
		DataType:    userattribute.DataTypeInt,
	},
	LastPostDate: {
		Name:        LastPostDate,
		Description: "Date oof last post",
		DataType:    userattribute.DataTypeDate,
	},
	PostCount: {
		Name:        PostCount,
		Description: "Amount of times that this user has posted a topic",
		DataType:    userattribute.DataTypeInt,
	},
	MemberOfGroups: {
		Name:        MemberOfGroups,
		Description: "The groups that this user is a member of",
		DataType:    userattribute.DataTypeStringSlice,
	},
}
