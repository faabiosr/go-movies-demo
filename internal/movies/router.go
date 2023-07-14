package movies

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type router struct {
	ds *Datasource
}

// Routes creates and expose the router.
func Routes(root *echo.Group, ds *Datasource) {
	r := &router{ds: ds}

	root.GET("", r.all)
	root.POST("", r.create)
	root.GET("/:movie-id", r.retrieve)
	root.PUT("/:movie-id", r.update)
	root.DELETE("/:movie-id", r.remove)
}

func (r *router) all(ec echo.Context) error {
	movies, err := r.ds.FetchAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, movies)
}

func (r *router) create(ec echo.Context) error {
	m := Movie{}
	if err := ec.Bind(&m); err != nil {
		return err
	}

	if err := m.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := r.ds.Store(m)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusCreated, m)
}

func (r *router) retrieve(ec echo.Context) error {
	id := ec.Param("movie-id")
	m, err := r.ds.Fetch(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return ec.JSON(http.StatusOK, m)
}

func (r *router) update(ec echo.Context) error {
	id := ec.Param("movie-id")

	m := Movie{ID: id}
	if err := ec.Bind(&m); err != nil {
		return err
	}

	if err := m.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := r.ds.Fetch(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if _, err := r.ds.Store(m); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.NoContent(http.StatusNoContent)
}

func (r *router) remove(ec echo.Context) error {
	id := ec.Param("movie-id")

	if err := r.ds.Remove(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.NoContent(http.StatusNoContent)
}
