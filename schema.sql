CREATE TABLE "lexicon" (
	"name"	        TEXT NOT NULL,
	"definition"	TEXT,
    "source"        TEXT,
    "createdAt"	    INTEGER NOT NULL,
    "updatedAt"	    INTEGER NOT NULL,
	PRIMARY KEY("name")
)
