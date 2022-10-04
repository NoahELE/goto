go build -tags=jsoniter

cd frontend
call npm run build
@echo on
cd ..

./goto.exe
