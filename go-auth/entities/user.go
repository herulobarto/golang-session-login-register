package entities

type User struct {
	Id          int64
	NamaLengkap string `validate:"required" label:"Nama Lengkap"`
	Email       string `validate:"required,email"`
	Username    string `validate:"required,gte=3"`
	Password    string `validate:"required,gte=6"`
	Cpassword   string `validate:"required,eqfield=Password" label:"Konfirmasi Password"`
}
