package model

// SecurityQuestions are the questions a member answers when registering.
//
// The list matches the questions of the previous application, so an account
// created here can also be recovered there.
var SecurityQuestions = []string{
	"Apa nama panggilan masa kecil Anda?",
	"Di kota mana Anda dilahirkan?",
	"Apa nama sekolah dasar Anda?",
	"Apa nama hewan peliharaan pertama Anda?",
	"Apa nama jalan tempat Anda dibesarkan?",
	"Siapa nama teman masa kecil Anda?",
}

// RoleAnggota is the role the previous application gave to a member that
// registered itself.
const RoleAnggotaLegacy = "ANGGOTA"

// Registration describes a member account that wants to register itself.
//
// It is written to the USERS sheet of the previous application, so the same
// account can be used in both applications.
type Registration struct {
	// Username is the account name.
	Username string
	// Password is the plaintext password; it is hashed before it is stored.
	Password string
	// Nama is the display name of the member.
	Nama string
	// Kecamatan is the district of the member.
	Kecamatan string
	// Question1 and Question2 are the two security questions.
	Question1 string
	// Answer1 and Answer2 are the answers to the security questions.
	Answer1   string
	Question2 string
	Answer2   string
}
