# GatorCLI

# Description

GatorCLI is a Go command line interface that allows users to follow RSS Feeds and see their posts. Users can register, login, follow and unfollow feeds as well as viewing their posts. More details in the instructions. This is also my first time writing a .md file in a long time as such it is improperly formatted. 

Guided project instructions from Boot.dev.

Note: This should be called gator instead of GatorCLI but I did not clue into that until writing this now. 

## Requirements
- postgresql@15
- Go

## Installation

1) Download this repository.

2) Move into sql/schema ```cd ./sql/schema```

3) Run the migrations with ```goose postgres postgres://[username]:[database] up```

4) From the root run ```go install github.com/CTK-code/GatorCLI ```

## Usage
```
GatorCLI login [username]
``` 
Sets current user to a registered user [username]. Fails if user is not registered

```
GatorCLI register [username]
```
Register a user named [username] and sets current user to that user

```
GatorCLI users
```
List all users

```
GatorCLI addfeed [feed_name] [feed_url]
```
Adds an RSS feed to the database. Requires the name of the feed and it's URL link.

```
GatorCLI feeds 
```
Lists all feeds along with who added them

```
GatorCLI follow [feed_url]
```
Adds [feed_url] to the list of feeds that current user is following.

```
GatorCLI unfollow [feed_url]
```
Removes [feed_url] from the list of feeds that current user is following.


```
GatorCLI following 
```
List all feeds current user is following.

```
GatorCLI agg [time_between] 
ex. GatorCLI agg 5s | 5m 
```
Fetch posts from current users followed feeds.

```
GatorCLI browse [limit]
```
Fetch posts from database for feeds that current user follows. Limits to the limit number most recent posts

