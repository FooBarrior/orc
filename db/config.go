package db

const user string = "admin"
const dbname string = "orc"
const password string = "admin"

var Tables = []string{
    "events",
    "event_types",
    "events_types",
    "persons",
    "users",
    "forms",
    "param_types",
    "params",
    "events_forms",
    "param_values",
    "faces",
    "registrations",
    "reg_param_vals",
    "groups",
    "group_registrations",
    "regs_groupregs",
    "docs",
    "events_docs",
}

var TableNames = []string{
    "Мероприятия",
    "Типы мероприятий",
    "Мероприятия - Типы",
    "Персоны",
    "Пользователи",
    "Формы",
    "Типы параметров",
    "Параметры",
    "Мероприятия - Формы",
    "Значения параметров",
    "Лица",
    "Регистрации",
    "Регистрации - Значения параметров",
    "Группы",
    "Регистрации групп",
    "Регистрации групп - Регистрации",
    "Документы",
    "Мероприятие - Документы",
}
