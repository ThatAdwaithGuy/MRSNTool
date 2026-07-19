package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ThatAdwaithGuy/req/db/query"
	"github.com/ThatAdwaithGuy/req/views"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Index(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.Index().Render(ctx, ctx.Writer); err != nil {
		fmt.Println("Error: ", err)
	}
}

func CreateDesignForm(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.CreateDesignPage().Render(ctx, ctx.Writer); err != nil {
		fmt.Println("Error: ", err)
	}
}

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

func (q *Query) PostDesigns(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	formName := ctx.PostForm("form_name")
	controlLevel := ctx.PostForm("control_level_name")
	controlType := ctx.PostForm("control_type")
	sequence, err := strconv.Atoi(ctx.PostForm("sequence"))
	if err != nil {
		err_msg := fmt.Sprintf("Sequence is not a int: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return

	}
	isMandatory := ctx.PostForm("is_mandatory") == "true"
	params := query.NewDesignParams{
		FormName:         formName,
		ControlLevelName: controlLevel,
		ControlType:      controlType,
		IsMandatory:      isMandatory,
		Sequence:         int32(sequence),
	}
	design, err := q.query.NewDesign(ctx.Request.Context(), params)
	if err != nil {
		err_msg := fmt.Sprintf("error while inserting design: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return
	}
	fmt.Println("Inserted design: ", design)
	str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(true, "Successfully inserted design"))
	if err != nil {
		fmt.Println("error while rendering success message for form design insert")
		return
	}

	ctx.String(http.StatusOK, str)
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
	r.POST("/designs", query.PostDesigns)

	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
