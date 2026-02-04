# Documentation Rules

Aturan untuk membuat dokumentasi kode Go di folder `docs/internal/`.

---

## Struktur Dokumentasi

Setiap file dokumentasi harus mengikuti struktur berikut:

1. **Judul** - Nama file dan deskripsi singkat
2. **Overview** - Ringkasan isi file (fungsi/struct apa saja)
3. **Import** - Penjelasan setiap package yang di-import
4. **Penjelasan per Komponen** - Struct, fungsi, method
5. **Contoh Penggunaan** - Code examples yang praktis

---

## Format Penjelasan Kode

### Untuk Struct

```markdown
## Struct NamaStruct

\`\`\`go
type NamaStruct struct {
Field1 string
Field2 int
}
\`\`\`

### Penjelasan Field

**`Field1 string`**

- Deskripsi field
- Contoh value
- Default value (jika ada)
```

### Untuk Fungsi

```markdown
## Fungsi NamaFungsi()

\`\`\`go
func NamaFungsi(param1 string, param2 int) (string, error) {
// kode
}
\`\`\`

### Penjelasan Baris per Baris

**`func NamaFungsi(param1 string, param2 int) (string, error)`**

- `param1` - Deskripsi parameter
- `param2` - Deskripsi parameter
- Return: deskripsi return value

**`baris kode selanjutnya`**

- Penjelasan apa yang dilakukan baris ini
```

---

## Aturan Penulisan

### Yang Harus Dilakukan

- Jelaskan setiap baris kode yang penting
- Gunakan bahasa Indonesia yang mudah dipahami
- Sertakan contoh penggunaan yang praktis
- Jelaskan "kenapa" bukan hanya "apa"
- Gunakan code block dengan syntax highlighting (`go`)

### Yang Tidak Boleh Dilakukan

- Jangan gunakan tabel markdown (render jelek)
- Jangan skip penjelasan untuk kode yang kompleks
- Jangan gunakan istilah teknis tanpa penjelasan

---

## Contoh Penjelasan Konsep Go

Jika ada konsep Go yang perlu dijelaskan, gunakan format:

```markdown
### Penjelasan `konsep`

Penjelasan singkat konsep tersebut.

\`\`\`go
// Contoh kode
\`\`\`
```

Contoh konsep yang perlu dijelaskan:

- Bit shift (`1<<20`)
- Blank identifier (`_ =`)
- Type assertion
- Interface
- Pointer vs value
- Goroutine dan channel

---

## Template File Dokumentasi

```markdown
# package/file.go - Deskripsi Singkat

Package `nama` berisi deskripsi package.

---

## Overview

File ini menyediakan:

- `Fungsi1()` - Deskripsi
- `Fungsi2()` - Deskripsi
- `Struct1` - Deskripsi

---

## Import

\`\`\`go
import (
"package1"
"package2"
)
\`\`\`

- `package1` - Deskripsi kenapa dipakai
- `package2` - Deskripsi kenapa dipakai

---

## Struct/Fungsi 1

(penjelasan detail)

---

## Contoh Penggunaan

(code examples)
```

---

## Referensi File

Dokumentasi yang sudah ada sebagai referensi:

- `docs/internal/httpx-response.md` - Contoh dokumentasi helper functions
- `docs/internal/config.md` - Contoh dokumentasi configuration
