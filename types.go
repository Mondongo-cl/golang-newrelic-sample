package main

type Configuration struct {
	username string
	password string
	port     int
	host     string
	database string
}

type Customer struct {
	Id            int64   `db:"id"`
	CreateAt      []uint8 `db:"created_at"`
	UpdateAt      []uint8 `db:"updated_at"`
	DeleteAt      []uint8 `db:"deleted_at"`
	FirstName     string  `db:"first_name"`
	LastName      string  `db:"last_name"`
	Email         string  `db:"email"`
	DateOfBirth   []uint8 `db:"dob"`
	CountryCode   string  `db:"country_code"`
	MovilNumber   string  `db:"mobile_number"`
	ProfilePicId  string  `db:"profile_pic_res_id"`
	BLId          string  `db:"bl_id"`
	BLQR          string  `db:"bl_qr"`
	RadixxId      string  `db:"radixx_id"`
	IOMobId       string  `db:"iomob_id"`
	IsVerified    string  `db:"is_verified"`
	MarkettingCom string  `db:"marketting_com"`
	Password      string  `db:"password"`
	MinorConcent  *string `db:"minor_consent"`
}
