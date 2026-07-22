package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ThatAdwaithGuy/req/db/query"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Query struct {
	query *query.Queries
}

func RenderToString(ctx context.Context, component templ.Component) (string, error) {
	var buf bytes.Buffer

	if err := component.Render(ctx, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
func main() {
	r := gin.Default()
	ctx := context.Background()

	dbURL := os.Getenv("DB_CONTAINER_URL")
	if dbURL == "" {
		fmt.Println("Can't get  DB_CONTAINER_URL")
		return
	}
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Println("error while opening connection", err)
		return
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		fmt.Println("Database ping failed:", err)
		return
	}

	queries := query.New(conn)
	query := Query{query: queries}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", Index)
	r.GET("/create_form", CreateDesignForm)
	r.GET("/view_designs", query.ViewDesigns)
	r.GET("/forms", query.FormsPage)
	r.GET("/enter_form/:id", query.EnterForm)

	r.POST("/api/post_designs", query.PostDesigns)
	r.GET("/api/get_form_designs", query.GetFormEntries)
	r.GET("/api/get_enterable_forms", query.GetEnterableForms)
	r.POST("/api/make_form_enterable", query.MakeFormEnterable)
	

	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
