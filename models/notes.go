package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Notes struct {
	Id        int       `orm:"column(id);pk;auto"`
	Body      string    `orm:"column(body);null"`
	IsDone    bool      `orm:"column(is_done);null"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(timestamp without time zone);null"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now;type(timestamp without time zone);null"`
}

func (t *Notes) TableName() string {
	return "notes"
}

func init() {
	orm.RegisterModel(new(Notes))
}

// AddNotes insert a new Notes into database and returns
// last inserted Id on success.
func AddNotes(m *Notes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetNotesById retrieves Notes by Id. Returns error if
// Id doesn't exist
func GetNotesById(id int) (v *Notes, err error) {
	o := orm.NewOrm()
	v = &Notes{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllNotes retrieves all Notes matches certain condition. Returns empty list if
// no records exist
func GetAllNotes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Notes, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Notes))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Notes
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		for _, v := range l {
			ml = append(ml, v)
		}
		return ml, nil
	}
	return nil, err
}

// UpdateNotes updates Notes by Id and returns error if
// the record to be updated doesn't exist
func UpdateNotesById(m *Notes) (err error) {
	o := orm.NewOrm()
	v := Notes{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteNotes deletes Notes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteNotes(id int) (err error) {
	o := orm.NewOrm()
	v := Notes{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Notes{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
