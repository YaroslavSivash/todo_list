package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	"todo_list/models"

	beego "github.com/beego/beego/v2/server/web"
)

// NotesController operations for Notes
type NotesController struct {
	beego.Controller
}

// URLMapping ...
func (c *NotesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

type NotesData struct {
	Id        int       `json:"id"`
	Body      string    `json:"body"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Post ...
// @Title Post
// @Description create Notes
// @Param	body		body 	models.Notes	true		"body for Notes content"
// @Success 201 {int} models.Notes
// @Failure 403 body is empty
// @router / [post]
func (c *NotesController) Post() {
	var v models.Notes
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if id, err := models.AddNotes(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = &NotesData{
				Id:        int(id),
				Body:      v.Body,
				IsDone:    v.IsDone,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Notes by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Notes
// @Failure 403 :id is empty
// @router /:id [get]
func (c *NotesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetNotesById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = &NotesData{
			Id:        id,
			Body:      v.Body,
			IsDone:    v.IsDone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Notes
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Notes
// @Failure 403
// @router / [get]
func (c *NotesController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllNotes(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		var res []NotesData

		for _, note := range l {
			res = append(res, NotesData{
				Id:        note.Id,
				Body:      note.Body,
				IsDone:    note.IsDone,
				CreatedAt: note.CreatedAt,
				UpdatedAt: note.UpdatedAt,
			})
		}
		c.Data["json"] = res
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Notes
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Notes	true		"body for Notes content"
// @Success 200 {object} models.Notes
// @Failure 403 :id is not int
// @router /:id [put]
func (c *NotesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := NotesData{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateNotesById(&models.Notes{
			Id:     v.Id,
			Body:   v.Body,
			IsDone: v.IsDone,
		}); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Notes
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *NotesController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteNotes(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
