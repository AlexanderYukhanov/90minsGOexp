# 90minsGOexp
useless experiment

## Before clock start
Initial version of swagger/swagger.yaml

## After clock ends (90 mins)
- The final commit and push
- Added forgoten go.mod
- This doc


## Not done, known defects, etc
- List available appointments not implemented
- Forgot to convert Conflict storage response to 409 response, so will give 500 or 503
- MySQL or other types of persistent storage not implemented
- Should have a service per route
- Not enough unit tests
- Consts should come from config
- Error codes must be introduced via swagger
- Times in error attributes must be in the same format as in input/output

## Instruction
- make swagger - generate server code
- make build - builds the service
- make start - starts the service on 8080 port
- make live-test - is supposed to send some queries and validate responses
