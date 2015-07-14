package models

import (
    "github.com/orc/db"
    "log"
    "strconv"
)

type GroupsModel struct {
    Entity
}

type Groups struct {
    Id    int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name  string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Owner int    `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) Groups() *GroupsModel {
    model := new(GroupsModel)

    model.TableName = "groups"
    model.Caption = "Группы"

    model.Columns = []string{"id", "name", "face_id"}
    model.ColNames = []string{"ID", "Название", "Владелец"}

    model.Fields = new(Groups)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = true
    model.SubTable = []string{"persons"}
    model.SubField = "group_id"

    return model
}

func (this *GroupsModel) Update(isAdmin bool, userId int, params, where map[string]interface{}) {
    faceId := -1
    query := `SELECT groups.face_id FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND groups.id = $2;`
    err := db.QueryRow(query, []interface{}{userId, where["id"]}).Scan(&faceId)

    if !isAdmin {
        if err != nil || faceId == -1 {
            log.Println("Нет прав редактировать эту группу")
            log.Println(err.Error())
            return
        }
        params["face_id"] = faceId
    }

    this.LoadModelData(params)
    this.LoadWherePart(where)
    db.QueryUpdate(this).Scan()
}

func (this *GroupsModel) Add(userId int, params map[string]interface{}) error {
    var faceId int
    query := `SELECT faces.id
        FROM registrations
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN users ON faces.user_id = users.id
        WHERE users.id = $1 AND events.id = $2;`
    db.QueryRow(query, []interface{}{userId, 1}).Scan(&faceId)
    params["face_id"] = faceId
    this.LoadModelData(params)
    db.QueryInsert(this, "").Scan()
    return nil
}

func (this *GroupsModel) Delete(id int) {
    query := `DELETE
        FROM persons
        WHERE persons.group_id = $1;`
    db.Query(query, []interface{}{id})

    query = `DELETE FROM groups WHERE id = $1;`
    db.Query(query, []interface{}{id})
}

func (this *GroupsModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "groups.id, "
            break
        case "name":
            query += "groups.name as group_name, "
            break
        case "face_id":
            query += "array_to_string(array_agg(param_values.value), ' ') as face_name, "
            break
        }
    }

    query = query[:len(query)-2]

    query += ` FROM param_values
        INNER JOIN registrations ON registrations.id = param_values.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN groups ON groups.face_id = faces.id`

    where, params, _ := this.Where(filters, 1)

    if where != "" {
        query += ` WHERE ` + where + ` AND params.id in (5, 6, 7) AND events.id = 1 GROUP BY groups.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) AND events.id = 1 GROUP BY groups.id`
    }

    if sidx != "" {
        query += ` ORDER BY groups.`+sidx
    }

    query += ` `+ sord

    if limit != -1 {
        params = append(params, limit)
        query += ` LIMIT $`+strconv.Itoa(len(params))
    }

    if offset != -1 {
        params = append(params, offset)
        query += ` OFFSET $`+strconv.Itoa(len(params))
    }

    query += `;`

    return db.Query(query, params)
}

func (this *GroupsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    var query, faces string
    if isAdmin {
        query = `SELECT array_to_string(
            array(
                SELECT f.id || ':' || f.id || '-' || array_to_string(
                array(
                    SELECT param_values.value
                    FROM param_values
                    INNER JOIN registrations ON registrations.id = param_values.reg_id
                    INNER JOIN faces ON faces.id = registrations.face_id
                    INNER JOIN events ON events.id = registrations.event_id
                    INNER JOIN params ON params.id = param_values.param_id
                    WHERE param_values.param_id IN (5, 6, 7) AND events.id = 1 AND faces.id = f.id ORDER BY param_values.param_id
                ), ' ')
                FROM param_values
                INNER JOIN registrations as reg ON reg.id = param_values.reg_id
                INNER JOIN faces as f ON f.id = reg.face_id
                INNER JOIN events ON events.id = reg.event_id
                INNER JOIN params as p ON p.id = param_values.param_id
                INNER JOIN users ON users.id = f.user_id GROUP BY f.id ORDER BY f.id
            ), ';') as name;`
    } else {
        query = `SELECT array_to_string(
            array(
                SELECT f.id || ':' || array_to_string(
                array(
                    SELECT param_values.value
                    FROM param_values
                    INNER JOIN registrations ON registrations.id = param_values.reg_id
                    INNER JOIN faces ON faces.id = registrations.face_id
                    INNER JOIN events ON events.id = registrations.event_id
                    INNER JOIN params ON params.id = param_values.param_id
                    WHERE param_values.param_id IN (5, 6, 7) AND events.id = 1 AND faces.id = f.id ORDER BY param_values.param_id
                ), ' ')
                FROM param_values
                INNER JOIN registrations as reg ON reg.id = param_values.reg_id
                INNER JOIN faces as f ON f.id = reg.face_id
                INNER JOIN events ON events.id = reg.event_id
                INNER JOIN params as p ON p.id = param_values.param_id
                INNER JOIN users ON users.id = f.user_id GROUP BY f.id ORDER BY f.id
            ), ';') as name;`
    }

    faces = db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "name",
            "name": "name",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
        },
        2: map[string]interface{} {
            "index": "face_id",
            "name": "face_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": faces},
            "searchoptions": map[string]string{"value": ":Все;"+faces},
        },
    }
}
