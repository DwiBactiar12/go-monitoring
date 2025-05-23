package main

import (
	"fmt"
	"os"
)

func main() {
	// List folder yang akan dibuat
	folders := []string{
		"cmd",
		"config",
		"internal/domain/auth/entity",
		"internal/domain/auth/repository",
		"internal/domain/auth/usecase",
		"internal/domain/auth/handler",

		"internal/domain/log/entity",
		"internal/domain/log/repository",
		"internal/domain/log/usecase",
		"internal/domain/log/handler",

		"internal/domain/transaction/entity",
		"internal/domain/transaction/repository",
		"internal/domain/transaction/usecase",
		"internal/domain/transaction/handler",

		"internal/domain/metric/entity",
		"internal/domain/metric/repository",
		"internal/domain/metric/usecase",
		"internal/domain/metric/handler",

		"internal/domain/cache", // opsional

		"internal/middleware",
		"internal/server",

		"pkg/db",
		"pkg/jwt",
		"pkg/utils",

		"scripts",
	}

	// File kosong yang ingin dibuat (relative dari root)
	files := []string{
		"cmd/main.go",
		"config/config.go",
		"internal/domain/cache/cache.go",
		"internal/domain/cache/cache_usecase.go",
		"internal/middleware/jwt.go",
		"internal/middleware/logger.go",
		"internal/server/router.go",
		"pkg/db/postgres.go",
		"pkg/db/influx.go",
		"pkg/db/redis.go",
		"pkg/jwt/jwt.go",
		"pkg/utils/utils.go",
		"scripts/migrate.sh",
		".env",
		"README.md",
		"go.mod",
		"go.sum",
	}

	// Buat folder
	for _, folder := range folders {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			fmt.Printf("Gagal buat folder %s: %v\n", folder, err)
		} else {
			fmt.Printf("Folder dibuat: %s\n", folder)
		}
	}

	// Buat file kosong jika belum ada
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			f, err := os.Create(file)
			if err != nil {
				fmt.Printf("Gagal buat file %s: %v\n", file, err)
			} else {
				fmt.Printf("File dibuat: %s\n", file)
				f.Close()
			}
		} else {
			fmt.Printf("File sudah ada: %s\n", file)
		}
	}
}
