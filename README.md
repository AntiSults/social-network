## SOCIAL NETWORK DEVELOPMENT

### SET UP

run `npm install`

### RUNNING FRONTEND SERVER

frontend runs on port :3000  
`npm run dev`

### RUNNING BACKEND

backend runs on port :8080  
`cd /backend`
`go run main.go`

### SO FAR IMPLEMENTED

Front end sends registration form to backend, including avatar  
Also creates a randomly named image file in the /public/uploads directory  
Early implementation of migrating sql tables  
Creates a session token for login, inserts to database  
Authenticates login  
Able to logout  
Changes homescreen based on if logged in or not  
Fetches avatar and displays it (also has a default avatar)
Checks email against regex, password has to be atleast 4 characters
Websockets connection established, handler in backend. Sends and returns the same message for checking the connection.

## Directories

Info about different directories

### Front end

### `/app`

this is where entire frontend is based, page.tsx here is html for home page

#### `/components`

reusable components for react

#### `/login` & `register`

Different pages (page.tsx in both these directories is the page itself)

#### `/profile`

User's profile page

#### `/testLoggedIn`

Test page for testing if user is logged in or not **MUST DELETE ON PRODUCTION**

#### `/utils`

Folder for utilities, so far authentication

####

### Back end

### `/backend`

Entire back-end is located here, also `main.go` is here which initializes the entire back end

#### `/db`

Database location

##### `/migrations/sqlite`

Migrations location

##### `/sqlite`

Calling the migrations, also used for opening the database

#### `/handlers`

Handles requests from http, they're set up in `/routes`  
Register, Login, Logout, GetAvatar

#### `/middleware`

Middleware functions between server and frontend

#### `/routes`

Sets up handlers

#### `/structs`

File for structs

#### `/public/uploads`

Right now holds avatar images for users

#### `/middleware`

Holds middleware functions between server and frontend  
So far validates if user is logged in, handles CORS, handles server errors

#### `/sockets/manager`

Add new handler (method) here for establishing new kind of connection.
Each handler will need corresponding map in the 'setupEventHandlers()' method.
This map will have the key(string with name of event) - the constant in the event.go file.


## Default readme

This is a [Next.js](https://nextjs.org/) project bootstrapped with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/basic-features/font-optimization) to automatically optimize and load Inter, a custom Google Font.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/deployment) for more details.
