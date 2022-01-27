Package intended for hashing new and checking exist user passwords

func HashPassword dedicated for creating password for encrypting received password string
using sha256 and returns hashed string or error if something went wrong

func CheckPasswordHash dedicated for compare received password and exist hash string
and returns boolean. If password and its hash are different would be returned false in another 
case would be returned true