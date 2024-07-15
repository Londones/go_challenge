package fixtures

import (
	"math/rand"
	"strconv"
	"time"

	"go-challenge/internal/models"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

func NewCatFixture(userID string, pictureIndex int, race models.Races) *models.Cats {
	names := []string{"Mittens", "Whiskers", "Shadow", "Bella", "Luna", "Simba"}
	sexes := []string{"Male", "Female"}
	colors := []string{"Black", "White", "Gray", "Orange", "Calico"}
	behaviors := []string{"Playful", "Lazy", "Aggressive", "Friendly", "Shy"}
	descriptions := []string{
		"Un chat très amical et joueur.",
		"Ce chat aime se prélasser au soleil.",
		"Très actif et adore chasser les jouets.",
		"Un peu timide mais très affectueux.",
		"Adore grimper et explorer.",
	}

	picturesURL := []string{
		"https://www.assuropoil.fr/wp-content/uploads/2023/07/avoir-un-chat-sante.jpg",
		"https://i.cbc.ca/1.7046192.1701492097!/fileImage/httpImage/african-wild-cat.jpg",
		"https://s.yimg.com/ny/api/res/1.2/C5uryMno9srLXTTHJWNllw--/YXBwaWQ9aGlnaGxhbmRlcjt3PTY0MDtoPTQ4MA--/https://s.yimg.com/os/en_US/News/BGR_News/funny-cat.jpg",
		"https://d2zp5xs5cp8zlg.cloudfront.net/image-53920-800.jpg",
	}

	birthDate := time.Now().AddDate(-rand.Intn(10), 0, 0)
	lastVaccineDate := time.Now().AddDate(-rand.Intn(5), 0, 0)
	description := randomChoice(descriptions)

	return &models.Cats{
		Name:            randomChoice(names),
		BirthDate:       &birthDate,
		Sexe:            randomChoice(sexes),
		LastVaccine:     &lastVaccineDate,
		LastVaccineName: "Rabies",
		Color:           randomChoice(colors),
		Behavior:        randomChoice(behaviors),
		Sterilized:      randomBool(),
		RaceID:          strconv.FormatUint(uint64(race.ID), 10),
		Description:     &description,
		Reserved:        randomBool(),
		PicturesURL:     pq.StringArray{picturesURL[pictureIndex%len(picturesURL)]},
		UserID:          userID,
	}
}

func CreateCatFixturesForUser(db *gorm.DB, count int, userID string) ([]*models.Cats, error) {
	var cats []*models.Cats
	var races []models.Races
	if err := db.Find(&races).Error; err != nil {

		return nil, err
	}
	for i := 0; i < count; i++ {
		cat := NewCatFixture(userID, i, races[i])
		if err := db.Create(cat).Error; err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}
