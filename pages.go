package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ThatAdwaithGuy/req/db/query"
	"github.com/ThatAdwaithGuy/req/views"
	"github.com/gin-gonic/gin"
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

func (q *Query) PostDesigns(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	formName := ctx.PostForm("form_name")
	controlLevel := ctx.PostForm("control_level_name")
	var controlType query.DataTypes = ""
	switch ctx.PostForm("control_type") {
	case "text":
		controlType = query.DataTypesText     
	case "number":
		controlType = query.DataTypesNumber   
	case "checkbox":
		controlType = query.DataTypesCheckbox 
	case "textarea":
		controlType = query.DataTypesTextarea 
	case "select":
		controlType = query.DataTypesSelect   
	}

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
		Column1:     formName,
		LabelName:   controlLevel,
		Column3:    controlType,
		IsMandatory: isMandatory,
		Sequence:    int32(sequence),
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

func (q *Query) GetFormEntries(ctx *gin.Context) {
	formName := ctx.Query("form")

	designs, err := q.query.GetDesignByFormName(ctx.Request.Context(), formName)
	if err != nil {
		err_msg := fmt.Sprintf("error while getting design: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return
	}

	fmt.Printf("Entries of %s: %s\n", formName, designs)

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.FormEntriesTable(designs).Render(ctx, ctx.Writer); err != nil {
		fmt.Println("Error: ", err)
	}
}

func (q *Query) ViewDesigns(ctx *gin.Context) {
	designs, err := q.query.GetAllFormNames(ctx.Request.Context())
	if err != nil {
		err_msg := fmt.Sprintf("error while getting design: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err := views.FormsListPage(designs).Render(ctx, ctx.Writer); err != nil {
		fmt.Println("Error: ", err)
	}
}
