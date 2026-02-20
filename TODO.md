# TODO

> [!NOTE]
> Install Extension TODO Tree for easly the checking todo list

- [x] Migrate to Hexagonal Architecture
- [ ] Implement reset password

## Reset Password (Link + SMTP) - Checklist

- [ ] Tambah usecase `RequestPasswordReset(ctx, email)` di service auth.
  - Validasi email.
  - Response generik (hindari user enumeration).
  - Generate raw token random, hash token, simpan hash ke `password_reset_tokens`.
  - Bangun reset URL dari env `FRONTEND_RESET_PASSWORD_URL`.
  - Kirim email via adapter SMTP (`EmailSender`).

- [ ] Tambah usecase `ConfirmPasswordReset(ctx, token, newPassword)`.
  - Validasi password baru.
  - Hash token, lookup token, cek expiry + used.
  - Update password user (`UpdateUserPassword`).
  - Mark token used.
  - (Disarankan) revoke seluruh refresh token user.

- [ ] Tambah endpoint HTTP:
  - `POST /auth/forgot-password`
  - `POST /auth/reset-password`

- [ ] Tambah DTO + validation di HTTP auth models:
  - `ForgotPasswordRequest { email }`
  - `ResetPasswordRequest { token, new_password }`

- [ ] Tambah SMTP config:
  - `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`
  - `SMTP_FROM_EMAIL`, `SMTP_FROM_NAME`, `SMTP_TLS`
  - `RESET_PASSWORD_TTL`, `FRONTEND_RESET_PASSWORD_URL`

- [ ] Tambah dokumentasi OpenAPI untuk forgot/reset password flow.

- [ ] Tambah test:
  - Unit test service reset password.
  - Handler test untuk validasi payload + status code.
