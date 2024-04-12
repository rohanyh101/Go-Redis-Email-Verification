# Email Verification Via Redis & "http/smtp" go package

## Run Instructions,
1. Docker must be installed on your machine,
2. Run the following,
```powershell
docker-compose up
```
3. configure the `.env` file as per your domain provider & email
4. Run,
```powershell
go run .
```

## Verify these endpoints,
- route => "/send-verification-email"
```powershell
curl --location --request GET 'localhost:3000/send-verification-email?email=example@example.com' `
>> --header 'Content-Type: application/json'
```

- route => "/verify-email"
```powershell
curl --location --request GET 'localhost:3000/verify-email?token=token=35f914255ca65296e8e027570f081cfde124d3daa7df7cb9be0ae21619bf208cexample@example.com' `
>> --header 'Content-Type: application/json'
```
