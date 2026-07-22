package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ThatAdwaithGuy/req/db/query"
	"github.com/ThatAdwaithGuy/req/views"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

	dropdown_id := pgtype.Int4{
		Int32: 0,
		Valid: false,
	}

	// DatatypeSelect is dropdown box
	if controlType == query.DataTypesSelect {
		options_string := ctx.PostForm("dropdown_options")
		options := strings.Split(options_string, ",")
		dropdown, err := q.query.NewDropDown(ctx.Request.Context(), options)
		if err != nil {
			err_msg := fmt.Sprintf("Error while inserting new dropdown box values: %s", err)
			str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
			if err != nil {
				fmt.Printf("error while rendering error message for %s\n", err_msg)
				return
			}
			ctx.String(http.StatusInternalServerError, str)
			return
		}
		dropdown_id.Int32 = dropdown.ID
		dropdown_id.Valid = true
	}

	params := query.NewDesignParams{
		Column1:     formName,
		LabelName:   controlLevel,
		Column3:     controlType,
		IsMandatory: isMandatory,
		Sequence:    int32(sequence),
		DropdownID:  dropdown_id,
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

	fmt.Printf("Entries of %s: %+v\n", formName, designs)

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

type JsonForms struct {
	ID        int32  `json:"id"`
	FormName  string `json:"form_name"`
	Enterable bool   `json:"enterable"`
}

func (q *Query) GetEnterableForms(ctx *gin.Context) {
	designs, err := q.query.GetAllEnterableForms(ctx.Request.Context())
	if err != nil {
		err_msg := fmt.Sprintf("error while querying enterable forms: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return
	}

	var json_forms []JsonForms
	for _, form := range designs {
		json_forms = append(json_forms, JsonForms{
			ID:       form.ID,
			FormName: form.FormName,
			// TODO:  Bool could be invalid so add a check for it.
			Enterable: form.Enterable.Bool,
		})
	}

	fmt.Println("json forms: ", json_forms)

	ctx.JSON(http.StatusOK, json_forms)
}

func (q *Query) MakeFormEnterable(ctx *gin.Context) {
	name := ctx.PostForm("name")
	fmt.Printf("Made %s enterable\n", name)
	err := q.query.SetFormEnterable(ctx.Request.Context(), name)
	if err != nil {
		err_msg := fmt.Sprintf("error while getting enterable forms: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return
	}

	if err := views.FormButton("", true).Render(ctx, ctx.Writer); err != nil {
		err_msg := fmt.Sprintf("error while rendering button: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return
	}

	ctx.Status(http.StatusOK)
}

func (q *Query) FormsPage(ctx *gin.Context) {
	forms, err := q.query.GetAllEnterableForms(ctx.Request.Context())
	if err != nil {
		err_msg := fmt.Sprintf("error while qurrying enterable forms: %s", err)
		str, err := RenderToString(ctx.Request.Context(), views.AlertMessage(false, err_msg))
		if err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.String(http.StatusInternalServerError, str)
		return
	}
	fmt.Println(forms)

	if err := views.FormsGridPage(forms).Render(ctx, ctx.Writer); err != nil {
		fmt.Println(err)
		err_msg := fmt.Sprintf("error while rendering forms grid page: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}
}


func (q *Query) EnterForm(ctx *gin.Context) {
}
