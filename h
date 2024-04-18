[33mcommit 22396eb3d48281e27859f7c0b73620974aac59df[m[33m ([m[1;36mHEAD -> [m[1;32mmigrate-to-api[m[33m, [m[1;31morigin/migrate-to-api[m[33m)[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Fri Apr 19 19:22:08 2024 -0700

    Write migration command

 lexapi/api.go  |  6 [32m+++++[m[31m-[m
 lexdb/lexdb.go | 13 [32m+++++++++++++[m
 main.go        | 36 [32m++++++++++++++++++++++++++++++[m[31m------[m
 3 files changed, 48 insertions(+), 7 deletions(-)

[33mcommit 1b023531797370e9af7e7abebc0789108fc2a63c[m[33m ([m[1;31morigin/dev[m[33m, [m[1;32mdev[m[33m)[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Thu Apr 18 11:21:22 2024 -0700

    Update API endpoints for the word of the day

 main.go | 2 [32m+[m[31m-[m
 1 file changed, 1 insertion(+), 1 deletion(-)

[33mcommit 9238d03ac73e0c8d889d2ec9fad7f5ae3176d51c[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Tue Apr 16 17:43:32 2024 -0700

    Add API service as an optional data store

 README.md                        |  12 [32m+++++++++++[m[31m-[m
 api/types.go                     | 163 [31m----------------------------------------------------------------------------------------------------------------------------------------------------------[m
 api/api.go => dictapi/dictapi.go | 236 [32m+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m[31m----------------------------------------------------------------------------------------------------------[m
 lexapi/api.go                    | 116 [32m++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m
 db/db.go => lexdb/lexdb.go       | 121 [32m+++++++++++++++++++++++++++++++++++++++++++++++++++++[m[31m-------------------------------------------------------------[m
 main.go                          | 219 [32m++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m[31m-------------------------------------------------------------------------------------------------[m
 types/types.go                   |  95 [32m++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m
 util/util.go                     |  15 [32m++++++++++++++[m[31m-[m
 8 files changed, 532 insertions(+), 445 deletions(-)

[33mcommit 15126d56333423c9057b57f5cb714b5758721932[m[33m ([m[1;31morigin/master[m[33m, [m[1;32mmaster[m[33m)[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 13 19:13:08 2024 -0700

    Add migrate-to-api command to move data to the Lexicon service

 api/api.go | 51 [32m+++++++++++++++++++++++++++++++++++++++++++++++++++[m
 db/db.go   | 30 [32m+++++++++++++++++++++++++[m[31m-----[m
 main.go    | 67 [32m+++++++++++++++++++++++++++++++++++++++++++++++++[m[31m------------------[m
 3 files changed, 125 insertions(+), 23 deletions(-)

[33mcommit b09f12f1a26b06c76480a53c7541ca4ab1b2d7af[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Fri Apr 12 17:21:49 2024 -0700

    Rename _utils folder

 {_utils => tools}/README.md           | 0
 {_utils => tools}/get-recent-words.sh | 0
 {_utils => tools}/get-word-count.sh   | 0
 {_utils => tools}/save-words-on-mw.sh | 0
 4 files changed, 0 insertions(+), 0 deletions(-)

[33mcommit 39e3955237a7d907e463294c8f4dd001364ec715[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Thu Apr 11 10:50:12 2024 -0700

    Display pronunciations

 main.go      | 17 [32m+++++++++++++++[m[31m--[m
 util/util.go |  4 [32m++[m[31m--[m
 2 files changed, 17 insertions(+), 4 deletions(-)

[33mcommit 1c0f59fc36089567f1603d1e0774140d14898bcf[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Thu Apr 11 10:15:32 2024 -0700

    Get WODs for a range of dates

 Makefile     |  2 [32m+[m[31m-[m
 api/api.go   |  7 [32m+++++[m[31m--[m
 main.go      | 41 [32m+++++++++++++++++++++++++++++++++[m[31m--------[m
 util/util.go |  7 [32m++++++[m[31m-[m
 4 files changed, 45 insertions(+), 12 deletions(-)

[33mcommit 47c8860a19ffa555fb0b36bacd59ba0cca783b5a[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Wed Apr 10 19:24:23 2024 -0700

    Include cross references

 Makefile     |  5 [32m+++++[m
 api/api.go   | 18 [32m++++++++++++++++++[m
 api/types.go | 17 [32m+++++++++++++++++[m
 main.go      | 29 [32m++++++++++++++++++++++++[m[31m-----[m
 util/util.go |  9 [32m+++++++++[m
 5 files changed, 73 insertions(+), 5 deletions(-)

[33mcommit 5b218e2cd6bb1b5c30de2ba474adb20343679c23[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Wed Apr 10 10:53:28 2024 -0700

    Escape 'word' parameter

 _utils/save-words-on-mw.sh | 7 [32m++++[m[31m---[m
 1 file changed, 4 insertions(+), 3 deletions(-)

[33mcommit 4f4ea1bc72e1cc405c069a3300180165c5363c5b[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Tue Apr 9 13:30:52 2024 -0700

    Implement `wod` command
    
    The `wod` command retrieves the word of the day (w/o definition) from
    WOD service. This commands accepts an optional `--date` parameter for
    retrieving a word on a specific date.

 api/api.go   | 34 [32m++++++++++++++++++++++++++++++++[m[31m--[m
 api/types.go |  6 [32m++++++[m
 main.go      | 68 [32m++++++++++++++++++++++++++++++++++++++++++++[m[31m------------------------[m
 3 files changed, 82 insertions(+), 26 deletions(-)

[33mcommit 9adc3d2d2dd71b10f38d70b0e74de63dfd8f83a2[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Mon Apr 8 13:24:48 2024 -0700

    Add get-word-count.sh utility

 _utils/get-word-count.sh | 2 [32m++[m
 1 file changed, 2 insertions(+)

[33mcommit 8ab6251ae8cf1a9c994f150cfbc37248a36a4853[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Mon Apr 8 12:45:57 2024 -0700

    Add utilities

 _utils/README.md           |  2 [32m++[m
 _utils/get-recent-words.sh |  6 [32m++++++[m
 _utils/save-words-on-mw.sh | 19 [32m+++++++++++++++++++[m
 3 files changed, 27 insertions(+)

[33mcommit 551e83c912e9b9b3b8d0d9dac5a16930484b2cc6[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sun Apr 7 21:36:20 2024 -0700

    Abridge output so as to fit on a screen w/o scrolling

 main.go | 95 [32m++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m[31m---------------------------------------[m
 1 file changed, 56 insertions(+), 39 deletions(-)

[33mcommit 2aef03ff812f75aa252ca8f814c5bcbe60f843a0[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sun Apr 7 17:14:33 2024 -0700

    Ignore lexicon name case

 main.go | 4 [32m++[m[31m--[m
 1 file changed, 2 insertions(+), 2 deletions(-)

[33mcommit a2e62576127dbacef565b4ee821d409e13022dcb[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sun Apr 7 08:50:33 2024 -0700

    Skip failed definitions in batch processing

 main.go | 11 [32m++++++++++[m[31m-[m
 1 file changed, 10 insertions(+), 1 deletion(-)

[33mcommit b51c1024192206fb34aee4afa5c84deaf038803c[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sun Apr 7 08:47:37 2024 -0700

    Fix URL encoding bug
    
    Spaces before the `?` mark should be replaced with `%20`, after it, you
    can use use either `+` or `%20`.

 api/api.go | 5 [32m++++[m[31m-[m
 1 file changed, 4 insertions(+), 1 deletion(-)

[33mcommit 3a1b394c3d3f33459f8edb37dbd37ebf70b2b7cf[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 23:40:07 2024 -0700

    Implement define-batch command
    
    This command is for importing words from another source. I personally
    have my saved words from the merriam-webster.com website which I want to
    bring to my personal database with their original saved dates.

 .gitignore |  1 [32m+[m
 db/db.go   | 17 [32m+++++++++++++++[m[31m--[m
 main.go    | 67 [32m++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m[31m---[m
 3 files changed, 80 insertions(+), 5 deletions(-)

[33mcommit 19b4a8ae1617fa6eb8cf6b1e65f7a05d4996fc27[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 21:29:14 2024 -0700

    Fix time stamp bugs

 db/db.go | 3 [32m++[m[31m-[m
 1 file changed, 2 insertions(+), 1 deletion(-)

[33mcommit e17132548659f796c5d47acd23628270018a44e0[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 20:57:16 2024 -0700

    Print quotes

 main.go | 13 [32m+++++++++++++[m
 1 file changed, 13 insertions(+)

[33mcommit 6c09d912023721dfe6d2ce6c8d17eea2f74c1161[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 20:42:16 2024 -0700

    Handles spelling suggestions

 api/api.go | 23 [32m++++++++++++++++++++++[m[31m-[m
 1 file changed, 22 insertions(+), 1 deletion(-)

[33mcommit af789d15433c65f597095b5d445f1d8caf5255f6[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 18:45:21 2024 -0700

    Fix DEntry.Et type
    
    The `et` field can contain strings or arrays, and probably other types.
    Ran into this problem when parsing the response for the "read" word.

 api/types.go | 16 [32m++++++++[m[31m--------[m
 1 file changed, 8 insertions(+), 8 deletions(-)

[33mcommit 81d25b9989a2aeb26fe1631b191121aad8bf0421[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 18:28:01 2024 -0700

    Color titles and subtitles

 README.md | 11 [32m+++++++++++[m
 go.mod    |  5 [32m++++[m[31m-[m
 go.sum    |  4 [32m++++[m
 main.go   | 67 [32m++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m[31m-----------[m
 4 files changed, 75 insertions(+), 12 deletions(-)

[33mcommit 75ff489fe323e1603ba1c6ad1fc842f5051e37fe[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 12:04:24 2024 -0700

    Parse API response and store definitions as JSON

 api/api.go   | 203 [32m+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m
 api/types.go | 140 [32m++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m
 main.go      |  85 [32m++++++++++++++++++++++++[m[31m-------------------------------------------------------------[m
 3 files changed, 367 insertions(+), 61 deletions(-)

[33mcommit be32b735f7a20ae58d58875a33bd5c0b9ecb46f0[m
Author: Rafael Rendon Pablo <rafaelrendonpablo@gmail.com>
Date:   Sat Apr 6 10:38:05 2024 -0700

    Implement basic interactive lexicon

 .gitignore   |   3 [32m+++[m
 README.md    |  23 [32m+++++++++++++++++++++++[m
 db/db.go     | 125 [32m+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m
 go.mod       |   5 [32m+++++[m
 go.sum       |  14 [32m++++++++++++++[m
 main.go      | 191 [32m+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++[m
 schema.sql   |   8 [32m++++++++[m
 util/util.go |  22 [32m++++++++++++++++++++++[m
 8 files changed, 391 insertions(+)
