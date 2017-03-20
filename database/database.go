package database

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/requires"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net"
	"net/url"
	"time"
)

var (
	Session *mgo.Session
)

type Database struct {
	session  *mgo.Session
	database *mgo.Database
}

func (d *Database) Close() {
	d.session.Close()
}

func (d *Database) getCollection(name string) (coll *Collection) {
	coll = &Collection{
		*d.database.C(name),
		d,
	}
	return
}

func (d *Database) Settings() (coll *Collection) {
	coll = d.getCollection("settings")
	return
}

func (d *Database) Events() (coll *Collection) {
	coll = d.getCollection("events")
	return
}

func (d *Database) Builds() (coll *Collection) {
	coll = d.getCollection("builds")
	return
}

func (d *Database) BuildsLog() (coll *Collection) {
	coll = d.getCollection("builds_log")
	return
}

func (d *Database) PkgGrid() *mgo.GridFS {
	return d.database.GridFS("pkg")
}

func (d *Database) PkgBuildGrid() *mgo.GridFS {
	return d.database.GridFS("pkg_build")
}

func Connect() (err error) {
	mgoUrl, err := url.Parse(config.Config.MongoUri)
	if err != nil {
		err = &ConnectionError{
			errors.Wrap(err, "database: Failed to parse mongo uri"),
		}
		return
	}

	vals := mgoUrl.Query()
	mgoSsl := vals.Get("ssl")
	mgoSslCerts := vals.Get("ssl_ca_certs")
	vals.Del("ssl")
	vals.Del("ssl_ca_certs")
	mgoUrl.RawQuery = vals.Encode()
	mgoUri := mgoUrl.String()

	if mgoSsl == "true" {
		info, e := mgo.ParseURL(mgoUri)
		if e != nil {
			err = &ConnectionError{
				errors.Wrap(e, "database: Failed to parse mongo url"),
			}
			return
		}

		info.DialServer = func(addr *mgo.ServerAddr) (
			conn net.Conn, err error) {

			tlsConf := &tls.Config{}

			if mgoSslCerts != "" {
				caData, e := ioutil.ReadFile(mgoSslCerts)
				if e != nil {
					err = &CertificateError{
						errors.Wrap(e, "database: Failed to load certificate"),
					}
					return
				}

				caPool := x509.NewCertPool()
				if ok := caPool.AppendCertsFromPEM(caData); !ok {
					err = &CertificateError{
						errors.Wrap(err,
							"database: Failed to parse certificate"),
					}
					return
				}

				tlsConf.RootCAs = caPool
			}

			conn, err = tls.Dial("tcp", addr.String(), tlsConf)
			return
		}
		Session, err = mgo.DialWithInfo(info)
		if err != nil {
			err = &ConnectionError{
				errors.Wrap(err, "database: Connection error"),
			}
			return
		}
	} else {
		Session, err = mgo.Dial(mgoUri)
		if err != nil {
			err = &ConnectionError{
				errors.Wrap(err, "database: Connection error"),
			}
			return
		}
	}

	Session.SetMode(mgo.Strong, true)

	return
}

func GetDatabase() (db *Database) {
	session := Session.Copy()
	database := session.DB("")

	db = &Database{
		session:  session,
		database: database,
	}
	return
}

func addIndexes() (err error) {
	db := GetDatabase()
	defer db.Close()

	coll := db.Builds()
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{
			"name",
			"version",
			"release",
			"repo",
			"arch",
		},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
		return
	}

	err = coll.EnsureIndex(mgo.Index{
		Key: []string{
			"state_rank",
			"name",
		},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
		return
	}

	coll = db.BuildsLog()
	err = coll.EnsureIndex(mgo.Index{
		Key: []string{
			"b",
			"t",
		},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
		return
	}

	coll = db.Events()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"channel"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	return
}

func addCollections() (err error) {
	db := GetDatabase()
	defer db.Close()
	coll := db.Events()

	names, err := db.database.CollectionNames()
	if err != nil {
		err = ParseError(err)
		return
	}

	for _, name := range names {
		if name == "events" {
			return
		}
	}

	err = coll.Create(&mgo.CollectionInfo{
		Capped:   true,
		MaxDocs:  1000,
		MaxBytes: 5242880,
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func init() {
	module := requires.New("database")
	module.After("config")

	module.Handler = func() {
		for {
			err := Connect()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("database: Connection")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}

		for {
			err := addCollections()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("database: Add collections")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}

		for {
			err := addIndexes()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("database: Add indexes")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}
	}
}
