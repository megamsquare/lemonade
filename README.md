# Lemonade Test

## To test the result
run this command

```bash
go run main.go
```

when the server starts, It should start at port 3000

you can visit

```text
# This endpoit takes in username to create a user: The method is POST
localhost:3000/user

# This endpoint returns all users information stored  in the memory
localhost:3000/user

#This endpoint takes in sender_id, receiver_id and amount: The method is post
localhost:3000/transaction
```