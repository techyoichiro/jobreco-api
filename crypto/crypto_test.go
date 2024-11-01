package crypto

import (
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// PasswordEncrypt のテスト
func TestPasswordEncrypt(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid password",
			args: args{
				password: "valid_password",
			},
			wantErr: false,
		},
		{
			name: "Empty password",
			args: args{
				password: "",
			},
			wantErr: false, // bcrypt handles empty passwords, so no error expected
		},
		{
			name: "Long password",
			args: args{
				password: "a very long password string that exceeds normal length",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PasswordEncrypt(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswordEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				err = bcrypt.CompareHashAndPassword([]byte(got), []byte(tt.args.password))
				if err != nil {
					t.Errorf("Generated hash does not match the original password")
				}
			}
		})
	}
}

// CompareHashAndPassword のテスト
func TestCompareHashAndPassword(t *testing.T) {
	password := "valid_password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	type args struct {
		hash     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid hash and password",
			args: args{
				hash:     string(hash),
				password: password,
			},
			wantErr: false,
		},
		{
			name: "Invalid password",
			args: args{
				hash:     string(hash),
				password: "wrong_password",
			},
			wantErr: true,
		},
		{
			name: "Empty password",
			args: args{
				hash:     string(hash),
				password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CompareHashAndPassword(tt.args.hash, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("CompareHashAndPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// EncryptEmail のテスト
func TestEncryptEmail(t *testing.T) {
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef") // ENCRYPTION_KEY を設定

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid email",
			args: args{
				email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "Empty email",
			args: args{
				email: "",
			},
			wantErr: false,
		},
		{
			name: "No ENCRYPTION_KEY",
			args: args{
				email: "test@example.com",
			},
			wantErr: true, // ENCRYPTION_KEY がない場合、エラーを期待
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ENCRYPTION_KEY のクリアと設定を切り替え
			if tt.wantErr {
				os.Setenv("ENCRYPTION_KEY", "")
			} else {
				os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
			}

			_, err := EncryptEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// DecryptEmail のテスト
func TestDecryptEmail(t *testing.T) {
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef") // ENCRYPTION_KEY を設定

	email := "test@example.com"
	encryptedEmail, _ := EncryptEmail(email)

	type args struct {
		encryptedEmail string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid encrypted email",
			args: args{
				encryptedEmail: encryptedEmail,
			},
			want:    email,
			wantErr: false,
		},
		{
			name: "Invalid encrypted email",
			args: args{
				encryptedEmail: "invalid_encrypted_email",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "No ENCRYPTION_KEY",
			args: args{
				encryptedEmail: encryptedEmail,
			},
			want:    "",
			wantErr: true, // ENCRYPTION_KEY がない場合、エラーを期待
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ENCRYPTION_KEY のクリアと設定を切り替え
			if tt.wantErr {
				os.Setenv("ENCRYPTION_KEY", "")
			} else {
				os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
			}

			got, err := DecryptEmail(tt.args.encryptedEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecryptEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
