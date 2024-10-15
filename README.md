## SOCIAL NETWORK DEVELOPMENT

## Authors

![fa-crown](fa-crown.svg)
_Aleksei Gedz_[gitea](https://01.kood.tech/git/agedz)

_Uljana Filippova_ [gitea](https://01.kood.tech/git/ufilippo)

_Dmitri Leljutin_ [gitea](https://01.kood.tech/git/dleljuti)

_Anti Sults_ [gitea](https://01.kood.tech/git/antisults)

_Andreas Veber_ [gitea](https://01.kood.tech/git/aveber)


### PROJECT OBJECTIVES
[link](https://github.com/01-edu/public/tree/master/subjects/social-network)

### AUDIT LINK
[link](https://github.com/01-edu/public/tree/master/subjects/social-network/audit)

### SET UP
run `npm install`
### RUNNING FRONTEND SERVER
frontend runs on port :3000  

`npm run dev`
[port:3000](http://localhost:3000)
(_server set to use range 3000-3010_)

### RUNNING BACKEND
backend runs on port :8080
`cd /backend`
`go run .`
[port:8080](http://localhost:3000)

### START PROJECT WITH DOCKER
Run Docker
Go to the folder "social-network" and use:

`docker build -t front-image .`

Go to the folder "backend" and use:

`docker build -t back-image .`

Go back to the "social-network" folder and use:

`docker-compose up`

[port:3000](http://localhost:3000)

## Directories
Info about different directories

### Front end
### `/app`

- this is where entire frontend is based, page.tsx here is html for home page

#### `/components`

- reusable components for react

#### `/context`

- user context

#### `/hooks`

- for loggingout user on browser tab close
- for WS chat frontend
- for WS notifications

#### `/login` & `register`

Different pages (page.tsx in both these directories is the page itself)

#### `/lib`

- api for posts

#### `users[id]/page.tsx`

- user's root page

#### `users[id]/chat`

- user's chat

#### `users[id]/events`

- user's group events

#### `users[id]/groups`

- user's groups

#### `users[id]/posts`

- user's posts

#### `users[id]/profile`

- user's profile

#### `/utils`

- authentication
- check login status
- style for chat group and regular
- search user handler
- interfaces

####

### Back end

### `/backend`

Entire back-end is located here, also `main.go` is here which initializes the entire back end

#### `/db`

- Database location

##### `/migrations/sqlite`

- Migrations location

##### `/sqlite`

- Calling the migrations, also used for opening the database, all methods

#### `/handlers`

- Handles requests from http, they're set up in `/routes`  

#### `/middleware`

- Holds middleware functions between server and frontend  
- validates if user is logged in, handles CORS, handles server errors, sets range of ports available for frontend, two funcs retreiving user id and user either from map or from DB. 

#### `/routes`

- Sets up handlers

#### `/structs`

- File for structs

#### `/public/uploads`

- Right now holds avatar images for users

#### `/sockets/`

- WebSocket backend server
