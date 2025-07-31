module github.com/jpfluger/alibs-slim/aconns/g-aconns

go 1.24.3

toolchain go1.24.5

replace github.com/jpfluger/alibs-slim => ../../../alibs-slim

require (
	github.com/jpfluger/alibs-slim v0.9.7
	github.com/jpfluger/alibs-slim/aconns/aclient-ftp v0.9.7
	github.com/jpfluger/alibs-slim/aconns/aclient-http v0.9.7
	github.com/jpfluger/alibs-slim/aconns/aclient-ldap v0.9.7
	github.com/jpfluger/alibs-slim/aconns/aclient-redis v0.9.7
	github.com/jpfluger/alibs-slim/aconns/aclient-smtp v0.9.7
	github.com/jpfluger/alibs-slim/aconns/adb-mssql v0.9.7
	github.com/jpfluger/alibs-slim/aconns/adb-mysql v0.9.7
	github.com/jpfluger/alibs-slim/aconns/adb-oracle v0.9.7
	github.com/jpfluger/alibs-slim/aconns/adb-pg v0.9.7
	github.com/stretchr/testify v1.10.0
)

require (
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/alexedwards/scs/redisstore v0.0.0-20250417082927-ab20b3feb5e9 // indirect
	github.com/anthonynsimon/bild v0.14.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.8-0.20250403174932-29230038a667 // indirect
	github.com/go-ldap/ldap/v3 v3.4.11 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/gofrs/uuid/v5 v5.3.2 // indirect
	github.com/gogs/chardet v0.0.0-20211120154057-b7413eaefb8f // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/gomig/avatar v1.0.3 // indirect
	github.com/gomig/utils v1.0.1 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hjson/hjson-go/v4 v4.5.0 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/jhillyerd/enmime/v2 v2.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jlaffaye/ftp v0.0.0-20220301011324-fed5bc26b7fa // indirect
	github.com/labstack/echo/v4 v4.13.4 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mileusna/timezones v0.0.0-20220627120747-ad570b2850c0 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354 // indirect
	github.com/olekukonko/tablewriter v1.0.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/otp v1.5.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/sethvargo/go-password v0.3.1 // indirect
	github.com/sijms/go-ora/v2 v2.9.0 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun v1.2.15 // indirect
	github.com/uptrace/bun/dialect/mssqldialect v1.2.15 // indirect
	github.com/uptrace/bun/dialect/mysqldialect v1.2.15 // indirect
	github.com/uptrace/bun/dialect/pgdialect v1.2.15 // indirect
	github.com/uptrace/bun/driver/pgdriver v1.2.15 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.2 // indirect
)

replace github.com/jpfluger/alibs-slim/aconns/aclient-ftp => ../aclient-ftp

replace github.com/jpfluger/alibs-slim/aconns/aclient-http => ../aclient-http

replace github.com/jpfluger/alibs-slim/aconns/aclient-ldap => ../aclient-ldap

replace github.com/jpfluger/alibs-slim/aconns/aclient-redis => ../aclient-redis

replace github.com/jpfluger/alibs-slim/aconns/aclient-sftp => ../aclient-sftp

replace github.com/jpfluger/alibs-slim/aconns/aclient-smtp => ../aclient-smtp

replace github.com/jpfluger/alibs-slim/aconns/adb-mssql => ../adb-mssql

replace github.com/jpfluger/alibs-slim/aconns/adb-mysql => ../adb-mysql

replace github.com/jpfluger/alibs-slim/aconns/adb-oracle => ../adb-oracle

replace github.com/jpfluger/alibs-slim/aconns/adb-pg => ../adb-pg
