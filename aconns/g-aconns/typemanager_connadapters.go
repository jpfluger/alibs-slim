package g_aconns

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	aclient_ftp "github.com/jpfluger/alibs-slim/aconns/aclient-ftp"
	aclient_http "github.com/jpfluger/alibs-slim/aconns/aclient-http"
	aclient_ldap "github.com/jpfluger/alibs-slim/aconns/aclient-ldap"
	aclient_redis "github.com/jpfluger/alibs-slim/aconns/aclient-redis"
	aclient_smtp "github.com/jpfluger/alibs-slim/aconns/aclient-smtp"
	adb_mssql "github.com/jpfluger/alibs-slim/aconns/adb-mssql"
	adb_mysql "github.com/jpfluger/alibs-slim/aconns/adb-mysql"
	adb_oracle "github.com/jpfluger/alibs-slim/aconns/adb-oracle"
	adb_pg "github.com/jpfluger/alibs-slim/aconns/adb-pg"
	"github.com/jpfluger/alibs-slim/areflect"
	"reflect"
)

// init registers the connection adapters with the TypeManager.
func init() {
	_ = areflect.TypeManager().Register(aconns.TYPEMANAGER_CONNADAPTERS, "gconns", returnTypeManagerConnAdapters)
}

// returnTypeManagerConnAdapters returns the reflect.Type for the given adapter type name.
func returnTypeManagerConnAdapters(typeName string) (reflect.Type, error) {
	var rtype reflect.Type

	switch aconns.AdapterType(typeName) {
	case aclient_ftp.ADAPTERTYPE_FTP:
		rtype = reflect.TypeOf(aclient_ftp.AClientFTP{})
	case aclient_http.ADAPTERTYPE_HTTP:
		rtype = reflect.TypeOf(aclient_http.AClientHTTP{})
	case aclient_ldap.ADAPTERTYPE_LDAP:
		rtype = reflect.TypeOf(aclient_ldap.AClientLDAP{})
	case aclient_redis.ADAPTERTYPE_REDIS:
		rtype = reflect.TypeOf(aclient_redis.AClientRedis{})
	case aclient_smtp.ADAPTERTYPE_SMTP:
		rtype = reflect.TypeOf(aclient_smtp.AClientSMTP{})
	case adb_mssql.ADAPTERTYPE_MSSQL:
		rtype = reflect.TypeOf(adb_mssql.ADBMSSql{})
	case adb_mysql.ADAPTERTYPE_MYSQL, adb_mysql.ADAPTERTYPE_MARIA:
		rtype = reflect.TypeOf(adb_mysql.ADBMysql{})
	case adb_oracle.ADAPTERTYPE_ORACLE:
		rtype = reflect.TypeOf(adb_oracle.ADBOracle{})
	case adb_pg.ADAPTERTYPE_PG:
		rtype = reflect.TypeOf(adb_pg.ADBPG{})
	default:
		return nil, fmt.Errorf("unknown adapter type: %s", typeName)
	}

	return rtype, nil
}
