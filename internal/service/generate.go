package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ChatService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i FakeTx -o ./mocks/ -s "_minimock.go"
