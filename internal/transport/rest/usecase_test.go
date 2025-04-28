package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	_ "github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"redirectServer/internal/database/migrations"
	repo "redirectServer/internal/database/repo"
	"redirectServer/internal/domain"
	"redirectServer/internal/service"
	"redirectServer/internal/transport/dto/params"
)

// Helper function to setup a test database
func setupTestDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		"0.0.0.0",
		"postgres",
		"postgres",
		"postgres",
		"5432",
		"disable",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	db, err := setupTestDB()
	if err != nil {
		panic(err)
	}

	if err := migrations.NewMigration1(db).Up(); err != nil {
		panic(err)
	}

	fpRepo := repo.NewFingerprintRepo(db)
	linkRepo := repo.NewLinkRepo(db)
	fpService := service.NewFingerprintService(fpRepo, linkRepo)
	linkService := service.NewLinkService(linkRepo)
	fpHandler := NewFingerprintHandler(linkService, fpService)
	linkHandler := NewLinkHandler(linkService)

	MapFPRoutes(router, fpHandler)
	MapLinkRoutes(router, linkHandler)

	return router
}

func TestFPRoutes(t *testing.T) {
	router := setupRouter()

	linkId := CreateLink(router, t)

	fp := CreateFP(linkId, router, t)

	FindByLinkId(linkId, router, t)

	FindByFP(fp, router, t)
}

func CreateLink(router *gin.Engine, t *testing.T) domain.NanoID {
	w := httptest.NewRecorder()

	dto := params.CreateSalonInviteLink{SalonId: uuid.New()}

	fpDTO, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/link/salon", strings.NewReader(string(fpDTO)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response = domain.DirectLink{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return response.NanoId
}

func CreateFP(linkId domain.NanoID, router *gin.Engine, t *testing.T) (fp params.Fingerprint) {
	w := httptest.NewRecorder()

	dto := params.Fingerprint{
		Language:     "ru",
		Languages:    []string{"ru", "en", "fr"},
		Cores:        0,
		Memory:       0,
		ScreenWidth:  0,
		ScreenHeight: 0,
		ColorDepth:   0,
		PixelRatio:   0,
		TimeZone:     "GMT",
	}
	fpDTO, _ := json.Marshal(dto)
	url := fmt.Sprintf("/fingerprint/%s", linkId)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(fpDTO)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	return dto
}

func FindByLinkId(linkId domain.NanoID, router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/fingerprint/find/%s", linkId)
	req, _ := http.NewRequest("POST", url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func FindByFP(fp params.Fingerprint, router *gin.Engine, t *testing.T) {
	badIP := "1.2.3.4:8080"
	goodIP := ""

	tests := []struct {
		name  string
		fp    params.Fingerprint
		ip    string
		found bool
	}{
		{
			name: "Identical Fp. FP found",
			fp: params.Fingerprint{
				Language:     fp.Language,
				Languages:    fp.Languages,
				Cores:        fp.Cores,
				Memory:       fp.Memory,
				ScreenWidth:  fp.ScreenWidth,
				ScreenHeight: fp.ScreenHeight,
				ColorDepth:   fp.ColorDepth,
				PixelRatio:   fp.PixelRatio,
				TimeZone:     fp.TimeZone,
			},
			ip:    goodIP,
			found: true,
		},
		{
			name: "Bad Fp. FP not found",
			fp: params.Fingerprint{
				Language:     "fake",
				Languages:    []string{"fake"},
				Cores:        100000,
				Memory:       100000,
				ScreenWidth:  100000,
				ScreenHeight: 100000,
				ColorDepth:   100000,
				PixelRatio:   100000,
				TimeZone:     "fake",
			},
			ip:    badIP,
			found: false,
		},
		{
			name: "bad lang, tz. FP found",
			fp: params.Fingerprint{
				Language:     "fake",
				Languages:    fp.Languages,
				Cores:        fp.Cores,
				Memory:       fp.Memory,
				ScreenWidth:  fp.ScreenWidth,
				ScreenHeight: fp.ScreenHeight,
				ColorDepth:   fp.ColorDepth,
				PixelRatio:   fp.PixelRatio,
				TimeZone:     "fake",
			},
			ip:    goodIP,
			found: true,
		},
		{
			name: "bad swidth, sheight, pixelr. FP found",
			fp: params.Fingerprint{
				Language:     fp.Language,
				Languages:    fp.Languages,
				Cores:        fp.Cores,
				Memory:       fp.Memory,
				ScreenWidth:  10000000,
				ScreenHeight: 10000000,
				ColorDepth:   fp.ColorDepth,
				PixelRatio:   10000000,
				TimeZone:     fp.TimeZone,
			},
			ip:    goodIP,
			found: true,
		},
		{
			name: "good all. IP VPN. FP found",
			fp: params.Fingerprint{
				Language:     fp.Language,
				Languages:    fp.Languages,
				Cores:        fp.Cores,
				Memory:       fp.Memory,
				ScreenWidth:  fp.ScreenWidth,
				ScreenHeight: fp.ScreenHeight,
				ColorDepth:   fp.ColorDepth,
				PixelRatio:   fp.PixelRatio,
				TimeZone:     fp.TimeZone,
			},
			ip:    badIP,
			found: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			fpDTO, _ := json.Marshal(&tt.fp)
			req, _ := http.NewRequest("POST", "/fingerprint/find", strings.NewReader(string(fpDTO)))
			req.RemoteAddr = tt.ip
			router.ServeHTTP(w, req)

			if tt.found {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})
	}
}
