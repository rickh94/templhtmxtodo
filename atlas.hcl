variable "dbpath" {
  type = string
  default = getenv("DB_PATH")
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://dev?mode=memory&_fk=1&_journal_mode=WAL"
  url = "sqlite://${var.dbpath}?_fk=1&_journal_mode=WAL"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  src = data.external_schema.gorm.url
  url = "sqlite:///${var.dbpath}?_fk=1&_journal_mode=WAL"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
