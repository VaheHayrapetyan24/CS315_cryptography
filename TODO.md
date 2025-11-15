Distribute

- [ ] Separate init script that writes to parameters
- [ ] Write G_t, D matrices to a file
- [ ] New endpoint for fetching a chunk of G_t
- [ ] Calculate A[i] when asked for, not before that
- [ ] When asked for an i > current n, add a new row to G_t



Node

- [ ] Initialize script that writes to a file
    - [ ] Should fetch current G_t, q, U
    - [ ] Second time called should do noop
    - [ ] Read/write permissions for current user only

- [ ] Rerfresh script for updating G_t
- [ ] When initing server, read U from the file, cache it
- [x] Specify port when server start command
- [ ] Send message endpoint that calls another node with given id and message
    - [ ] Encryption with resulting key
- [ ] Message endpoint that receives message, decrypts it, prints


TurboBlom 
- [ ] matrix consisting of 1s and 0s
- [ ] matrix/vector operations with bit shift/flips
