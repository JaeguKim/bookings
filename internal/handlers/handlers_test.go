package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tsawler/bookings-app/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq","/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET",http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	//{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	//{"post-search-avail", "/search-availability","POST",[]postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
	//{"post-search-avail-json", "/search-availability-json","POST",[]postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
	//{"make-reservation-post", "/make-reservation","POST",[]postData{
	//	{key: "first_name", value: "John"},
	//	{key: "last_name", value: "Smith"},
	//	{key: "email", value: "me@here.com"},
	//	{key: "phone", value: "555-555-5555"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d",e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room {
			ID: 1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Rerservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Rerservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Rerservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room {
			ID: 1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("RerservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/reservation-summary",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("RerservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room {
			ID: 1,
			RoomName: "General's Quarters",
		},
	}

	// test case where id is not in url param
	req, _ := http.NewRequest("GET", "/choose-room/fish", nil)
	req.RequestURI = "/choose-room/fish"
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/choose-room/1",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case for success
	req, _ = http.NewRequest("GET", "/choose-room/1",nil)
	req.RequestURI = "/choose-room/1"
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	// test case where room is not found in database
	req, _ := http.NewRequest("GET", "/book-room?id=123&s=2050-01-01&e=2050-01-03", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case for success
	req, _ = http.NewRequest("GET", "/book-room?id=1&s=2050-01-01&e=2050-01-03", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation",strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostRerservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostRerservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostRerservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostRerservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostRerservation handler returned wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostRerservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert reservation into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostRerservation handler failed when trying to insert reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure to insert room restriction into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody, "room_id=1000")

	req, _ = http.NewRequest("POST", "/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostRerservation handler failed when trying to insert room restriction: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// test for missing post body
	req, _ := http.NewRequest("POST", "/post-availability",nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for success
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for error when parsing start in form data
	reqBody = "start=invalid-data"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for error when parsing end in form data
	reqBody = "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end=invalid-data")
	req, _ = http.NewRequest("POST", "/search-availability",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for error when searching availability for all rooms from database
	reqBody = "start=2050-11-11"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end=2050-11-12")
	req, _ = http.NewRequest("POST", "/search-availability",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for case when rooms slice is empty
	reqBody = "start=2050-10-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end=2050-10-02")
	req, _ = http.NewRequest("POST", "/search-availability",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// test for missing post body
	req, _ := http.NewRequest("POST", "/search-availability-json",nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("AvailabilityJSON handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for success
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability-json",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusOK {
		t.Errorf("AvailabilityJSON handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}


func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}