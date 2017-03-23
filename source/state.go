package source

import "github.com/autoabs/autoabs/database"

type state struct {
	Repo   string `bson:"_id"`
	Commit string `bson:"commit"`
}

func getCommit(db *database.Database, repo string) (commit string, err error) {
	coll := db.SourcesState()

	doc := &state{}

	err = coll.FindOneId(repo, doc)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	commit = doc.Commit

	return
}

func setCommit(db *database.Database, repo, commit string) (err error) {
	coll := db.SourcesState()

	doc := &state{
		Repo:   repo,
		Commit: commit,
	}

	_, err = coll.UpsertId(repo, doc)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
