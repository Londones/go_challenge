package queries

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-challenge/internal/database"
	"go-challenge/internal/models"
	"strconv"
	"testing"
	"time"
)

var dtb, dbErr = database.TestDatabaseInit()
var db = DatabaseService{s: database.Service{Db: dtb.Db}}

// Mise en place d'objet globaux qui serviront aux tests
var annonceDescription = "annonce de test"
var associationIsVerified = true
var descriptionCat string = "Il s'appel PABLO et c'est un chat de TEST"
var dateCat = time.Time{}

var users, _ = db.GetAllUsers()
var raceUnique, _ = db.FindRaceByID("1")
var annonces, _ = db.GetAllAnnonces()

var annonce = models.Annonce{
	Title:       "Annonce de test",
	Description: &annonceDescription,
	UserID:      users[1].ID,
	CatID:       "1",
}
var association = models.Association{
	Name:       "Assoc de TEST",
	AddressRue: "10 rue du TEST",
	Cp:         "12345",
	Ville:      "TEST-CITY",
	Phone:      "0101010101",
	Email:      "emailDeTest@gmail.com",
	OwnerID:    users[1].ID,
	Verified:   &associationIsVerified,
}
var cat = models.Cats{
	Name:            "PABLO",
	BirthDate:       &dateCat,
	Sexe:            "male",
	LastVaccine:     &dateCat,
	LastVaccineName: "ANTI-PABLO",
	Color:           "VIOLET",
	Behavior:        "PABLOCITO",
	Sterilized:      false,
	RaceID:          strconv.Itoa(int(raceUnique.ID)),
	Description:     &descriptionCat,
	Reserved:        false,
	UserID:          users[1].ID,
}
var favorite = models.Favorite{
	UserID:    users[1].ID,
	AnnonceID: strconv.Itoa(int(annonces[1].ID)),
}
var race = models.Races{
	RaceName: "TESTOSAURUS",
}
var rating = models.Rating{
	Mark:     4,
	Comment:  "Rating de test",
	UserID:   users[1].ID,
	AuthorID: users[2].ID,
}
var room = models.Room{
	User1ID:   users[1].ID,
	User2ID:   users[2].ID,
	AnnonceID: strconv.Itoa(int(annonces[1].ID)),
}

func TestDatabaseService_CreateAnnonce(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var annonceTest test
	annonceTest.name = "Test creation annonce"
	annonceTest.fields = db
	annonceTest.wantId = 31
	annonceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: annonceTest.fields.s,
	}

	t.Run(annonceTest.name, func(t *testing.T) {

		gotId, err := s.CreateAnnonce(&annonce)
		if !annonceTest.wantErr(t, err, fmt.Sprintf("CreateAnnonce(%v)", annonce)) {
			return
		}
		assert.Equalf(t, annonceTest.wantId, gotId, "CreateAnnonce(%v)", annonce)
	})
}

func TestDatabaseService_CreateAssociation(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var associationTest test
	associationTest.name = "Test creation association"
	associationTest.fields = db
	associationTest.wantId = 1
	associationTest.wantErr = assert.NoError

	t.Run(associationTest.name, func(t *testing.T) {

		gotId, err := db.CreateAssociation(&association)
		if !associationTest.wantErr(t, err, fmt.Sprintf("CreateAssociation(%v)", association)) {
			return
		}
		assert.Equalf(t, associationTest.wantId, gotId, "CreateAssociation(%v)", association)
	})
}

func TestDatabaseService_CreateCat(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var catTest test
	catTest.name = "Test creation cat"
	catTest.fields = db
	catTest.wantId = 31
	catTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: catTest.fields.s,
	}

	t.Run(catTest.name, func(t *testing.T) {

		gotId, err := s.CreateCat(&cat)
		if !catTest.wantErr(t, err, fmt.Sprintf("CreateCat(%v)", cat)) {
			return
		}
		assert.Equalf(t, catTest.wantId, gotId, "CreateCat(%v)", cat)
	})
}

func TestDatabaseService_CreateFavorite(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var favoriteTest test
	favoriteTest.name = "Test creation favorite"
	favoriteTest.fields = db
	favoriteTest.wantId = 1
	favoriteTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: favoriteTest.fields.s,
	}

	t.Run(favoriteTest.name, func(t *testing.T) {

		gotId, err := s.CreateFavorite(&favorite)
		if !favoriteTest.wantErr(t, err, fmt.Sprintf("CreateFavorite(%v)", favorite)) {
			return
		}
		assert.Equalf(t, favoriteTest.wantId, gotId, "CreateFavorite(%v)", favorite)
	})
}

func TestDatabaseService_CreateRace(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var raceTest test
	raceTest.name = "Test creation race"
	raceTest.fields = db
	raceTest.wantId = 6
	raceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: raceTest.fields.s,
	}

	t.Run(raceTest.name, func(t *testing.T) {

		gotId, err := s.CreateRace(&race)
		if !raceTest.wantErr(t, err, fmt.Sprintf("CreateRace(%v)", race)) {
			return
		}
		assert.Equalf(t, raceTest.wantId, gotId, "CreateRace(%v)", race)
	})
}

func TestDatabaseService_CreateRating(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var ratingTest test
	ratingTest.name = "Test creation rating"
	ratingTest.fields = db
	ratingTest.wantId = 1
	ratingTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: ratingTest.fields.s,
	}

	t.Run(ratingTest.name, func(t *testing.T) {

		gotId, err := s.CreateRating(&rating)
		if !ratingTest.wantErr(t, err, fmt.Sprintf("CreateRating(%v)", rating)) {
			return
		}
		assert.Equalf(t, ratingTest.wantId, gotId, "CreateRating(%v)", rating)
	})
}

func TestDatabaseService_CreateRoom(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var roomTest test
	roomTest.name = "Test creation room"
	roomTest.fields = db
	roomTest.wantId = 1
	roomTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: roomTest.fields.s,
	}

	t.Run(roomTest.name, func(t *testing.T) {

		gotId, err := s.CreateRoom(&room)
		if !roomTest.wantErr(t, err, fmt.Sprintf("CreateRoom(%v)", room)) {
			return
		}
		assert.Equalf(t, roomTest.wantId, gotId, "CreateRoom(%v)", room)
	})
}

// FindBy
func TestDatabaseService_FindAnnonceByCatID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID uint
		wantId   uint
		wantErr  assert.ErrorAssertionFunc
	}

	var catId int = 23

	var annonceTest test
	annonceTest.name = "Test trouver annonce avec id chat"
	annonceTest.toFindID = 23
	annonceTest.fields = db
	annonceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: annonceTest.fields.s,
	}

	fmt.Println(annonce.ID)
	fmt.Println(annonce.CatID)

	t.Run(annonceTest.name, func(t *testing.T) {
		annonceFounded, err := s.FindAnnonceByCatID(strconv.Itoa(catId))
		fmt.Println(annonceFounded.ID)
		fmt.Println(annonceTest.toFindID)
		assert.NotNil(t, annonceFounded)
		assert.Equal(t, annonceTest.toFindID, annonceFounded.ID)
		//assert.IsType(t, models.Annonce{}, annonceFounded)
		annonceTest.wantErr(t, err, fmt.Sprintf("FindAnnonce(%v)", annonceTest.toFindID))
	})
}

func TestDatabaseService_FindAnnonceByID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID uint
		wantId   uint
		wantErr  assert.ErrorAssertionFunc
	}

	var annonceTest test
	annonceTest.name = "Test trouver annonce avec ID"
	annonceTest.toFindID = annonce.ID
	annonceTest.fields = db
	annonceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: annonceTest.fields.s,
	}

	t.Run(annonceTest.name, func(t *testing.T) {
		annonceFounded, err := s.FindAnnonceByID(strconv.Itoa(int(annonceTest.toFindID)))
		assert.NotNil(t, annonceFounded)
		assert.Equal(t, annonceTest.toFindID, annonceFounded.ID)
		//assert.IsType(t, models.Annonce{}, annonceFounded)
		annonceTest.wantErr(t, err, fmt.Sprintf("FindAnnonce(%v)", annonceTest.toFindID))
	})
}

func TestDatabaseService_FindAssociationById(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID uint
		wantId   uint
		wantErr  assert.ErrorAssertionFunc
	}

	var associationTest test
	associationTest.name = "Test trouver association par ID"
	associationTest.toFindID = association.ID
	associationTest.fields = db
	associationTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: associationTest.fields.s,
	}

	t.Run(associationTest.name, func(t *testing.T) {
		associationFounded, err := s.FindAssociationById(int(associationTest.toFindID))
		assert.NotNil(t, associationFounded)
		assert.Equal(t, associationTest.toFindID, associationFounded.ID)
		//assert.IsType(t, models.Association{}, associationFounded)
		associationTest.wantErr(t, err, fmt.Sprintf("FindAssociation(%v)", associationTest.toFindID))
	})
}

func TestDatabaseService_FindCatByID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var catTest test
	catTest.name = "Test trouver cat par ID"
	catTest.toFindID = strconv.Itoa(int(cat.ID))
	catTest.fields = db
	catTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: catTest.fields.s,
	}

	t.Run(catTest.name, func(t *testing.T) {
		catFounded, err := s.FindCatByID(catTest.toFindID)
		assert.NotNil(t, catFounded)
		assert.Equal(t, catTest.toFindID, strconv.Itoa(int(catFounded.ID)))
		//assert.IsType(t, models.Cat{}, catFounded)
		catTest.wantErr(t, err, fmt.Sprintf("FindCat(%v)", catTest.toFindID))
	})
}

func TestDatabaseService_FindCatsByUserID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var catTest test
	catTest.name = "Test trouver cat par userID"
	catTest.toFindID = strconv.Itoa(int(cat.ID))
	catTest.fields = db
	catTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: catTest.fields.s,
	}

	var idCatsFounded []string

	t.Run(catTest.name, func(t *testing.T) {
		catsFounded, err := s.FindCatsByUserID(cat.UserID)
		assert.NotNil(t, catsFounded)
		for _, as := range catsFounded {
			idCatsFounded = append(idCatsFounded, strconv.Itoa(int(as.ID)))
		}
		assert.Contains(t, idCatsFounded, catTest.toFindID)
		//assert.IsType(t, models.Cat{}, catFounded)
		catTest.wantErr(t, err, fmt.Sprintf("FindCat(%v)", catTest.toFindID))
	})
}

func TestDatabaseService_FindFavoriteByID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var favoriteTest test
	favoriteTest.name = "Test trouver favorite par ID"
	favoriteTest.toFindID = strconv.Itoa(int(favorite.ID))
	favoriteTest.fields = db
	favoriteTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: favoriteTest.fields.s,
	}

	t.Run(favoriteTest.name, func(t *testing.T) {
		favoriteFounded, err := s.FindFavoriteByID(favoriteTest.toFindID)
		assert.NotNil(t, favoriteFounded)
		assert.Equal(t, favoriteTest.toFindID, strconv.Itoa(int(favoriteFounded.ID)))
		//assert.IsType(t, models.Favorite{}, favoriteFounded)
		favoriteTest.wantErr(t, err, fmt.Sprintf("FindFavorite(%v)", favoriteTest.toFindID))
	})
}

func TestDatabaseService_FindFavoritesByAnnonceID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var favoriteTest test
	favoriteTest.name = "Test trouver favorite par userID"
	favoriteTest.toFindID = strconv.Itoa(int(favorite.ID))
	favoriteTest.fields = db
	favoriteTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: favoriteTest.fields.s,
	}

	var idFavoritesFounded []string

	t.Run(favoriteTest.name, func(t *testing.T) {
		favoritesFounded, err := s.FindFavoritesByAnnonceID(favorite.AnnonceID)
		assert.NotNil(t, favoritesFounded)
		for _, as := range favoritesFounded {
			idFavoritesFounded = append(idFavoritesFounded, strconv.Itoa(int(as.ID)))
		}
		assert.Contains(t, idFavoritesFounded, favoriteTest.toFindID)
		//assert.IsType(t, models.Favorite{}, favoriteFounded)
		favoriteTest.wantErr(t, err, fmt.Sprintf("FindFavorite(%v)", favoriteTest.toFindID))
	})
}

func TestDatabaseService_FindFavoritesByUserID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var favoriteTest test
	favoriteTest.name = "Test trouver favorite par userID"
	favoriteTest.toFindID = strconv.Itoa(int(favorite.ID))
	favoriteTest.fields = db
	favoriteTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: favoriteTest.fields.s,
	}

	var idFavoritesFounded []string

	t.Run(favoriteTest.name, func(t *testing.T) {
		favoritesFounded, err := s.FindFavoritesByUserID(favorite.UserID)
		assert.NotNil(t, favoritesFounded)
		for _, as := range favoritesFounded {
			idFavoritesFounded = append(idFavoritesFounded, strconv.Itoa(int(as.ID)))
		}
		assert.Contains(t, idFavoritesFounded, favoriteTest.toFindID)
		//assert.IsType(t, models.Favorite{}, favoriteFounded)
		favoriteTest.wantErr(t, err, fmt.Sprintf("FindFavorite(%v)", favoriteTest.toFindID))
	})
}

func TestDatabaseService_FindRaceByID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var raceTest test
	raceTest.name = "Test trouver race par ID"
	raceTest.toFindID = strconv.Itoa(int(race.ID))
	raceTest.fields = db
	raceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: raceTest.fields.s,
	}

	t.Run(raceTest.name, func(t *testing.T) {
		raceFounded, err := s.FindRaceByID(raceTest.toFindID)
		assert.NotNil(t, raceFounded)
		assert.Equal(t, raceTest.toFindID, strconv.Itoa(int(raceFounded.ID)))
		//assert.IsType(t, models.Race{}, raceFounded)
		raceTest.wantErr(t, err, fmt.Sprintf("FindRace(%v)", raceTest.toFindID))
	})
}

func TestDatabaseService_FindRatingByID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var ratingTest test
	ratingTest.name = "Test trouver rating par ID"
	ratingTest.toFindID = strconv.Itoa(int(rating.ID))
	ratingTest.fields = db
	ratingTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: ratingTest.fields.s,
	}

	t.Run(ratingTest.name, func(t *testing.T) {
		ratingFounded, err := s.FindRatingByID(ratingTest.toFindID)
		assert.NotNil(t, ratingFounded)
		assert.Equal(t, ratingTest.toFindID, strconv.Itoa(int(ratingFounded.ID)))
		//assert.IsType(t, models.Rating{}, ratingFounded)
		ratingTest.wantErr(t, err, fmt.Sprintf("FindRating(%v)", ratingTest.toFindID))
	})
}

func TestDatabaseService_FindRoomsByUserID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name     string
		fields   DatabaseService
		toFindID string
		wantErr  assert.ErrorAssertionFunc
	}

	var roomTest test
	roomTest.name = "Test trouver room par userID"
	roomTest.toFindID = strconv.Itoa(int(room.ID))
	roomTest.fields = db
	roomTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: roomTest.fields.s,
	}

	var idRoomsFounded []string

	t.Run(roomTest.name, func(t *testing.T) {
		roomsFounded, err := s.FindRoomsByUserID(room.User1ID)
		assert.NotNil(t, roomsFounded)
		for _, as := range roomsFounded {
			idRoomsFounded = append(idRoomsFounded, strconv.Itoa(int(as.ID)))
		}
		assert.Contains(t, idRoomsFounded, roomTest.toFindID)
		//assert.IsType(t, models.Room{}, roomFounded)
		roomTest.wantErr(t, err, fmt.Sprintf("FindRoom(%v)", roomTest.toFindID))
	})
}

func TestDatabaseService_GetAllAnnonces(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		toCount int
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var annonceTest test
	annonceTest.name = "Test trouver toutes les annonces"
	annonceTest.toCount = 31
	annonceTest.fields = db
	annonceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: annonceTest.fields.s,
	}

	t.Run(annonceTest.name, func(t *testing.T) {
		annonceFounded, err := s.GetAllAnnonces()
		assert.NotNil(t, annonceFounded)
		assert.Equal(t, annonceTest.toCount, len(annonceFounded))
		annonceTest.wantErr(t, err, fmt.Sprintf("FindAnnonce(%v)", annonceTest.toCount))
	})
}

func TestDatabaseService_GetAllAssociations(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		toCount int
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var associationTest test
	associationTest.name = "Test trouver toutes les associations"
	associationTest.toCount = 1
	associationTest.fields = db
	associationTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: associationTest.fields.s,
	}

	t.Run(associationTest.name, func(t *testing.T) {
		associationFounded, err := s.GetAllAssociations()
		assert.NotNil(t, associationFounded)
		assert.Equal(t, associationTest.toCount, len(associationFounded))
		associationTest.wantErr(t, err, fmt.Sprintf("FindAssociation(%v)", associationTest.toCount))
	})
}

func TestDatabaseService_GetAllCats(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		toCount int
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var catTest test
	catTest.name = "Test trouver toutes les cats"
	catTest.toCount = 31
	catTest.fields = db
	catTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: catTest.fields.s,
	}

	t.Run(catTest.name, func(t *testing.T) {
		catFounded, err := s.GetAllCats()
		assert.NotNil(t, catFounded)
		assert.Equal(t, catTest.toCount, len(catFounded))
		catTest.wantErr(t, err, fmt.Sprintf("FindCat(%v)", catTest.toCount))
	})
}

func TestDatabaseService_GetAllRace(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		toCount int
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var raceTest test
	raceTest.name = "Test trouver toutes les races"
	raceTest.toCount = 6
	raceTest.fields = db
	raceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: raceTest.fields.s,
	}

	t.Run(raceTest.name, func(t *testing.T) {
		raceFounded, err := s.GetAllRace()
		assert.NotNil(t, raceFounded)
		assert.Equal(t, raceTest.toCount, len(raceFounded))
		raceTest.wantErr(t, err, fmt.Sprintf("FindRace(%v)", raceTest.toCount))
	})
}

func TestDatabaseService_GetAllRatings(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		toCount int
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var ratingTest test
	ratingTest.name = "Test trouver toutes les ratings"
	ratingTest.toCount = 1
	ratingTest.fields = db
	ratingTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: ratingTest.fields.s,
	}

	t.Run(ratingTest.name, func(t *testing.T) {
		ratingFounded, err := s.GetAllRatings()
		assert.NotNil(t, ratingFounded)
		assert.Equal(t, ratingTest.toCount, len(ratingFounded))
		ratingTest.wantErr(t, err, fmt.Sprintf("FindRating(%v)", ratingTest.toCount))
	})
}

func TestDatabaseService_GetAllUsers(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		toCount int
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var userTest test
	userTest.name = "Test trouver toutes les users"
	userTest.toCount = 6
	userTest.fields = db
	userTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: userTest.fields.s,
	}

	t.Run(userTest.name, func(t *testing.T) {
		userFounded, err := s.GetAllUsers()
		assert.NotNil(t, userFounded)
		assert.Equal(t, userTest.toCount, len(userFounded))
		userTest.wantErr(t, err, fmt.Sprintf("FindUser(%v)", userTest.toCount))
	})
}

func TestDatabaseService_DeleteAssociation(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name       string
		fields     DatabaseService
		toDeleteID uint
		wantId     uint
		wantErr    assert.ErrorAssertionFunc
	}

	var associationTest test
	associationTest.name = "Test suppression association"
	associationTest.toDeleteID = 1
	associationTest.fields = db
	associationTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: associationTest.fields.s,
	}

	t.Run(associationTest.name, func(t *testing.T) {
		associationTest.wantErr(t, s.DeleteAssociation(int(associationTest.toDeleteID)), fmt.Sprintf("DeleteAssociation(%v)", associationTest.toDeleteID))
		association, _ := s.FindAssociationById(int(associationTest.toDeleteID))
		assert.Nil(t, association)
	})
}

func TestDatabaseService_DeleteCatByID(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name       string
		fields     DatabaseService
		toDeleteID uint
		wantId     uint
		wantErr    assert.ErrorAssertionFunc
	}

	var catTest test
	catTest.name = "Test suppression cat"
	catTest.toDeleteID = 31
	catTest.fields = db
	catTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: catTest.fields.s,
	}

	t.Run(catTest.name, func(t *testing.T) {
		catTest.wantErr(t, s.DeleteCatByID(strconv.Itoa(int(catTest.toDeleteID))), fmt.Sprintf("DeleteCat(%v)", catTest.toDeleteID))
		cat, _ := s.FindCatByID(strconv.Itoa(int(catTest.toDeleteID)))
		assert.Nil(t, cat)
	})
}

func TestDatabaseService_DeleteFavorite(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name       string
		fields     DatabaseService
		toDeleteID uint
		wantId     uint
		wantErr    assert.ErrorAssertionFunc
	}

	var favoriteTest test
	favoriteTest.name = "Test suppression favorite"
	favoriteTest.toDeleteID = 1
	favoriteTest.fields = db
	favoriteTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: favoriteTest.fields.s,
	}

	favorite, _ := s.FindFavoriteByID(strconv.Itoa(int(favoriteTest.toDeleteID)))

	t.Run(favoriteTest.name, func(t *testing.T) {
		favoriteTest.wantErr(t, s.DeleteFavorite(favorite), fmt.Sprintf("DeleteFavorite(%v)", favoriteTest.toDeleteID))
		deletedFavorite, _ := s.FindFavoriteByID(strconv.Itoa(int(favoriteTest.toDeleteID)))
		assert.Nil(t, deletedFavorite)
	})
}

func TestDatabaseService_DeleteRace(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name       string
		fields     DatabaseService
		toDeleteID uint
		wantId     uint
		wantErr    assert.ErrorAssertionFunc
	}

	var raceTest test
	raceTest.name = "Test suppression race"
	raceTest.toDeleteID = 6
	raceTest.fields = db
	raceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: raceTest.fields.s,
	}

	t.Run(raceTest.name, func(t *testing.T) {
		raceTest.wantErr(t, s.DeleteRace(strconv.Itoa(int(raceTest.toDeleteID))), fmt.Sprintf("DeleteRace(%v)", raceTest.toDeleteID))
		race, _ := s.FindRaceByID(strconv.Itoa(int(raceTest.toDeleteID)))
		assert.Nil(t, race)
	})
}

func TestDatabaseService_DeleteRating(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name       string
		fields     DatabaseService
		toDeleteID uint
		wantId     uint
		wantErr    assert.ErrorAssertionFunc
	}

	var ratingTest test
	ratingTest.name = "Test suppression rating"
	ratingTest.toDeleteID = 31
	ratingTest.fields = db
	ratingTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: ratingTest.fields.s,
	}

	t.Run(ratingTest.name, func(t *testing.T) {
		ratingTest.wantErr(t, s.DeleteRating(strconv.Itoa(int(ratingTest.toDeleteID))), fmt.Sprintf("DeleteRating(%v)", ratingTest.toDeleteID))
		rating, _ := s.FindRatingByID(strconv.Itoa(int(ratingTest.toDeleteID)))
		assert.Nil(t, rating)
	})
}
