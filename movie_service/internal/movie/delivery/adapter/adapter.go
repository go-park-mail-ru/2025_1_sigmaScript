package adapter

import (
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	pb "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/api/movie_api_v1/proto"
)

func SrvPersonToDescPerson(srv *mocks.PersonJSON) *pb.PersonJSON {
	if srv == nil {
		return nil
	}
	mc := make(map[int32]*pb.MovieJSON, len(srv.MovieCollection))
	for k, mv := range srv.MovieCollection {
		mc[int32(k)] = &pb.MovieJSON{
			Id:              int32(mv.ID),
			Name:            mv.Title,
			OriginalName:    mv.Title,
			About:           "",
			Poster:          mv.PreviewURL,
			PromoUrl:        "",
			ReleaseYear:     mv.ReleaseDate,
			Country:         "",
			Slogan:          "",
			Director:        "",
			Budget:          0,
			BoxOfficeUs:     0,
			BoxOfficeGlobal: 0,
			BoxOfficeRussia: 0,
			PremierRussia:   "",
			PremierGlobal:   "",
			Rating:          mv.Rating,
			Duration:        mv.Duration,
			Genres:          "",
			Staff:           nil,
			Reviews:         nil,
			RatingKp:        0,
			RatingImdb:      0,
			Watchability:    nil,
			Logo:            "",
			Backdrop:        "",
		}
	}
	return &pb.PersonJSON{
		Id:              int32(srv.ID),
		FullName:        srv.FullName,
		EnFullName:      srv.EnFullName,
		Photo:           srv.Photo,
		About:           srv.About,
		Sex:             srv.Sex,
		Growth:          srv.Growth,
		Birthday:        srv.Birthday,
		Death:           srv.Death,
		Career:          srv.Career,
		Genres:          srv.Genres,
		TotalFilms:      srv.TotalFilms,
		MovieCollection: mc,
	}
}

func SrvPersonsToDescPersons(srvs []mocks.PersonJSON) []*pb.PersonJSON {
	out := make([]*pb.PersonJSON, 0, len(srvs))
	for i := range srvs {
		out = append(out, SrvPersonToDescPerson(&srvs[i]))
	}
	return out
}

func SrvReviewToDescReview(srv *mocks.ReviewJSON) *pb.ReviewJSON {
	if srv == nil {
		return nil
	}
	return &pb.ReviewJSON{
		Id:         int32(srv.ID),
		User:       &pb.ReviewUserDataJSON{Login: srv.User.Login, Avatar: srv.User.Avatar},
		ReviewText: srv.ReviewText,
		Score:      srv.Score,
		CreatedAt:  srv.CreatedAt,
	}
}

func SrvReviewsToDescReviews(srvs []mocks.ReviewJSON) []*pb.ReviewJSON {
	out := make([]*pb.ReviewJSON, 0, len(srvs))
	for i := range srvs {
		out = append(out, SrvReviewToDescReview(&srvs[i]))
	}
	return out
}

func SrvWatchProvidersToDesc(srv []mocks.WatchProviderJSON) []*pb.WatchProviderJSON {
	out := make([]*pb.WatchProviderJSON, 0, len(srv))
	for _, w := range srv {
		out = append(out, &pb.WatchProviderJSON{Name: w.Name, Url: w.Url})
	}
	return out
}

func SrvMovieToDescMovie(srv *mocks.MovieJSON) *pb.MovieJSON {
	if srv == nil {
		return nil
	}
	return &pb.MovieJSON{
		Id:              int32(srv.ID),
		Name:            srv.Name,
		OriginalName:    srv.OriginalName,
		About:           srv.About,
		Poster:          srv.Poster,
		PromoUrl:        srv.PromoURL,
		ReleaseYear:     srv.ReleaseYear,
		Country:         srv.Country,
		Slogan:          srv.Slogan,
		Director:        srv.Director,
		Budget:          srv.Budget,
		BoxOfficeUs:     srv.BoxOfficeUS,
		BoxOfficeGlobal: srv.BoxOfficeGlobal,
		BoxOfficeRussia: srv.BoxOfficeRussia,
		PremierRussia:   srv.PremierRussia,
		PremierGlobal:   srv.PremierGlobal,
		Rating:          srv.Rating,
		Duration:        srv.Duration,
		Genres:          srv.Genres,
		Staff:           SrvPersonsToDescPersons(srv.Staff),
		Reviews:         SrvReviewsToDescReviews(srv.Reviews),
		RatingKp:        srv.RatingKP,
		RatingImdb:      srv.RatingIMDB,
		Watchability:    SrvWatchProvidersToDesc(srv.Watchability),
		Logo:            srv.Logo,
		Backdrop:        srv.Backdrop,
	}
}

func DescNewReviewToSrvNewReview(desc *pb.NewReviewDataJSON) mocks.NewReviewDataJSON {
	if desc == nil {
		return mocks.NewReviewDataJSON{}
	}
	return mocks.NewReviewDataJSON{ReviewText: desc.ReviewText, Score: desc.Score}
}

func ToSrvGetMovieByID(req *pb.GetMovieByIDRequest) int {
	if req == nil {
		return 0
	}
	return int(req.MovieId)
}

func ToSrvGetAllReviewsOfMovieByID(req *pb.GetAllReviewsOfMovieByIDRequest) int {
	if req == nil {
		return 0
	}
	return int(req.MovieId)
}

func ToSrvNewReviewDataFromDesc(req *pb.CreateNewMovieReviewRequest) mocks.NewReviewDataJSON {
	if req == nil || req.NewReview == nil {
		return mocks.NewReviewDataJSON{}
	}
	return DescNewReviewToSrvNewReview(req.NewReview)
}

func ToSrvUpdateReviewDataFromDesc(req *pb.UpdateMovieReviewRequest) mocks.NewReviewDataJSON {
	if req == nil || req.NewReview == nil {
		return mocks.NewReviewDataJSON{}
	}
	return DescNewReviewToSrvNewReview(req.NewReview)
}

func ToSrvPersonID(req *pb.GetPersonByIDRequest) int {
	if req == nil {
		return 0
	}
	return int(req.PersonId)
}

func ToSrvGenreID(req *pb.GetGenreByIDRequest) string {
	if req == nil {
		return ""
	}
	return req.GenreId
}

func ToSrvSearchQuery(req *pb.SearchActorsAndMoviesRequest) string {
	if req == nil {
		return ""
	}
	return req.Query
}

func ToDescGetMovieByIDResponse(srv *mocks.MovieJSON) *pb.GetMovieByIDResponse {
	return &pb.GetMovieByIDResponse{Movie: SrvMovieToDescMovie(srv)}
}

func ToDescGetAllReviewsOfMovieByIDResponse(sr *[]mocks.ReviewJSON) *pb.GetAllReviewsOfMovieByIDResponse {
	if sr == nil {
		return &pb.GetAllReviewsOfMovieByIDResponse{}
	}
	return &pb.GetAllReviewsOfMovieByIDResponse{Reviews: SrvReviewsToDescReviews(*sr)}
}

func ToDescCreateNewMovieReviewResponse(srv *mocks.NewReviewDataJSON) *pb.CreateNewMovieReviewResponse {
	return &pb.CreateNewMovieReviewResponse{Review: &pb.NewReviewDataJSON{ReviewText: srv.ReviewText, Score: srv.Score}}
}

func ToDescUpdateMovieReviewResponse(srv *mocks.NewReviewDataJSON) *pb.UpdateMovieReviewResponse {
	return &pb.UpdateMovieReviewResponse{Review: &pb.NewReviewDataJSON{ReviewText: srv.ReviewText, Score: srv.Score}}
}

func ToDescGetPersonByIDResponse(p *mocks.PersonJSON) *pb.GetPersonByIDResponse {
	return &pb.GetPersonByIDResponse{Person: SrvPersonToDescPerson(p)}
}

func movieToMovieJSON(m *mocks.Movie) *mocks.MovieJSON {
	if m == nil {
		return nil
	}
	return &mocks.MovieJSON{
		ID:          m.ID,
		Name:        m.Title,
		Poster:      m.PreviewURL,
		Duration:    m.Duration,
		ReleaseYear: m.ReleaseDate,
		Rating:      m.Rating,
	}
}

func ToDescGetGenreByIDResponse(g *mocks.Genre) *pb.GetGenreByIDResponse {
	if g == nil {
		return &pb.GetGenreByIDResponse{}
	}
	movies := make([]*pb.MovieJSON, 0, len(g.Movies))
	for i := range g.Movies {
		movies = append(movies, SrvMovieToDescMovie(movieToMovieJSON(&g.Movies[i])))
	}
	return &pb.GetGenreByIDResponse{Genre: &pb.Genre{Id: g.ID, Name: g.Name, Movies: movies}}
}

func ToDescGetAllGenresResponse(list *[]mocks.Genre) *pb.GetAllGenresResponse {
	out := &pb.GetAllGenresResponse{}
	if list == nil {
		return out
	}
	for _, g := range *list {
		movies := make([]*pb.MovieJSON, 0, len(g.Movies))
		for i := range g.Movies {
			movies = append(movies, SrvMovieToDescMovie(movieToMovieJSON(&g.Movies[i])))
		}
		out.Genres = append(out.Genres, &pb.Genre{Id: g.ID, Name: g.Name, Movies: movies})
	}
	return out
}

func ToDescSearchActorsAndMoviesResponse(sr *models.SearchResponseJSON) *pb.SearchActorsAndMoviesResponse {
	out := &pb.SearchActorsAndMoviesResponse{}
	if sr == nil {
		return out
	}
	for _, m := range sr.MovieCollection {
		out.Movies = append(out.Movies, &pb.MovieJSON{
			Id:              int32(m.ID),
			Name:            m.Title,
			OriginalName:    m.Title,
			About:           "",
			Poster:          m.PreviewURL,
			PromoUrl:        "",
			ReleaseYear:     m.ReleaseDate,
			Country:         "",
			Slogan:          "",
			Director:        "",
			Budget:          0,
			BoxOfficeUs:     0,
			BoxOfficeGlobal: 0,
			BoxOfficeRussia: 0,
			PremierRussia:   "",
			PremierGlobal:   "",
			Rating:          m.Rating,
			Duration:        m.Duration,
			Genres:          "",
			Staff:           nil,
			Reviews:         nil,
			RatingKp:        0,
			RatingImdb:      0,
			Watchability:    nil,
			Logo:            "",
			Backdrop:        "",
		})
	}
	for _, a := range sr.Actors {
		out.Actors = append(out.Actors, &pb.PersonJSON{
			Id:              int32(a.ID),
			FullName:        a.FullName,
			EnFullName:      a.EnFullName,
			Photo:           a.Photo,
			About:           a.About,
			Sex:             a.Sex,
			Growth:          a.Growth,
			Birthday:        a.Birthday,
			Death:           a.Death,
			Career:          a.Career,
			Genres:          a.Genres,
			TotalFilms:      a.TotalFilms,
			MovieCollection: nil,
		})
	}
	return out
}

func ToDescGetMainPageCollectionsResponse(cols mocks.Collections) *pb.GetMainPageCollectionsResponse {
	out := &pb.GetMainPageCollectionsResponse{}
	for name, coll := range cols {
		movies := make([]*pb.MovieJSON, 0, len(coll))
		for _, m := range coll {
			movies = append(movies, SrvSimpleMovieToDescMovie(&m))
		}
		out.Collections = append(out.Collections, &pb.Collection{
			Name:   name,
			Movies: movies,
		})
	}
	return out
}

func SrvSimpleMovieToDescMovie(srv *mocks.Movie) *pb.MovieJSON {
	if srv == nil {
		return nil
	}
	return &pb.MovieJSON{
		Id:              int32(srv.ID),
		Name:            srv.Title,
		OriginalName:    srv.Title,
		About:           "",
		Poster:          srv.PreviewURL,
		PromoUrl:        "", // …
		ReleaseYear:     srv.ReleaseDate,
		Country:         "",
		Slogan:          "",
		Director:        "",
		Budget:          0,
		BoxOfficeUs:     0,
		BoxOfficeGlobal: 0,
		BoxOfficeRussia: 0,
		PremierRussia:   "",
		PremierGlobal:   "",
		Rating:          srv.Rating,
		Duration:        srv.Duration,
		Genres:          "",
		Staff:           nil,
		Reviews:         nil,
		RatingKp:        0,
		RatingImdb:      0,
		Watchability:    nil,
		Logo:            "",
		Backdrop:        "",
	}
}
