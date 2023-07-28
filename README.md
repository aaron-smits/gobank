# gobank

This project stems from a tutorial I found from [Anthony GG](https://www.youtube.com/watch?v=pwZuNmAzaH8).

I also referenced [this youtube video](https://www.youtube.com/watch?v=p08c0-99SyU) to learn more about using docker-compose.

My main goal with this project is the build a basic JSON API with Go to understand how developing web backends in Go works and the idioms Go developers use.

On the to do list:
- [x] Transfer endpoint
  - [x] Implement storage method for transfer
  - [X] Implement add balance method
  - [X] Implement subtract balance method
  - [X] Implement seeding method for balance
  - [X] Implement transfer endpoint
- [x] Error handling enhancements
- [x] Add ability for admins to update accounts 
- [ ] Logout endpoint
- [ ] Method to invalidate JWT tokens

Future Tasks
- [ ] Transaction history table
- [ ] Testing (unit, integration, end-to-end)
- [ ] Create a client UI using Go html templates or HTMX
