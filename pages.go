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
	formIDstr := ctx.Param("id")
	formID, err := strconv.Atoi(formIDstr)
	if err != nil {
		err_msg := fmt.Sprintf("formID is not a number: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	rows, err := q.query.GetDesignByFormID(ctx.Request.Context(), int32(formID))
	if err != nil {
		err_msg := fmt.Sprintf("Error while accessing rows: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}
	form_name, err := q.query.GetFormNameFromID(ctx.Request.Context(), int32(formID))
	if err != nil {
		err_msg := fmt.Sprintf("Error while accessing form_name: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	fields := []views.Field{}
	for _, row := range rows {
		fields = append(fields, views.Field{
			LabelName: row.LabelName,
			DataType:  row.DataType,
			Required:  row.IsMandatory,
		})
	}
	fmt.Println(fields)
	fmt.Printf("Enter page forms info:\nfields: %+v,\nform_name: %s,\nform_id: %d\nform_id_str: %s", fields, form_name, formID, formIDstr)

	if err := views.EnterFormPage(fields, form_name, int32(formID)).Render(ctx, ctx.Writer); err != nil {
		err_msg := fmt.Sprintf("Error while rendering enter form page: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (q *Query) enterForm(ctx *gin.Context, arg query.NewDataEntryParams) *query.Datum {
	ret, err := q.query.NewDataEntry(ctx.Request.Context(), arg)
	if err != nil {
		err_msg := fmt.Sprintf("Error while inserting data form entry: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return nil
		}
		ctx.Status(http.StatusInternalServerError)
		return nil
	}
	return &ret
}

func (q *Query) NewFormEntry(ctx *gin.Context) {
	formIDstr := ctx.Param("id")
	formID, err := strconv.Atoi(formIDstr)
	if err != nil {
		err_msg := fmt.Sprintf("Form ID is not a number: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	schema, err := q.query.GetDesignByFormID(ctx.Request.Context(), int32(formID))
	if err != nil {
		err_msg := fmt.Sprintf("Error querying designs for form entry: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	form_entry_id, err := q.query.GetNextFormEntryId(ctx.Request.Context(), int32(formID))
	if err != nil {
		err_msg := fmt.Sprintf("Error generating form entry id: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	for k, v := range ctx.Request.PostForm {
		val := v[0]

		var design query.GetDesignByFormIDRow
		for i, des := range schema {
			if des.LabelName == k {
				design = des
				schema[i] = schema[len(schema)-1]
				schema = schema[:len(schema)-1]
				break
			}
		}

		fmt.Println("val: ", val)

		data_entry := q.enterForm(ctx, query.NewDataEntryParams{
			Data:        val,
			DataType:    design.DataType,
			FormID:      int32(formID),
			FormEntryID: form_entry_id,
		})
		if data_entry != nil {
			fmt.Printf("Data entry: %+v\n", *data_entry)
		} else {
			fmt.Println("Data entry is nil", data_entry)
		}
	}

	if len(schema) != 0 {
		for _, des := range schema {
		data_entry := q.enterForm(ctx, query.NewDataEntryParams{
			Data:        "",
			DataType:    des.DataType,
			FormID:      int32(formID),
			FormEntryID: form_entry_id,
		})
		if data_entry != nil {
			fmt.Printf("Filled empty Data: %+v\n", *data_entry)
		} else {
			fmt.Println("Data entry is nil", data_entry)
		}

		}
	}

	// TODO: render and write a grey out button for successful state
	ctx.Status(http.StatusOK)
}

func (q *Query) GetSelectOptions(ctx *gin.Context) {
	field := ctx.Query("field")
	form_id_str := ctx.Query("form_id")
	form_id, err := strconv.Atoi(form_id_str)
	if err != nil {
		err_msg := fmt.Sprintf("form_id ins't a number: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	params := query.GetOptionsDropDownBoxParams{
		LabelName: field,
		FormID:    int32(form_id),
	}

	options, err := q.query.GetOptionsDropDownBox(ctx.Request.Context(), params)
	if err != nil {
		err_msg := fmt.Sprintf("Error while querying options for dropdown: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err := views.SelectOptions(options).Render(ctx, ctx.Writer); err != nil {
		err_msg := fmt.Sprintf("Error while rendering select options page: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (q *Query) GetFormStats(ctx *gin.Context) {
	form_name := ctx.Query("form")
	fmt.Println("form name:", form_name)

	fieldCount, err := q.query.GetRowCountFromFormName(ctx.Request.Context(), form_name)
	if err != nil {
		fmt.Println("ERROR while access row count: ", err)
		fieldCount = 0
	}

	rowCount, err := q.query.GetDataRowCountFromFormName(ctx.Request.Context(), form_name)
	if err != nil {
		fmt.Println("ERROR while access data row count: ", err)
		rowCount = 0
	}

	if err := views.FormStatsPartial(fieldCount, rowCount / fieldCount ).Render(ctx, ctx.Writer); err != nil {
		err_msg := fmt.Sprintf("Error while rendering form stats: %s", err)
		if err := views.AlertMessage(false, err_msg).Render(ctx, ctx.Writer); err != nil {
			fmt.Printf("error while rendering error message for %s\n", err_msg)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return

	}
}
