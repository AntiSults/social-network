# ToDo list

- [x] Select a framework (Next.js)

- [ ] Create Basic Authentication
    - [ ] Create User Table in SQlite
    -- Fields: 
    - Email
    - Password
    - First Name
    - Last Name
    - Date of Birth
    - Avatar/Image (Optional)
    - Nickname (Optional)
    - About Me (Optional)

- [ ] Menu
    - Pages: My profile, Followers, Groups, Chat

- [ ] Create a user profile
    - [ ] User information field
    - [ ] Public and private user profile
    - [ ] User activity field 
    - [ ] Followers and following field

- [ ] Posts
    - [ ] Create Post Table in SQlite (private and public fields!)
    - [ ] Checkbox for private posts 
    - [ ] Create Comments Table in SQlite (comments need to be linked to post id)
    - [ ] Ability to comment on posts

- [ ] Groups
    - [ ] Create Groups Table in SQlite
    - [ ] Ability to create a group
    - [ ] Section in the menu for groups
    - [ ] Title and description fields
    - [ ] Invitation
    - [ ] Posts for groups (Posts can only be seen by group members)
    - [ ] Field for the event
    - [ ] Event -- Fields:
        - Title
        - Description
        - Day/Time
        - 2 Options (at least):
            - Going
            - Not going
    - [ ] Chat for group members

- [ ] Chat (websocket)
    - [ ] Create Messages Table in SQlite
    - [ ] Ability to send private messages to other users that they are following or being followed
    - [ ] Ability to send emojis

- [ ] Notifications (websocket)
    - [ ] Create Notifications Table in SQlite
    - [ ] Create a special field for notifications
        - A user should be notified if he/she:
            - has a private profile and some other user sends him/her a following request
            - receives a group invitation, so he can refuse or accept the request
            - is the creator of a group and another user requests to join the group, so he can refuse or accept the request
            - is member of a group and an event is created