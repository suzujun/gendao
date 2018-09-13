# Gendao

[![CircleCI](https://circleci.com/gh/suzujun/gendao/tree/master.svg?style=shield&circle-token=a6f399037d30cbc591f227f63163f77f8c5ac976)](https://circleci.com/gh/suzujun/gendao/tree/master)
[![Build Status](https://travis-ci.org/suzujun/gendao.svg?branch=master)](https://travis-ci.org/suzujun/gendao)
[![Coverage Status](https://coveralls.io/repos/github/suzujun/gendao/badge.svg)](https://coveralls.io/github/suzujun/gendao) [![Code Climate](https://codeclimate.com/github/suzujun/gendao/badges/gpa.svg)](https://codeclimate.com/github/suzujun/gendao) [![codecov](https://codecov.io/gh/suzujun/gendao/branch/master/graph/badge.svg)](https://codecov.io/gh/suzujun/gendao)


Gendao is generate DAO and Model source code using templates.


## Install

``` bash
go get github.com/suzujun/gendao
```

## Usage

Gendao provides these commands.

* `init` - Create initialized JSON file
* `pull` - Generate JSON of table struct from database
* `addtype` - Set your own type for the column in the table
* `gen` - Generate Dao and Model from json schema

### Example

``` bash
# Setting the database information
$ gendao init -d database_name > config.json

# Generate JSON schema for table
$ gendao pull config.json

# Generate source code
$ gendao gen config.json -t tablename1,tablename2
```

## gendao commands
### gendao init
Create initialized JSON file.

* `host` or `H` - host name to connect to the database (`localhost` by default)
* `port` or `P` - port to connect to the database (`3306` by default)
* `user` or `u` - user name to connect to the database (`root` by default)
* `password` or `p` - password to connect to the database (empty value by default)
* `database` or `d` - database to be processed (The value of the config is used as the default)

### gendao pull [config name]
Generate a JSON of table struct. This command has these flag options.

* `database` - database to be processed (The value of the config is used as the default)

### gendao addtype [config name]
Set your own type for the column in the table.
Follow the wizard and enter necessary items.

### gendao gen [config name]
Generate a source code. This command has these flag options.

* `database` - database to be processed (The value of the config is used as the default)
* `table` - tables to be processed (select all by default)

# License

MIT
