#！/bin/bash

#solc --abi -o ./abi/a ./sol/a.sol
#solc --bin -o ./bin/a ./sol/a.sol
#abigen --abi=./abi/a/BABYTOKEN.abi --pkg=BABYTOKEN --out=babytoken.go --bin=./bin/a/BABYTOKEN.bin

#solc --abi -o ./abi/ATH ./sol/ATH.sol
#abigen --abi=./abi/ATH/IERC20.abi --pkg=ATH --out=ath.go

#solc --abi -o ./abi/USDT ./sol/USDT.sol
#abigen --abi=./abi/USDT/IERC20.abi --pkg=USDT --out=usdt.go

#solc --abi -o ./abi/Invitation ./sol/Invitation.sol
#abigen --abi=./abi/Invitation/invitation.abi --pkg=Invitation --out=invitation.go

#solc --abi -o ./abi/IDO ./sol/IDO.sol
#abigen --abi=./abi/IDO/IDO.abi --pkg=IDO --out=ido.go

#solc --abi -o ./abi/SAT ./sol/SAT.sol
#abigen --abi=./abi/SAT/IERC20.abi --pkg=SAT --out=sat.go

#abigen --abi=./abi/router.abi --pkg=router --out=router.go

abigen --abi=./abi/pair.abi --pkg=pair --out=pair.go
