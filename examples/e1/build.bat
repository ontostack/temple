..\..\temple.exe src\e1.go
go build
if not exist "e1" mkdir e1
cd e1
..\e1.exe >> e1.go
go fmt
go build
cd ..
